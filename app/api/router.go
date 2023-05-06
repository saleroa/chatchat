package api

import (
	"chatchat/app/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() error {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/register", register)
	r.POST("/login", login)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/:uid/changePassword", ChangePassword)
	}
	err := r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}
}
