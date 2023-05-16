package global

import (
	"chatchat/app/internal/model/config"
	"chatchat/model"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

var (
	Config  *config.Config
	Logger  *zap.Logger
	MysqlDB *sql.DB
	Rdb     *redis.Client
	//
	MapLock   sync.RWMutex
	OnlineMap = make(map[int]*model.OnLineUser)

	GReadChannel  = make(chan *model.Message)
	GWriteChannel = make(chan *model.Message)
	Upgrade       = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
