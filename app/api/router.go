package api

import (
	"chatchat/app/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
)

func InitRouter() error {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("span_root")
	defer span.Finish()

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(ginhttp.Middleware(opentracing.GlobalTracer()))

	r.POST("/register", register)
	r.POST("/login", login, GetOfflineMessage) // 钩子函数
	r.POST("/verificationID", SendMail)
	r.POST("/RVerificationID", RSendMail)

	r.GET("/getCookie", GetEmptyCookie)
	r.GET("/oauth2login", Oauth2Login)
	r.POST("/oauth2Register", Oauth2Register)
	r.GET("/oauth2", Oauth2)
	r.GET("/oauth2/refresh", Oauth2Refresh)
	r.GET("/oauth2/try", Oauth2Try)
	r.POST("/oauth2/pwd", Oauth2Pwd)
	r.GET("/oauth2/client", Oauth2Client)
	r.GET("/oauth2/logout", Oauth2Logout)

	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/changePassword", ChangePassword)
		UserRouter.POST("/changeNickname", ChangeNickname)
		UserRouter.POST("/changeIntroduction", ChangeIntroduction)
		UserRouter.POST("/changeAvatar", ChangeAvatar)
		UserRouter.GET("/getUser", GetUser)
	}

	GroupRouter := r.Group("/group")
	{
		GroupRouter.Use(middleware.JWTAuthMiddleware())
		GroupRouter.POST("/createGroup", CreateGroup)
		GroupRouter.POST("/joinInGroup", JoinGroup)
		GroupRouter.POST("/exitGroup", ExitGroup)
		GroupRouter.POST("/kickOut", KickOut)
		GroupRouter.DELETE("/deleteGroup", DeleteGroup)
		//查找群组没做
	}

	//FriendRouter := r.Group("/group")
	//{
	//	FriendRouter.Use(middleware.JWTAuthMiddleware())
	//
	//	//添加好友，删除好友，获取所有好友
	//	FriendRouter.POST("/addFriend", AddFriend)
	//	FriendRouter.DELETE("/deleteFriend", DeleteFriend)
	//	//查找好友没做
	//}

	ChatRouter := r.Group("/chat")
	{
		//ChatRouter.GET("/getall", GetAllYour)
		ChatRouter.GET("/conn", GetConn)
	}

	err := r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}
}
