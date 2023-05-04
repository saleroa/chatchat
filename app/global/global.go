package global

import (
	"chatchat/app/internal/model/config"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Config  *config.Config
	Logger  *zap.Logger
	MysqlDB *sql.DB
	Rdb     *redis.Client
)
