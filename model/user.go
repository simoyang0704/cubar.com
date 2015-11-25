package model

import (
	"cubar.com/lib/util"
	"fmt"
)

const (
	GENDER_MALE   = "male"
	GENDER_FEMALE = "female"

	USER_STATUS_NORMAL   = 1
	USER_STATUS_UNACTIVE = 2
	USER_STATUS_SILENCE  = 3
	USER_STATUS_DELETED  = 9
)

type User struct {
	UserId   int    `json:"user_id" db:"user_id"`
	Name     string `json:"name" db:"name"`
	Pinyin   string `json:"pinyin" db:"pinyin"`
	Password string `json:"password" db:"password"`
	Gender   string `json:"gender" db:"gender"`
	Age      int    `json:"age" db:"age"`
	Birthday string `json:"birthday" db:"birthday"`
	Phone    string `json:"phone" db:"phone"`
	Mail     string `json:"mail" db:"mail"`
	Status   int    `json:"status" db:"status"`
	Created  string `json:"created" db:"created"`
	Updated  string `json:"updated" db:"updated"`
	Deleted  string `json:"deleted" db:"deleted"`
}

func (user *User) Add() error {

	if user.Status == 0 {
		user.Status = USER_STATUS_NORMAL
	}

	if user.Created == "" {
		user.Created = util.Now()
	}

	// 密码调整为md5编码
	user.Password = util.Md5Hash(user.Password)

	return __handle.Insert(user)
}

func (user *User) Get() (bool, error) {

	exists := false
	if user.UserId == 0 {
		return exists, fmt.Errorf("get user but user_id is 0")
	}

	obj, err := __handle.Get(User{}, user.UserId)
	if err != nil {
		return exists, nil
	}

	if objUser, ok := obj.(*User); ok {
		*user = *objUser
		exists = true
	}

	return exists, nil
}

func (user *User) Update() error {

	oldUser := &User{
		UserId: user.UserId,
	}

	if exists, err := oldUser.Get(); err != nil {
		return fmt.Errorf("update user, get old user error(%v)", err)
	} else if !exists {
		return fmt.Errorf("update user, old user not exists(%v)", oldUser)
	}

	// todo 需要合并一些信息
	if user.Password != "" {
		user.Password = util.Md5Hash(user.Password)
	}

	user.Updated = util.Now()

	if _, err := __handle.Update(user); err != nil {
		return fmt.Errorf("update user(%v), error(%v)", user, err)
	}

	return nil
}

func (user *User) Delete() error {

	oldUser := &User{
		UserId: user.UserId,
	}

	if exists, err := oldUser.Get(); err != nil {
		return fmt.Errorf("delete user, get old user error(%v)", err)
	} else if !exists {
		return fmt.Errorf("delete user, old user not exists(%v)", oldUser)
	}

	if _, err := __handle.Delete(user); err != nil {
		return fmt.Errorf("delete user(%v), error(%v)", user, err)
	}

	return nil
}

func GetByAccount(account string) (user *User) {

	if account == "" {
		return nil
	}

	// 先获取user_account
	userAccount := &UserAccount{
		Account: account,
	}

	if exists, err := userAccount.GetByAccount(); err != nil || !exists {
		return nil
	}

	// 判断是否验证
	if userAccount.Verified != VERIFIED_CONFIRM {
		return nil
	}

	// 获取user
	user = &User{}
	user.UserId = userAccount.UserId
	if exists, err := user.Get(); err != nil || !exists {
		return nil
	}

	return
}
