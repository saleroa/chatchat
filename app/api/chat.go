package api

import (
	"chatchat/app/global"
	"chatchat/dao"
	"chatchat/dao/redis"
	"chatchat/model"
	"chatchat/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func GetConn(c *gin.Context) {
	coon, err := global.Upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		log.Println(err)
		return
	}
	id := c.GetInt("id")
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

					err := redis.SaveOfflineMessage(*message, cli)
					if err != nil {
						return
					}
				}
				//用户在线
				onLineUser.ReadChannel <- *message
				//缓存
				err := dao.InsertAndCacheData(db, cli, *message)
				if err != nil {
					return
				}
			}
			//群发逻辑

			result, err := cli.HGetAll(context.Background(), fmt.Sprintf("users:%s", string(message.TargetId))).Result()
			if err != nil {
				log.Println(err)
				return
			}
			for _, uidstring := range result {
				uid, _ := strconv.Atoi(uidstring)
				//群聊用户在线
				if global.OnlineMap[uid] != nil {

					global.OnlineMap[uid].ReadChannel <- *message
					err := dao.InsertAndCacheData(db, cli, *message)
					if err != nil {
						return
					}
				}
				//群聊用户不在线
				err := redis.SaveOfflineMessage(*message, cli)
				if err != nil {
					return
				}

			}
		}
	}
}

// 上线后第一件事，读取离线消息
func GetOfflineMessage(c *gin.Context) {
	id := c.GetInt("id")
	cli := global.Rdb
	len, err := cli.LLen(context.Background(), string(id)).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
	}
	result, err := cli.LRange(context.Background(), "", 0, len-1).Result()
	if err != nil {
		utils.ResponseFail(c, err.Error())
	}
	defer cli.Del(context.Background(), string(id))

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

//func GetAllYours(c *gin.Context) {
//	uid := c.GetInt("uid")
//
//	//返回好友的id
//	//返回群组的id
//
//}
