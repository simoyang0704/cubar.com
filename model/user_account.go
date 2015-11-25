package model

import (
	"cubar.com/lib/util"
	"fmt"
)

const (
	VERIFIED_CONFIRM  = 1
	VERIFIED_REGISTER = 0
)

type UserAccount struct {
	UserAccountId int    `json:"user_account_id" db:"user_account_id"`
	Account       string `json:"account" db:"account"`
	UserId        int    `json:"user_id" db:"user_id"`
	Verified      int    `json:"verified" db:"verified"`
	Created       string `json:"created" db:"created"`
	Deleted       string `json:"deleted" db:"deleted"`
}

type Account struct {
	User
	Account string `json:"account"`
}

type Login struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (userAccount *UserAccount) Add() error {

	userAccount.Created = util.Now()
	userAccount.Verified = VERIFIED_CONFIRM

	return __handle.Insert(userAccount)
}

func (userAccount *UserAccount) GetByAccount() (bool, error) {

	exists := false
	if userAccount.Account == "" {
		return exists, fmt.Errorf("invalid account of null")
	}

	if err := __handle.SelectOne(userAccount, "SELECT * FROM user_accounts WHERE account=?", userAccount.Account); err != nil {
		// return fmt.Errorf("get user by account(%v), error(%v)", user.Account, err)
		// 异常就是不存在
		return exists, nil
	}

	exists = true
	return exists, nil
}

func (userAccount *UserAccount) DeleteByAccount() error {

	if userAccount.Account == "" {
		return nil
	}

	if exists, err := userAccount.GetByAccount(); err != nil {
		return err
	} else if !exists {
		return nil
	}

	_, err := __handle.Delete(userAccount)

	return err
}

func (account *Account) Register() error {

	if account.Account == "" {
		return fmt.Errorf("register user but account is null")
	}

	userAccount := &UserAccount{
		Account: account.Account,
	}

	exists, _ := userAccount.GetByAccount()
	if exists {
		return fmt.Errorf("register error, account(%v) has been used", userAccount.Account)
	}

	// 添加一个user
	user := &account.User
	if err := user.Add(); err != nil {
		return err
	}

	userAccount.UserId = user.UserId
	userAccount.Verified = VERIFIED_CONFIRM

	// 添加一个account
	if err := userAccount.Add(); err != nil {
		return err
	}

	// TODO 发送邮件验证码

	return nil
}
