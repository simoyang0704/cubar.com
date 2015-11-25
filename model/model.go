package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"os"
)

var (
	__handle     *gorp.DbMap
	__loggedUser *User
)

func Start() {

	err := initDb("127.0.0.1", "3306", "cubar", "root", "root12")
	if err != nil {
		panic(err)
	}
}

func Stop() {

	__handle.TraceOff()
	__handle.Db.Close()
}

func initDb(host, port, dbname, user, password string) error {

	format := "%v:%v@tcp(%v:%v)/%v?charset=utf8"
	connectionStr := fmt.Sprintf(format, user, password, host, port, dbname)
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return fmt.Errorf("链接数据库出错：", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("数据库无法链接", err)
	}

	db.SetMaxOpenConns(100)
	__handle = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	__handle.AddTableWithName(User{}, "users").SetKeys(true, "UserId")
	__handle.AddTableWithName(Community{}, "communities").SetKeys(true, "CommunityId")
	__handle.AddTableWithName(UserAccount{}, "user_accounts").SetKeys(true, "UserAccountId")
	__handle.AddTableWithName(CommunityUser{}, "community_users").SetKeys(false, "CommunityUserId")

	__handle.TraceOn("[gorp]", log.New(os.Stdout, "cubar:", log.Lmicroseconds))

	return nil
}
