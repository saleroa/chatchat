package redis

import (
	"chatchat/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func InsertIntoRedis(cli *redis.Client, message model.Message, key string) error {
	// 将数据缓存到 Redis 中
	id := fmt.Sprintf("%sto%s", message.FromId, message.TargetId)

	err := cli.HSet(context.Background(), key, "time", message.Time, "content", message.Content, "id", id, "sendtype", message.SendType).Err()
	if err != nil {
		log.Println()
		return err
	}
	err = cli.Expire(context.Background(), key, 24*7*time.Hour).Err()
	if err != nil {
		log.Println()
		return err
	}
	return nil
}

// 插入离线消息
func SaveOfflineMessage(message model.Message, cli *redis.Client) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = cli.LPush(context.Background(), string(message.TargetId), msg).Err()
	if err != nil {
		return err
	}

	return nil
}
