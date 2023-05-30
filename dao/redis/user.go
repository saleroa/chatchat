package redis

import (
	"chatchat/app/global"
	"chatchat/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"time"
)

var user *model.User

func HGet(ctx context.Context, key, field string) (interface{}, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis:user:HGet")
	defer span.Finish()
	GetKey := global.Rdb.HGet(ctx, key, field)
	if GetKey.Err() != nil && GetKey.Err() != redis.Nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("redis HGet value from %s in %s failed", key, field))
		global.Logger.Error(GetKey.Err().Error())
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
	span, _ := opentracing.StartSpanFromContext(ctx, "redis:user:HSet")
	defer span.Finish()
	SetKV := global.Rdb.HSet(ctx, key, value)
	if SetKV.Err() != nil && SetKV.Err() != redis.Nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("redis set %s to %s failed", value, key))
		global.Logger.Error(SetKV.Err().Error())
	}
	return SetKV.Err()
}

func Set(ctx context.Context, key string, value interface{}, time time.Duration) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis:user:Set")
	defer span.Finish()
	SetKV := global.Rdb.Set(ctx, key, value, time)
	if SetKV.Err() != nil && SetKV.Err() != redis.Nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("redis set value from %s to %s failed", value, key))
		global.Logger.Error(SetKV.Err().Error())
	}
	return SetKV.Err()
}

func Get(ctx context.Context, key string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis:user:Get")
	defer span.Finish()
	GetKey := global.Rdb.Get(ctx, key)
	if GetKey.Err() != nil && GetKey.Err() != redis.Nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("redis get value from %s failed", key))
		global.Logger.Error(GetKey.Err().Error())
		return "", GetKey.Err()
	}
	return GetKey.Val(), nil
}

func ZSetUserID(ctx context.Context, username string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis:user:ZSetUserID")
	defer span.Finish()
	ID := global.Rdb.ZCard(ctx, "userID")

	var z = &redis.Z{
		Score:  float64(ID.Val() + 1),
		Member: username,
	}

	SetKV := global.Rdb.ZAdd(ctx, "userID", z)
	if SetKV.Err() != nil && SetKV.Err() != redis.Nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("redis ZSetUserID from %s failed", username))
		global.Logger.Error(SetKV.Err().Error())
	}
	return SetKV.Err()
}
func ZGet(ctx context.Context, key string, member string) (error, interface{}) {
	span, _ := opentracing.StartSpanFromContext(ctx, "redis:user:ZGet")
	defer span.Finish()
	GetKey := global.Rdb.ZScore(ctx, key, member)
	if GetKey.Err() != nil && GetKey.Err() != redis.Nil {
		span.SetTag("error", true)
		span.SetTag("error_info", fmt.Sprintf("redis get value from %s failed", key))
		global.Logger.Error(GetKey.Err().Error())
		return GetKey.Err(), ""
	}
	return nil, GetKey.Val()
}
