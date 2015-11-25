package model

import (
	"fmt"
	// "log"

	"cubar.com/conf"
	"cubar.com/lib/util"
	"strings"
)

const (
	COMMUNITY_STATUS_NORMAL  = 1
	COMMUNITY_STATUS_DELETED = 9

	COMMUNITY_TABLE      = "communities"
	COMMUNITY_USER_TABLE = "community_users"

	QUERY_CYCLE   = 0.5
	DEFAULT_LIMIT = 20
)

type Community struct {
	CommunityId   int     `json:"community_id" db:"community_id"`
	CityId        int     `json:"city_id" db:"city_id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	Latitude      float64 `json:"latitude" db:"latitude"`
	Longitude     float64 `json:"longitude" db:"longitude"`
	CreatedUserId int     `json:"created_user_id" db:"created_user_id"`
	UpdatedUserId int     `json:"updated_user_id" db:"updated_user_id"`
	DeletedUserId int     `json:"deleted_user_id" db:"deleted_user_id"`
	OwnerUserId   int     `json:"owner_user_id" db:"owner_user_id"`
	CategoryId    int     `json:"category_id" db:"category_id"`
	SubCategoryId int     `json:"sub_category_id" db:"sub_category_id"`
	Status        int     `json:"status" db:"status"`
	Created       string  `json:"created" db:"created"`
	Updated       string  `json:"updated" db:"updated"`
	Deleted       string  `json:"deleted" db:"deleted"`
	Activated     string  `json:"activated" db:"activated"`
	Distance      float64 `json:"distance"`
}

type CommunityQuery struct {
	CityId        int     `json:"city_id"`
	Name          string  `json:"name"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	OwnerUserId   int     `json:"owner_user_id"`
	CategoryId    int     `json:"category_id"`
	SubCategoryId int     `json:"sub_category_id"`
	Offset        int     `json:"offset"`
	Limit         int     `json:"limit"`
}

func (community *Community) Add() error {

	if community.Status == 0 {
		community.Status = COMMUNITY_STATUS_NORMAL
	}

	if community.CategoryId == 0 {
		community.CategoryId = conf.CATE_OTHER
	}

	if community.OwnerUserId == 0 {
		community.OwnerUserId = community.CreatedUserId
	}

	community.Created = util.Now()
	community.Activated = community.Created
	community.Status = COMMUNITY_STATUS_NORMAL

	if err := __handle.Insert(community); err != nil {
		return err
	}

	// 将创建人加入到该关系中
	communityUser := &CommunityUser{
		CommunityId: community.CommunityId,
		UserId:      community.OwnerUserId,
		Status:      COMMUNITY_USER_STATUS_JOINED,
		Role:        COMMUNITY_ROLE_ADMIN,
		Joined:      util.Now(),
	}

	if err := communityUser.Add(); err != nil {
		return err
	}

	return nil
}

func (community *Community) Get() (bool, error) {

	exists := false
	if community.CommunityId == 0 {
		return exists, fmt.Errorf("get community, invalid cid")
	}

	obj, err := __handle.Get(Community{}, community.CommunityId)
	if err != nil {
		return exists, nil
	}

	*community = *(obj.(*Community))

	if community.Status != COMMUNITY_STATUS_DELETED {
		exists = true
	}

	return exists, nil
}

func (community *Community) Update() error {

	oldCommunity := &Community{
		CommunityId: community.CommunityId,
	}

	if exists, err := oldCommunity.Get(); err != nil {
		return fmt.Errorf("update community(%v), get old community error(%v)", community, err)
	} else if !exists {
		return fmt.Errorf("update community(%v), community not exists", community)
	}

	community.Updated = util.Now()

	if _, err := __handle.Update(community); err != nil {
		return fmt.Errorf("update community(%v), db error(%v)", community, err)
	}

	return nil
}

func (community *Community) Delete() error {

	oldCommunity := &Community{
		CommunityId: community.CommunityId,
	}

	if exists, err := oldCommunity.Get(); err != nil {
		return fmt.Errorf("delete community(%v), get old community error(%v)", community, err)
	} else if !exists {
		return fmt.Errorf("delete community(%v), community not exists", community)
	}

	community.Deleted = util.Now()
	community.Status = COMMUNITY_STATUS_DELETED

	if _, err := __handle.Update(community); err != nil {
		return fmt.Errorf("delete community(%v), db error(%v)", community, err)
	}

	return nil
}

func GetCommunity(communityId int) *Community {

	if communityId == 0 {
		return nil
	}

	// todo 从缓存中获取

	// 从db获取
	community := &Community{
		CommunityId: communityId,
	}

	if exi, err := community.Get(); err != nil || !exi {
		return nil
	}

	return community
}

func (commmunityQuery *CommunityQuery) Query() (communities []Community) {

	sqlstr, params := commmunityQuery.formatSql(false, nil, 0)

	if _, err := __handle.Select(&communities, sqlstr, params...); err != nil {
		return nil
	}

	return
}

func (communityQuery *CommunityQuery) Count() (count int64) {

	sqlstr, params := communityQuery.formatSql(true, nil, 0)

	count, _ = __handle.SelectInt(sqlstr, params...)

	return
}

func (communityQuery *CommunityQuery) QueryByUserAndStatus(user *User, status int) (communities []Community) {

	sqlstr, params := communityQuery.formatSql(false, user, status)

	if _, err := __handle.Select(&communities, sqlstr, params...); err != nil {

		return nil
	}

	return
}

func (communityQuery *CommunityQuery) formatSql(count bool, user *User, status int) (sqlstr string, params []interface{}) {

	selectSql, filterSql, orderSql, filterParams := communityQuery.formatFilterSql(count)

	sqlstr = "SELECT " + selectSql + " FROM " + COMMUNITY_TABLE
	if count {
		sqlstr = "SELECT count(*) FROM " + COMMUNITY_TABLE
	}

	if user != nil {

		sqlstr += " LEFT JOIN " + COMMUNITY_USER_TABLE + " ON " + COMMUNITY_TABLE + ".community_id = " + COMMUNITY_USER_TABLE + ".community_id"
		if filterSql != "" {
			filterSql += " AND"
		}

		filterSql += " " + COMMUNITY_USER_TABLE + ".user_id = ? AND " + COMMUNITY_USER_TABLE + ".status = ?"
		filterParams = append(filterParams, user.UserId, status)
		orderSql = " ORDER BY " + COMMUNITY_TABLE + ".activated DESC"
	}

	if filterSql != "" {
		sqlstr += " WHERE " + filterSql
	}

	if !count {
		if orderSql != "" {
			sqlstr += orderSql
		}

		if communityQuery.Limit == 0 {
			communityQuery.Limit = DEFAULT_LIMIT
		}

		sqlstr += " LIMIT " + util.IntToString(communityQuery.Limit)
		if communityQuery.Offset > 0 {
			sqlstr += " OFFSET " + util.IntToString(communityQuery.Offset)
		}
	}

	params = filterParams

	return
}

func (communityQuery *CommunityQuery) formatFilterSql(count bool) (selectSql string, filterSql string, orderSql string, params []interface{}) {

	selectSql = COMMUNITY_TABLE + ".*"
	sqlstrs := make([]string, 0)
	if communityQuery.CityId != 0 {
		sqlstrs = append(sqlstrs, COMMUNITY_TABLE+".city_id = ?")
		params = append(params, communityQuery.CityId)
	}

	if communityQuery.Name != "" {
		sqlstrs = append(sqlstrs, COMMUNITY_TABLE+".name LIKE ?")
		params = append(params, "%"+communityQuery.Name+"%")
	}

	if communityQuery.OwnerUserId != 0 {
		sqlstrs = append(sqlstrs, COMMUNITY_TABLE+".owner_user_id = ?")
		params = append(params, communityQuery.OwnerUserId)
	}

	if communityQuery.CategoryId != 0 {
		sqlstrs = append(sqlstrs, COMMUNITY_TABLE+".category_id = ?")
		params = append(params, communityQuery.CategoryId)
	}

	if communityQuery.SubCategoryId != 0 {
		sqlstrs = append(sqlstrs, COMMUNITY_TABLE+".sub_category_id = ?")
		params = append(params, communityQuery.SubCategoryId)
	}

	if (communityQuery.Latitude != 0 || communityQuery.Longitude != 0) && !count {
		selectSql += ", calc_distance(?, ?, " + COMMUNITY_TABLE + ".latitude, " + COMMUNITY_TABLE + ".longitude) AS distance"
		orderSql = " ORDER BY distance ASC "
		params = append([]interface{}{communityQuery.Latitude, communityQuery.Longitude}, params...)

		sqlstrs = append(sqlstrs, "abs("+COMMUNITY_TABLE+".latitude - ?) < ? AND abs("+COMMUNITY_TABLE+".longitude - ?) < ?")
		params = append(params, communityQuery.Latitude, QUERY_CYCLE, communityQuery.Longitude, QUERY_CYCLE)
	}

	filterSql = strings.Join(sqlstrs, " AND ")

	return
}
