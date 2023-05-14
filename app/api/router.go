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
	r.POST("/verificationID", SendMail)
	r.POST("/RVerificationID", RSendMail)

	r.GET("/oauth2login", Oauth2Login)
	r.POST("/oauth2Register", Oauth2Register)
	r.GET("/oauth2", Oauth2)
	r.GET("/oauth2/refresh", Oauth2Refresh)
	r.GET("/oauth2/try", Oauth2Try)
	//r.GET("/oauth2/pwd", Oauth2Pwd)
	r.GET("/oauth2/client", Oauth2Client)

	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/changePassword", ChangePassword)
		UserRouter.POST("/changeNickname", ChangeNickname)
		UserRouter.POST("/changeIntroduction", ChangeIntroduction)
		UserRouter.POST("/changeAvatar", ChangeAvatar)
		UserRouter.GET("/getUser", GetUser)
	}

	err := r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}
}
