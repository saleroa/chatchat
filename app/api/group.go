package api

import (
	"chatchat/app/global"
	"chatchat/dao/redis"
	"chatchat/model"
	"chatchat/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateGroup(c *gin.Context) {
	uid := c.GetInt("uid")
	name := c.GetString("name")

	err := redis.CreateGroup(model.Group{
		Name:     name,
		Time:     time.Now(),
		MangerID: uid,
	})
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func DeleteGroup(c *gin.Context) {
	uid := c.GetInt("uid")
	gidString := c.GetString("gid")

	client := global.Rdb

	managerid, err := client.HGet(context.Background(), fmt.Sprintf("groups:%s", gidString), "managerid").Int()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	if managerid != uid {
		utils.ResponseFail(c, "you are not manager")
		return
	}

	err = redis.DeleteGroup(gidString, client)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}

	utils.ResponseSuccess(c, "success")
}

func JoinGroup(c *gin.Context) {
	uid := c.GetInt("uid")
	gidString := c.GetString("gid")
	uidString := string(uid)

	client := global.Rdb
	bool, err := client.HSetNX(context.Background(), fmt.Sprintf("users:%s", gidString), uidString, uid).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	if !bool {
		utils.ResponseFail(c, "failed to insert user into group")
		return
	}
	utils.ResponseSuccess(c, "success")
}

func ExitGroup(c *gin.Context) {
	uidString := c.GetString("uid")
	gidString := c.GetString("gid")

	client := global.Rdb

	err := client.HDel(context.Background(), fmt.Sprintf("users:%s", gidString), uidString).Err()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}

func KickOut(c *gin.Context) {

	uid := c.GetInt("uid")
	gidString := c.GetString("gid")
	kickedidString := c.GetString("kickedid")

	client := global.Rdb
	managerid, err := client.HGet(context.Background(), fmt.Sprintf("users:%s", gidString), "managerid").Int()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	if managerid != uid {
		utils.ResponseFail(c, "you are not manager")
		return
	}
	err = client.HDel(context.Background(), gidString, kickedidString).Err()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, "success")
}
