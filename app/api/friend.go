package api

import (
	"chatchat/app/global"
	"chatchat/dao/redis"
	"chatchat/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

func AddFriend(c *gin.Context) {
	uid, _ := c.Get("id")
	fmt.Println(uid)
	Fusername, flag := c.GetPostForm("username")
	if flag == false {
		utils.ResponseFail(c, "请输入添加好友的用户名")
		return
	}
	db := global.MysqlDB
	fid, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", Fusername), "id")
	var test int64
	err := db.QueryRow("SELECT 1 FROM friend WHERE uid = ? AND fid = ?", uid, fid).Scan(&test)
	if err != nil && err != sql.ErrNoRows {
		utils.ResponseFail(c, err.Error())
		return
	}
	if err != sql.ErrNoRows {
		utils.ResponseFail(c, "this friend has been added")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	_, err = tx.Exec("insert into `friend` (fid,uid) values (?,?)", fid, uid)
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	_, err = tx.Exec("insert into `friend` (fid,uid) values (?,?)", uid, fid)
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "add friend success")
}

func DeleteFriend(c *gin.Context) {
	uid, _ := c.Get("id")
	Fusername, flag := c.GetPostForm("username")
	if flag == false {
		utils.ResponseFail(c, "请输入添加好友的用户名")
		return
	}
	fid, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", Fusername), "id")
	db := global.MysqlDB
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	_, err = tx.Exec("delete from `friend` where fid = ? and uid = ?", fid, uid)
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	_, err = tx.Exec("delete from `friend` where fid = ? and uid = ?", uid, fid)
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}
