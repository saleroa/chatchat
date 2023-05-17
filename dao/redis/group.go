package redis

import (
	"chatchat/app/global"
	"chatchat/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

// group 的 name 应该设置为唯一
func CreateGroup(group model.Group) (err error) {

	client := global.Rdb
	//创建group
	pipe := client.TxPipeline()

	Id, err := pipe.Incr(context.Background(), "group:id").Result()
	if err != nil {
		log.Println(err)
		return err
	}
	gidString := string(Id)
	_, err = client.HSet(context.Background(), fmt.Sprintf("groups:%s", gidString), "name", group.Name, "time", group.Time, "managerid", group.MangerID).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func DeleteGroup(gidString string, client *redis.Client) (err error) {
	pipe := client.TxPipeline()
	err = pipe.Del(context.Background(), fmt.Sprintf("users:%s", gidString)).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	err = pipe.Del(context.Background(), fmt.Sprintf("groups:%s", gidString)).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
