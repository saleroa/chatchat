package api

import (
	"chatchat/app/global"
	"chatchat/utils"
	"github.com/gin-gonic/gin"
)

func AddFriend(c *gin.Context) {
	uid := c.GetInt("uid")
	fid := c.GetInt("fid")
	db := global.MysqlDB
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
	utils.ResponseSuccess(c, "success")
}

func DeleteFriend(c *gin.Context) {
	uid := c.GetInt("uid")
	fid := c.GetInt("fid")
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
