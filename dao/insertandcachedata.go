package dao

import (
	"chatchat/dao/mysql"
	rdb "chatchat/dao/redis"
	"chatchat/model"

	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
)

func InsertAndCacheData(db *sql.DB, cli *redis.Client, message model.Message) error {
	//需要对单发和群聊的消息做区分
	key := uuid.New().String()
	// 插入新数据到 MySQL 中
	err := mysql.InsertIntoMysql(db, message, key)
	if err != nil {
		log.Println(err)
		return err
	}
	err = rdb.InsertIntoRedis(cli, message, key)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
