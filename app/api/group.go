package api

import (
	"chatchat/app/global"
	"chatchat/dao/redis"
	"chatchat/model"
	"chatchat/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

func CreateGroup(c *gin.Context) {

	id, _ := c.Get("id")
	name := c.PostForm("name")
	err := Create(name, int(id.(int64)))
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "create group success")
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
	uid, _ := c.Get("id")
	gid, _ := c.GetPostForm("gid")

	db := global.MysqlDB
	var identity int

	err := db.QueryRow("select identity from `group_members` where group_id = ? and user_id = ?", gid, uid).Scan(&identity)
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

	utils.ResponseSuccess(c, "delete group success")
}

func JoinGroup(c *gin.Context) {
	uid, _ := c.Get("id")
	gid, _ := strconv.Atoi(c.PostForm("gid"))
	db := global.MysqlDB
	_, err := db.Exec("insert into `group_members` (group_id,user_id) values (?,?)", gid, uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "join group success")
}

func ExitGroup(c *gin.Context) {

	uid, _ := c.Get("id")
	gid, _ := strconv.Atoi(c.PostForm("gid"))

	db := global.MysqlDB
	var identity int
	err := db.QueryRow("select identity from `group_members` where group_id = ? and user_id = ?", gid, uid).Scan(&identity)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	if identity == 1 {
		utils.ResponseFail(c, "manager can't leave the group")
		return
	}
	_, err = db.Exec("delete from `group_members` where group_id = ? and user_id = ?", gid, uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "exit group success")
}

func KickOut(c *gin.Context) {

	uid, _ := c.Get("id")
	gid, _ := strconv.Atoi(c.PostForm("gid"))
	kickedname, _ := c.GetPostForm("kickedname")
	kid, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", kickedname), "id")
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
	_, err = db.Exec("delete from `group_members` where group_id = ? and user_id = ?", gid, kid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func SearchGroup(c *gin.Context) {

	name, flag := c.GetPostForm("group_name")
	if flag == false {
		utils.ResponseFail(c, "请输入要查找的群聊名字")
		return
	}

	db := global.MysqlDB
	var group model.Group
	group.Name = name
	err := db.QueryRow("select gid, created_at  from `groups` where group_name = ? ", name).Scan(&group.Id, &group.Time)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"group":  group,
	})
}

func GetMembers(c *gin.Context) {
	gid, flag := c.GetPostForm("gid")
	if flag == false {
		utils.ResponseFail(c, "请输入群聊id")
		return
	}
	db := global.MysqlDB
	type member struct {
		Nickname     string
		Avatar       string
		Introduction string
	}
	var (
		members        []member
		memberusername string
		memberid       int64
	)

	rows, err := db.Query("select user_id from `group_members` where group_id = ?", gid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&memberid)
		if err != nil {
			utils.ResponseFail(c, err.Error())
			return
		}
		_ = db.QueryRow("select username from `user_bases` where id = ?", memberid).Scan(&memberusername)
		nickname, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", memberusername), "nickname")
		avatar, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", memberusername), "avatar")
		introduction, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", memberusername), "introduction")
		member := member{
			Nickname:     nickname.(string),
			Avatar:       avatar.(string),
			Introduction: introduction.(string),
		}
		members = append(members, member)
	}
	c.JSON(200, gin.H{
		"status":  200,
		"friends": members,
	})
}
