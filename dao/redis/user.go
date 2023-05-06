package redis

import (
	"chatchat/app/global"
	"chatchat/model"
	"context"
	"github.com/go-redis/redis/v8"
)

var user *model.User

func HGet(ctx context.Context, key, field string) (interface{}, error) {
	GetKey := global.Rdb.HGet(ctx, key, field)
	if GetKey.Err() != nil {
		return "", GetKey.Err()
	}
	return GetKey.Val(), nil
}

//func Check(ctx context.Context, field, context string) (bool, error) {
//	var cursor uint64
//	for {
//		data, cursor, err := global.Rdb.GeoSearch(ctx, cursor, context, 20).Result()
//		if err != nil {
//			return true, err
//		}
//		for _, k := range data {
//			value, _ := global.Rdb.Get(ctx, k).Result()
//			json.Unmarshal([]byte(value), user)
//		}
//		if cursor == 0 {
//			break
//		}
//	}
//	return false, nil
//}

func HSet(ctx context.Context, key string, value ...interface{}) error {
	SetKV := global.Rdb.HSet(ctx, key, value)
	return SetKV.Err()
}

func ZSetUserID(ctx context.Context, username string) error {
	ID := global.Rdb.ZCard(ctx, "userID")

	var z = &redis.Z{
		Score:  float64(ID.Val() + 1),
		Member: username,
	}

	SetKV := global.Rdb.ZAdd(ctx, "userID", z)
	return SetKV.Err()
}
func ZGet(ctx context.Context, key string, member string) (error, interface{}) {
	GetKey := global.Rdb.ZScore(ctx, key, member)
	if GetKey.Err() != nil {
		return GetKey.Err(), ""
	}
	return nil, GetKey.Val()
}
