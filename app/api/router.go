package api

import (
	"chatchat/app/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() error {
	r := gin.Default()
	r.Use(middleware.CORS())
	err := r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}
}
