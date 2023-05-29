package dao

import (
	"chatchat/app/global"
	"chatchat/dao/mysql"
	rdb "chatchat/dao/redis"
	"chatchat/model"
	"fmt"

	"database/sql"
	"github.com/go-redis/redis/v8"
)

func InsertAndCacheData(db *sql.DB, cli *redis.Client, message model.Message) error {
	//需要对单发和群聊的消息做区分
	fromid := message.FromId
	targetid := message.TargetId
	if message.SendType == 1 {
		if fromid > targetid {
			fromid, targetid = targetid, fromid
		}
		key := fmt.Sprintf("%dto%d", fromid, targetid)
		// 插入新数据到 MySQL 中
		err := mysql.InsertIntoMysql(db, message, key)
		if err != nil {
			global.Logger.Error(err.Error())
			return err
		}
		err = rdb.InsertIntoRedis(cli, message, key)
		if err != nil {
			global.Logger.Error(err.Error())
			return err
		}
	} else if message.SendType == 2 {
		key := fmt.Sprintf("%d", message.TargetId)
		// 插入新数据到 MySQL 中
		err := mysql.InsertIntoMysql(db, message, key)
		if err != nil {
			global.Logger.Error(err.Error())
			return err
		}
		err = rdb.InsertIntoRedis(cli, message, key)
		if err != nil {
			global.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
