package model

import (
	"cubar.com/lib/util"
	"fmt"
)

const (
	ACTIVITY_STATUS_NORMAL  = 1
	ACTIVITY_STATUS_DELETED = 9

	ACTIVITY_TYPE_ACTIVITY = 1
	ACTIVITY_TYPE_SHARE    = 2

	ACTIVITY_TABLE      = "activities"
	ACTIVITY_USER_TABLE = "activity_users"
)

type Activity struct {
	ActivityId       int     `json:"activity_id" db:"activity_id"`
	CommunityId      int     `json:"community_id" db:"community_id"`
	CityId           int     `json:"city_id" db:"city_id"`
	FollowActivityId int     `json:"follow_activity_id" db:"follow_activity_id"`
	Type             int     `json:"type" db:"type"`
	Name             string  `json:"name" db:"name"`
	Content          string  `json:"content" db:"content"`
	UserNum          int     `json:"user_num" db:"user_num"`
	StartTime        string  `json:"start_time" db:"start_time"`
	EndTime          string  `json:"end_time" db:"end_time"`
	Location         string  `json:"location" db:"location"`
	Latitude         float64 `json:"latitude" db:"latitude"`
	Longitude        float64 `json:"longitude" db:"longitude"`
	CreatedUserId    int     `json:"created_user_id" db:"created_user_id"`
	Status           int     `json:"status" db:"status"`
	Created          string  `json:"created" db:"created"`
	Deleted          string  `json:"deleted" db:"deleted"`
}

func (activity *Activity) Add() error {

	if activity.CommunityId == 0 {
		return fmt.Errorf("活动必须附属于社区")
	}

	community := GetCommunity(activity.CommunityId)
	if community == nil {
		return fmt.Errorf("附属于不存在的社区")
	}

	activity.CityId = community.CityId
	activity.Type = ACTIVITY_TYPE_ACTIVITY
	if activity.FollowActivityId != 0 {
		activity.Type = ACTIVITY_TYPE_SHARE
	}
}
