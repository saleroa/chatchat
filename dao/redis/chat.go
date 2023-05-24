package redis

import (
	"chatchat/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

func InsertIntoRedis(cli *redis.Client, message model.Message, key string) error {
	// 将数据缓存到 Redis 中
	msg, _ := json.Marshal(message)
	if message.SendType == 1 {
		err := cli.RPush(context.Background(), fmt.Sprintf("friend:%s", key), msg).Err()
		if err != nil {
			return err
		}
	} else if message.SendType == 2 {
		err := cli.RPush(context.Background(), fmt.Sprintf("group:%s", key), msg).Err()
		if err != nil {
			return err
		}
	}
	err := cli.Expire(context.Background(), key, 24*7*time.Hour).Err()
	if err != nil {
		log.Println()
		return err
	}
	return nil
}

// 插入离线消息
func SaveOfflineMessage(message model.Message, cli *redis.Client, uid int64) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if message.SendType == 1 {
		err = cli.RPush(context.Background(), strconv.Itoa(message.TargetId), msg).Err()
		if err != nil {
			return err
		}
	} else if message.SendType == 2 {
		err = cli.RPush(context.Background(), strconv.Itoa(int(uid)), msg).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
