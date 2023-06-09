package api

import (
	"chatchat/app/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
)

func InitRouter() error {
	tracer, _, err := middleware.InitJaeger("chatchat")
	if err != nil {
		//global.Logger.Fatal("initialize jaeger failed", zap.Error(err))
	}
	opentracing.SetGlobalTracer(tracer)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(ginhttp.Middleware(tracer))

	r.POST("/register", register)
	r.POST("/login", login)
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
		GroupRouter.POST("/searchGroup", SearchGroup)
		GroupRouter.POST("/getMembers", GetMembers)

	}

	FriendRouter := r.Group("/friend")
	{
		FriendRouter.Use(middleware.JWTAuthMiddleware())
		FriendRouter.POST("/addFriend", AddFriend)
		FriendRouter.DELETE("/deleteFriend", DeleteFriend)
	}

	ChatRouter := r.Group("/chat")
	{
		ChatRouter.Use(middleware.JWTAuthMiddleware())
		ChatRouter.GET("/getGroups", GetGroups)
		ChatRouter.GET("/getFriends", GetFriends)
		ChatRouter.GET("/getFriendMessage", GetFriendMessage)
		ChatRouter.GET("/getGroupMessage", GetGroupMessage)
		ChatRouter.GET("/getOffMsg", GetOfflineMessage)
	}
	r.GET("/chat/conn", GetConn)

	err = r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}
}
