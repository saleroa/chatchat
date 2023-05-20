package api

import (
	"chatchat/app/global"
	"chatchat/model"
	"chatchat/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

func CreateGroup(c *gin.Context) {

	id, _ := strconv.Atoi(c.PostForm("uid"))
	name := c.PostForm("name")
	err := Create(name, id)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func Create(name string, id int) (err error) {
	db := global.MysqlDB
	//开启事务
	tx, err := db.Begin()

	if err != nil {
		panic(err.Error())
	}
	_, err = tx.Exec("insert into `groups` (group_name,created_at) values (?,?)", name, time.Now())
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	var gid int
	err = tx.QueryRow("SELECT gid FROM `groups` WHERE group_name= ?  LIMIT 1", name).Scan(&gid)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	//群主为 1 ,成员为默认 0
	_, err = tx.Exec("insert into `group_members` (group_id,user_id,identity) values (?,?,?)", gid, id, 1)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	tx.Commit()

	return nil
}

func DeleteGroup(c *gin.Context) {
	uid := c.GetInt("uid")
	gid := c.GetInt("gid")

	db := global.MysqlDB
	var identity int

	err := db.QueryRow("select identidy from `groups` where gid = ? and uid = ?", gid, uid).Scan(&identity)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	if identity != 1 {
		utils.ResponseFail(c, "you are not manager")
		return
	}

	_, err = db.Exec("delete from `groups` where gid = ? ", gid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}

	utils.ResponseSuccess(c, "success")
}

func JoinGroup(c *gin.Context) {
	uid, _ := strconv.Atoi(c.PostForm("uid"))
	gid, _ := strconv.Atoi(c.PostForm("gid"))
	db := global.MysqlDB
	_, err := db.Exec("insert into `group_members` (group_id,user_id) values (?,?)", gid, uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func ExitGroup(c *gin.Context) {

	uid, _ := strconv.Atoi(c.PostForm("uid"))
	gid, _ := strconv.Atoi(c.PostForm("gid"))

	db := global.MysqlDB
	_, err := db.Exec("delete from `group_members` where group_id = ? and user_id = ?", gid, uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func KickOut(c *gin.Context) {

	uid, _ := strconv.Atoi(c.PostForm("uid"))
	gid, _ := strconv.Atoi(c.PostForm("gid"))
	kickedid, _ := strconv.Atoi(c.PostForm("kickedid"))

	db := global.MysqlDB
	var identity int
	err := db.QueryRow("select  identity from `group_members` where group_id = ? and user_id = ? ", gid, uid).Scan(&identity)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseFail(c, "no such user")
			return
			//没有该用户存在
		}
		utils.ResponseFail(c, err.Error())
		return
		//查询出错
	}
	if identity != 1 {
		utils.ResponseFail(c, err.Error())
		return
	}
	_, err = db.Exec("delete from `group_members` where group_id = ? and user_id = ?", gid, kickedid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func SearchGroup(c *gin.Context) {

	name := c.GetString("groupname")

	db := global.MysqlDB
	var group model.Group
	err := db.QueryRow("select gid, created_at  from `groups` where group_name = ? ", name).Scan(&group.Id, &group.Time)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}

	c.JSON(200, group)
}
