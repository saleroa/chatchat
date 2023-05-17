package redis

//import (
//	"chatchat/model"
//	"context"
//	"encoding/json"
//	"errors"
//	"github.com/go-redis/redis/v8"
//	"log"
//)
//
//// group 的 name 应该设置为唯一
//func CreateGroup(group model.Group) (err error) {
//
//	//创建group
//	//需要先检验是否有同名的group
//	marshal, err := json.Marshal(group)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	pipe := client.TxPipeline()
//	defer pipe.Close()
//
//	result, _ := pipe.HExists(context.Background(), "groups", group.Name).Result()
//	if result {
//		return errors.New("the same name group already exits")
//		return
//	}
//	pipe.HSet(context.Background(), "groups", group.Name, marshal)
//
//	_, err = pipe.Exec(context.Background())
//	if err != nil {
//		log.Println(err)
//		pipe.Discard()
//		return err
//	}
//	return nil
//
//}
//
//func DeleteGroup(gidString string, groupname string, client *redis.Client) (err error) {
//	pipe := client.TxPipeline()
//
//	result, _ := pipe.HExists(context.Background(), "groups", groupname).Result()
//	if !result {
//		return errors.New("the group do not exits")
//		return
//	}
//
//	pipe.HDel(context.Background(), "groups", groupname)
//
//	_, err = pipe.Exec(context.Background())
//	if err != nil {
//		log.Println(err)
//		pipe.Discard()
//		return err
//	}
//
//	return nil
//}
