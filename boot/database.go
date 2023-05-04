package boot

import (
	"chatchat/app/global"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
)

func MysqlDBSetup() {
	config := global.Config.Database.Mysql

	db, err := sql.Open("mysql", config.GetDsn())
	if err != nil {
		global.Logger.Fatal("initialize database failed", zap.Error(err))
	}
	//defer func(db *sql.DB) {
	//	err := db.Close()
	//	if err != nil {
	//
	//	}
	//}(db)关闭数据库

	db.SetConnMaxLifetime(global.Config.Database.Mysql.GetConnMaxIDleTime())
	db.SetConnMaxIdleTime(global.Config.Database.Mysql.GetConnMaxIDleTime())
	db.SetMaxIdleConns(global.Config.Database.Mysql.MaxOpenConns)
	db.SetMaxOpenConns(global.Config.Database.Mysql.MaxIdleConns)
	err = db.Ping()
	if err != nil {
		global.Logger.Fatal("connected failed", zap.Error(err))
	}
	global.MysqlDB = db

	global.Logger.Info("initialize database success")
}

func RedisSetup() {
	config := global.Config.Database.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Addr, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.Db,
		PoolSize: config.PoolSize,
	})

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()//关闭redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Fatal("connect to redis failed", zap.Error(err))
	}
	global.Rdb = rdb

	global.Logger.Info("initialize redis success")
}
