package global

import (
	"chatchat/app/internal/model/config"
	"chatchat/model"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
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
			//authHeader := r.Header.Get("Authorization")
			//if authHeader == "" {
			//	return false
			//}
			//// 按空格分割
			//parts := strings.SplitN(authHeader, " ", 2)
			//if !(len(parts) == 2 && parts[0] == "Bearer") {
			//
			//	return false
			//}
			//// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
			//mc, err := ParseToken(parts[1])
			//if err != nil {
			//	return false
			//}
			//// 将当前请求的username信息保存到请求的上下文c上
			//context.WithValue(r.Context(), "id", mc.ID)
			//context.WithValue(r.Context(), "username", mc.Username)
			return true
		},
	}
)
var Secret = []byte("WZY")

func ParseToken(tokenString string) (*model.MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
