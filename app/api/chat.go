package api

import (
	"chatchat/app/global"
	"chatchat/dao"
	"chatchat/dao/redis"
	"chatchat/model"
	"chatchat/utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

func GetConn(c *gin.Context) {
	coon, err := global.Upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		log.Println(err)
		return
	}
	ID, _ := c.GetQuery("id")
	id, _ := strconv.Atoi(ID)
	client := &model.OnLineUser{
		UserId:       id,
		Coon:         coon,
		ReadChannel:  make(chan model.Message),
		WriteChannel: make(chan model.Message),
	}
	global.MapLock.RLock()
	global.OnlineMap[id] = client
	global.MapLock.RUnlock()
	//携程处理，为了数据的可靠需要加锁

	go client.Read()
	go client.Write()
	Run(client)
}

// 读取websocket的消息，发送到管道
func Run(client *model.OnLineUser) {
	go GRead()
	go GWrite()

	for {
		//连接websocket后，消息会发送到这里
		//消息需要json格式
		_, msg, err := client.Coon.ReadMessage()
		if err != nil {
			log.Println("coon close")
			global.MapLock.RLock()
			delete(global.OnlineMap, client.UserId)
			global.MapLock.RUnlock()
			close(client.ReadChannel)
			close(client.WriteChannel)
			return
		}
		message := &model.Message{}
		err = json.Unmarshal(msg, message)
		if err != nil {
			log.Println(err)
			return
		}
		message.Time = time.Now()
		fmt.Println(message.Time)
		global.GReadChannel <- message

	}
}

func GRead() {
	for {
		select {
		case msg, ok := <-global.GReadChannel:
			if !ok {
				log.Println("ReadService err...")
				return
			}
			global.GWriteChannel <- msg
		}
	}
}

// 这里对消息的类型判断会有些问题
// 如果是 string 的消息的话ok,但如果是 图片的url就不该存进数据库

func GWrite() {
	db := global.MysqlDB
	cli := global.Rdb
	for {
		select {
		case message, ok := <-global.GWriteChannel:

			if !ok {
				log.Println("ReadService err...")
				return
			}
			//私发
			if message.SendType == 1 {
				onLineUser := global.OnlineMap[message.TargetId]
				if onLineUser == nil {
					//用户不在线

					err := redis.SaveOfflineMessage(*message, cli, 0)
					err1 := dao.InsertAndCacheData(db, cli, *message)
					if err != nil || err1 != nil {
						return
					}
				} else {
					//用户在线
					onLineUser.ReadChannel <- *message
					//缓存
					err := dao.InsertAndCacheData(db, cli, *message)
					if err != nil {
						return
					}
				}
			} else if message.SendType == 2 {

				//群发逻辑

				rows, err := db.Query("select user_id from `group_members` where group_id = ?", message.TargetId)
				if err != nil {
					log.Println(err)
					return
				}
				for rows.Next() {
					var uid int
					rows.Scan(&uid)

					//群聊用户在线
					if global.OnlineMap[uid] == nil {
						//群聊用户不在线
						err := redis.SaveOfflineMessage(*message, cli, int64(uid))
						if err != nil {
							return
						}
					} else {
						if uid == message.FromId {
							err := dao.InsertAndCacheData(db, cli, *message)
							if err != nil {
								return
							}
						}
						global.OnlineMap[uid].ReadChannel <- *message
					}
				}
			}
		}
	}
}

// 上线后第一件事，读取离线消息
func GetOfflineMessage(c *gin.Context) {
	v, _ := c.Get("id")
	ID, _ := v.(int64)
	id := int(ID)
	//username, _ := c.Get("username")
	//fmt.Println(username)
	//fmt.Println(ID)
	//fmt.Println(id)
	cli := global.Rdb
	len, err := cli.LLen(c.Request.Context(), strconv.Itoa(id)).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	result, err := cli.LRange(c.Request.Context(), strconv.Itoa(id), 0, len-1).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer cli.Del(c.Request.Context(), strconv.Itoa(id))
	for _, res := range result {
		var s model.Message
		if err := json.Unmarshal([]byte(res), &s); err != nil {
			utils.ResponseFail(c, err.Error())
		}
		//读取到离线消息了，发送给targetid
		global.OnlineMap[id].ReadChannel <- s

		continue
	}
}

type Group struct {
	Gid  int
	Name string
}

func GetGroups(c *gin.Context) {
	var (
		groupid int
		group   Group
		groups  []Group
	)
	uid, _ := c.Get("id")
	db := global.MysqlDB

	rows, err := db.Query("select group_id from `group_members` where user_id = ?", uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&groupid)
		if err != nil {
			utils.ResponseFail(c, err.Error())
			return
		}
		_ = db.QueryRow("select group_name from `groups` where gid = ?", groupid).Scan(&group.Name)
		group.Gid = groupid
		groups = append(groups, group)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"groups": groups,
	})
}
func GetFriends(c *gin.Context) {
	type friend struct {
		ID           int64
		Nickname     string
		Avatar       string
		Introduction string
	}
	var (
		friends        []friend
		friendUsername string
		friendid       int64
	)
	uid, _ := c.Get("id")

	db := global.MysqlDB

	rows, err := db.Query("select fid from `friend` where uid = ?", uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&friendid)
		if err != nil {
			utils.ResponseFail(c, err.Error())
			return
		}
		_ = db.QueryRow("select username from `user_bases` where id = ?", friendid).Scan(&friendUsername)
		nickname, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", friendUsername), "nickname")
		avatar, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", friendUsername), "avatar")
		introduction, _ := redis.HGet(c.Request.Context(), fmt.Sprintf("user:%s", friendUsername), "introduction")
		friend := friend{
			ID:           friendid,
			Nickname:     nickname.(string),
			Avatar:       avatar.(string),
			Introduction: introduction.(string),
		}
		friends = append(friends, friend)
	}

	c.JSON(200, gin.H{
		"status":  200,
		"friends": friends,
	})
}
func GetAll(c *gin.Context) {
	var (
		friends  []string
		friend   string
		friendid int64
		groupid  int
		group    string
		groups   []string
	)
	uid, _ := c.GetPostForm("uid")

	db := global.MysqlDB

	rows, err := db.Query("select fid from `friend` where uid = ?", uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&friendid)
		if err != nil {
			utils.ResponseFail(c, err.Error())
			return
		}
		_ = db.QueryRow("select nickname from `user_bases` where id = ?", friendid).Scan(&friend)
		friends = append(friends, friend)
	}

	rows, err = db.Query("select group_id from `group_members` where user_id = ?", uid)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&groupid)
		if err != nil {
			utils.ResponseFail(c, err.Error())
			return
		}
		_ = db.QueryRow("select group_name from `groups` where gid = ?", groupid).Scan(&group)
		fmt.Println(group)
		groups = append(groups, group)
	}

	c.JSON(200, gin.H{
		"status":  200,
		"friends": friends,
		"groups":  groups,
	})
}
func GetFriendMessage(c *gin.Context) {
	ID, _ := c.Get("id")
	id := ID.(int64)
	value2, _ := c.GetQuery("toid")
	Size, _ := c.GetQuery("size")
	size, _ := strconv.ParseInt(Size, 10, 64)
	Offset, _ := c.GetQuery("offset")
	offset, _ := strconv.ParseInt(Offset, 10, 64)
	toid, _ := strconv.ParseInt(value2, 10, 64)
	if id > toid {
		id, toid = toid, id
	}
	key := fmt.Sprintf("friend:%dto%d", id, toid)
	cli := global.Rdb
	lenth, err := cli.LLen(context.Background(), key).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	//fmt.Println(len, offset, size)
	result, err := cli.LRange(context.Background(), key, lenth-(offset+1)*size, lenth-offset*size-1).Result()

	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	var msgs []model.Message
	for _, res := range result {
		var msg model.Message
		if err := json.Unmarshal([]byte(res), &msg); err != nil {
			utils.ResponseFail(c, err.Error())
		}
		msgs = append(msgs, msg)
		continue
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": msgs,
	})
}
func GetGroupMessage(c *gin.Context) {
	ID, _ := c.Get("id")
	id := ID.(int64)
	value2, _ := c.GetQuery("gid")
	Size, _ := c.GetQuery("size")
	size, _ := strconv.ParseInt(Size, 10, 64)
	Offset, _ := c.GetQuery("offset")
	offset, _ := strconv.ParseInt(Offset, 10, 64)
	gid, _ := strconv.ParseInt(value2, 10, 64)
	key := fmt.Sprintf("group:%d", gid)
	cli := global.Rdb
	db := global.MysqlDB
	identity := -1
	err := db.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?", gid, id).Scan(&identity)
	if err != nil && err != sql.ErrNoRows {
		utils.ResponseFail(c, err.Error())
		return
	}
	if identity == -1 {
		utils.ResponseFail(c, "you do not have the access to search the message")
		return
	}
	lenth, err := cli.LLen(context.Background(), key).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	result, err := cli.LRange(context.Background(), key, lenth-(offset+1)*size, lenth-offset*size-1).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	var msgs []model.Message
	for _, res := range result {
		var msg model.Message
		if err := json.Unmarshal([]byte(res), &msg); err != nil {
			utils.ResponseFail(c, err.Error())
		}
		msgs = append(msgs, msg)
		continue
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": msgs,
	})
}
