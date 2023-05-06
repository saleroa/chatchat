package mysql

import (
	"chatchat/app/global"
	"chatchat/model"
	"chatchat/utils"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var u model.User

func AddUser(username, password, nickname string, ID int64) (bool, string) {
	EncryptPassword, err1 := utils.EncryptPassword(password) //加密密码
	if err1 != nil {
		return false, "encrypt password failed"
	}
	sqlStr := "insert into user_bases(id,username,password,nickname,vip,create_time,update_time) values (?,?,?,?,?,?,?)"
	_, err := global.MysqlDB.Exec(sqlStr, ID, username, EncryptPassword, nickname, 0, time.Now(), time.Now())
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return false, "another error"
	}
	log.Println("insert success")
	return true, ""
}

func SelectUser(username string) bool {
	//boot.MysqlDBSetup()
	u.Username = username
	sqlStr := "select  username from user_bases where username=?"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := global.MysqlDB.QueryRow(sqlStr, username).Scan(&u.Username)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	return true
}

func SelectPassword(username string) string {
	sqlStr := "select id, username, password from user_bases where username=?"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	u.Username = username
	err := global.MysqlDB.QueryRow(sqlStr, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return ""
	}
	return u.Password
}
func SelectID(username string) int64 {
	sqlStr := "select id, username, password from user_bases where username=?"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	u.Username = username
	err := global.MysqlDB.QueryRow(sqlStr, username).Scan(&u.ID, u.Username, &u.Password)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return 0
	}
	return u.ID
}
func ChangePassword(st model.User) (bool, string) {
	sqlStr := "update user_bases set password=? where username=?"
	EncryptPassword, err1 := utils.EncryptPassword(st.Password) //加密密码
	if err1 != nil {
		return false, "encrypt password failed"
	}
	_, err := global.MysqlDB.Exec(sqlStr, EncryptPassword, st.Username)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func ChangeNickname(st model.User) (bool, string) {
	sqlStr := "update user_bases set nickname=? where username=?"
	_, err := global.MysqlDB.Exec(sqlStr, st.Nickname, st.Username)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func ChangeIntroduction(st model.User) (bool, string) {
	sqlStr := "update user_bases set introduction=? where username=?"
	_, err := global.MysqlDB.Exec(sqlStr, st.Introduction, st.Username)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func ChangeAvatar(st model.User) (bool, string) {
	sqlStr := "update user_bases set avatar=? where username=?"
	_, err := global.MysqlDB.Exec(sqlStr, st.Avatar, st.Username)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func FindNickname(nickname string) bool {
	sqlStr := "select id,nickname from user_bases where nickname=?"
	u.Nickname = nickname
	u.ID = 0
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := global.MysqlDB.QueryRow(sqlStr, nickname).Scan(&u.ID, &u.Nickname)
	if err == sql.ErrNoRows {
		//无事发生
	} else if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	//if err != nil || err != errors.New("sql: no rows in result set") {
	//	fmt.Printf("scan failed, err:%v\n", err)
	//	return false
	//}
	if u.ID == 0 {
		return true
	} else {
		return false
	}
}
func FindID() int {
	sqlStr := "select id from user_bases where id >=?"
	rows, err := global.MysqlDB.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return 0
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	i := 1
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return 0
		}
		i++
	}
	return i
}
