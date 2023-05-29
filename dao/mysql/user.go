package mysql

import (
	"chatchat/app/global"
	"chatchat/model"
	"chatchat/utils"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"log"
	"time"
)

var u model.User

func AddUser(ctx context.Context, username, password, nickname string, ID int64) (bool, string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:addUser")
	defer span.Finish()
	EncryptPassword, err1 := utils.EncryptPassword(password) //加密密码
	if err1 != nil {
		return false, "encrypt password failed"
	}
	sqlStr := "insert into user_bases(id,username,password,nickname,vip,create_time,update_time) values (?,?,?,?,?,?,?)"
	_, err := global.MysqlDB.Exec(sqlStr, ID, username, EncryptPassword, nickname, 0, time.Now(), time.Now())
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", err)
		fmt.Printf("insert failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return false, err.Error()
	}
	log.Println("insert success")
	return true, ""
}
func SelectID(ctx context.Context, username string) int64 {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:selectUser")
	defer span.Finish()
	sqlStr := "select id, username, password from user_bases where username=?"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	u.Username = username
	err := global.MysqlDB.QueryRow(sqlStr, username).Scan(&u.ID, u.Username, &u.Password)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", err)
		fmt.Printf("scan failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return 0
	}
	return u.ID
}
func ChangePassword(ctx context.Context, st model.User) (bool, string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:changePassword")
	defer span.Finish()
	sqlStr := "update user_bases set password=? where username=?"
	EncryptPassword, err1 := utils.EncryptPassword(st.Password) //加密密码
	if err1 != nil {
		return false, "encrypt password failed"
	}
	_, err := global.MysqlDB.Exec(sqlStr, EncryptPassword, st.Username)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", err)
		fmt.Printf("update failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func ChangeNickname(ctx context.Context, st model.User) (bool, string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:changeNickname")
	defer span.Finish()
	sqlStr := "update user_bases set nickname=? where username=?"
	_, err := global.MysqlDB.Exec(sqlStr, st.Nickname, st.Username)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", err.Error())
		fmt.Printf("update failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func ChangeIntroduction(ctx context.Context, st model.User) (bool, string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:changeIntroduction")
	defer span.Finish()
	sqlStr := "update user_bases set introduction=? where username=?"
	_, err := global.MysqlDB.Exec(sqlStr, st.Introduction, st.Username)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("change introduction failed,error is %s", err))
		fmt.Printf("update failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}
func ChangeAvatar(ctx context.Context, st model.User) (bool, string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:changeAvatar")
	defer span.Finish()
	sqlStr := "update user_bases set avatar=? where username=?"
	_, err := global.MysqlDB.Exec(sqlStr, st.Avatar, st.Username)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("change avatar failed,err is %s", err))
		fmt.Printf("update failed, err:%v\n", err)
		global.Logger.Error(err.Error())
		return false, "update failed"
	}
	log.Println("update success")
	return true, ""
}

func FindID(ctx context.Context) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql:user:findID")
	defer span.Finish()
	sqlStr := "select id from user_bases where id >=?"
	rows, err := global.MysqlDB.Query(sqlStr, 0)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("find ID failed,err is %s", err))
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
