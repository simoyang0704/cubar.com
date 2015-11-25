package model

import (
	"cubar.com/lib/util"
	"fmt"
)

const (
	COMMUNITY_USER_STATUS_JOINED  = 1
	COMMUNITY_USER_STATUS_APPLY   = 2
	COMMUNITY_USER_STATUS_DELETED = 9

	COMMUNITY_ROLE_ADMIN   = 1
	COMMUNITY_ROLE_MANAGER = 2
	COMMUNITY_ROLE_MEMBER  = 3
)

type CommunityUser struct {
	CommunityUserId string `json:"community_user_id" db:"community_user_id"`
	CommunityId     int    `json:"community_id" db:"community_id"`
	UserId          int    `json:"user_id" db:"user_id"`
	Status          int    `json:"status" db:"status"`
	Role            int    `json:"role" db:"role"`
	Joined          string `json:"joined" db:"joined"`
	Created         string `json:"created" db:"created"`
	Deleted         string `json:"deleted" db:"deleted"`
}

func (communityUser *CommunityUser) Add() error {

	if communityUser.CommunityId == 0 || communityUser.UserId == 0 {
		return fmt.Errorf("社区关系添加错误，用户或者社区有误")
	}

	if communityUser.Status == 0 {
		communityUser.Status = COMMUNITY_USER_STATUS_APPLY
	}

	if communityUser.Role == 0 {
		communityUser.Role = COMMUNITY_ROLE_MEMBER
	}

	communityUser.CommunityUserId = generateCommunityUserId(communityUser.CommunityId, communityUser.UserId)
	communityUser.Created = util.Now()

	return __handle.Insert(communityUser)
}

func (communityUser *CommunityUser) Delete() error {

	if communityUser.CommunityId == 0 || communityUser.UserId == 0 {
		return fmt.Errorf("要退出的社区有误")
	}

	communityUser.CommunityUserId = generateCommunityUserId(communityUser.CommunityId, communityUser.UserId)
	_, err := __handle.Delete(communityUser)
	return err
}

func (communityUser *CommunityUser) Get() (bool, error) {

	exists := false
	if communityUser.CommunityUserId == "" {
		communityUser.CommunityUserId = generateCommunityUserId(communityUser.CommunityId, communityUser.UserId)
	}

	obj, err := __handle.Get(Community{}, communityUser.CommunityUserId)
	if err != nil {
		return exists, nil
	}

	*communityUser = *(obj.(*CommunityUser))
	exists = true

	return exists, nil
}

func ApplyCommunity(communityId int, userId int) error {

	communityUser := &CommunityUser{
		CommunityId: communityId,
		UserId:      userId,
	}

	exists, err := communityUser.Get()
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return communityUser.Add()
}

func JoinCommunity(communityId int, userId int) error {

	communityUser := &CommunityUser{
		CommunityId: communityId,
		UserId:      userId,
	}

	exists, err := communityUser.Get()
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("没有申请过这个社区")
	}

	if communityUser.Status != COMMUNITY_USER_STATUS_JOINED {
		communityUser.Status = COMMUNITY_USER_STATUS_JOINED
		communityUser.Joined = util.Now()
		_, err := __handle.Update(communityUser)

		return err
	}

	return nil
}

func generateCommunityUserId(communityId, userId int) string {

	return fmt.Sprintf("%v_%v", communityId, userId)
}
