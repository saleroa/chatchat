package mysql

import (
	"chatchat/app/global"
	"fmt"
	"log"
)

func AddOauth2User(username, oauth2Username string) (bool, string) {
	sqlStr := "insert into user_auths(username,oauth2_username) values (?,?)"
	_, err := global.MysqlDB.Exec(sqlStr, username, oauth2Username)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return false, "another error"
	}
	log.Println("insert success")
	return true, ""
}
