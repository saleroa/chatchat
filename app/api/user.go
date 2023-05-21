package api

import (
	"chatchat/app/api/middleware"
	"chatchat/app/global"
	"chatchat/dao"
	"chatchat/dao/mysql"
	"chatchat/dao/redis"
	"chatchat/model"
	"chatchat/utils"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func register(c *gin.Context) {
	e, b := sentinel.Entry("chatchat_user", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		// 请求被拒绝，在此处进行处理
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  429,
			"message": "too many request",
		})
	} else {
		var user *model.User
		if err := c.ShouldBind(&user); err != nil {
			fmt.Println(err)
			utils.ResponseFail(c, "verification failed")
			return
		}
		// 传入用户名和密码
		username := user.Username
		password := user.Password
		nickname := user.Nickname
		mailID := user.MailID
		EncryptPassword, err1 := utils.EncryptPassword(password) //加密密码
		if err1 != nil {
			utils.ResponseFail(c, "encrypt password failed")
			return
		}
		flag2, msg := dao.AddUserCheck(username, password, nickname)
		if !flag2 {
			utils.ResponseFail(c, msg)
			return
		}
		user.Password = EncryptPassword
		//_ = redis.Set(c, fmt.Sprintf("%s:vip", username), "0", 0)
		user.Nickname = nickname
		user.Introduction = "这个人很懒，什么都没留下~"
		user.Avatar = "http://test.violapioggia.cn/chatchatUsers/empty_avatar.png"
		uid, err := redis.Get(c.Request.Context(), fmt.Sprintf("Rmail:%s", username))
		if err != nil {
			utils.ResponseFail(c, "verification code has expired")
			return
		}
		if uid == "" || uid != mailID {
			utils.ResponseFail(c, "wrong mailID")
			return
		}

		user.ID = global.Rdb.ZCard(c.Request.Context(), "userID").Val() + 1
		flag1, msg := mysql.AddUser(c.Request.Context(), username, password, nickname, user.ID) //写入数据库
		if flag1 {
			err := redis.ZSetUserID(c.Request.Context(), username)
			if err != nil {
				utils.ResponseFail(c, err.Error())
				return
			}
			err = redis.HSet(c.Request.Context(), fmt.Sprintf("user:%s", username), "id", user.ID, "password", user.Password, "nickname", user.Nickname, "introduction", user.Introduction, "avatar", user.Avatar)
			if err != nil {
				utils.ResponseFail(c, err.Error())
				return
			}
			utils.ResponseSuccess(c, "register success")

		} else {
			utils.ResponseFail(c, fmt.Sprintf("register failed,%s", msg))
			return
		}
		e.Exit()
	}
}

func login(c *gin.Context) {
	e, b := sentinel.Entry("chatchat_user", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		// 请求被拒绝，在此处进行处理
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  429,
			"message": "too many request",
		})
	} else {
		var user model.User
		if err := c.ShouldBind(&user); err != nil {
			fmt.Println(err)
			utils.ResponseFail(c, "verification failed")
			return
		}
		username := user.Username
		password := user.Password
		id, _ := redis.HGet(c, fmt.Sprintf("user:%s", username), "id")
		user.ID, _ = strconv.ParseInt(id.(string), 10, 64)

		flag, _ := redis.HGet(c, fmt.Sprintf("user:%s", username), "password")
		if flag == "" {
			utils.ResponseFail(c, "user doesn't exists")
			return
		}
		RedisPassword, _ := redis.HGet(c, fmt.Sprintf("user:%s", username), "password")
		if !utils.EqualsPassword(password, RedisPassword.(string)) {
			utils.ResponseFail(c, "wrong password")
			return
		}
		claim := model.MyClaims{
			Username: username,
			ID:       user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
				Issuer:    "Wzy",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, _ := token.SignedString(middleware.Secret)
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "login success",
			"token":   tokenString,
		})
		c.Set("id", user.ID)
		c.Next()
		e.Exit()
		return
	}
}

func ChangePassword(c *gin.Context) {
	username, _ := c.Get("username")
	//id, _ := redis.HGet(c, fmt.Sprintf("user:%s", username), "")
	var user model.User
	c.ShouldBind(&user)
	newPassword := user.NewPassword
	mailID := user.MailID
	uid, _ := redis.Get(c.Request.Context(), fmt.Sprintf("mail:%s", username))
	if mailID != uid {
		utils.ResponseFail(c, "wrong mailID")
		return
	}
	flag1, msg := dao.AddUserCheck("2604296771@qq.com", newPassword, "ziYu")
	if !flag1 {
		utils.ResponseFail(c, msg)
		return
	}
	EncryptPassword, err1 := utils.EncryptPassword(newPassword) //加密密码
	if err1 != nil {
		utils.ResponseFail(c, "encrypt password failed")
		return
	}
	u := model.User{
		Username: fmt.Sprintf("%s", username),
		Password: newPassword,
	}
	flag2, msg := mysql.ChangePassword(c.Request.Context(), u)
	if flag2 {
		utils.ResponseSuccess(c, "password change success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("password change failed,%s", msg))
		return
	} //更新数据库数据
	err := redis.HSet(c.Request.Context(), fmt.Sprintf("user:%s", username), "password", EncryptPassword) //重新写入到redis
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
	}
}

func ChangeNickname(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	c.ShouldBind(&user)
	nickname := user.Nickname
	flag1, msg := dao.AddUserCheck("2604296771@qq.com", "1234512345w", nickname)
	if !flag1 {
		utils.ResponseFail(c, msg)
		return
	}
	u := model.User{
		Username: fmt.Sprintf("%s", username),
		Nickname: nickname,
	}
	flag2, msg := mysql.ChangeNickname(c.Request.Context(), u)
	if flag2 {
	} else {
		utils.ResponseFail(c, fmt.Sprintf("nickname change failed,%s", msg))
		return
	} //更新数据库数据
	//global.Rdb.Del(c, fmt.Sprintf("%s:nickname", username), nickname)      //删除缓存中的键值对
	err := redis.HSet(c, fmt.Sprintf("user:%s", username), "nickname", nickname) //重新写入到redis
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
		return
	}
	utils.ResponseSuccess(c, "change nickname success")
}

func ChangeAvatar(c *gin.Context) {
	var user model.User
	// 解析表单数据，设置最大文件大小
	username, _ := c.Get("username")
	var i int
	for k, v := range username.(string) {
		if v == '@' {
			i = k
			break
		}
	}
	user.Username = username.(string)
	num := username.(string)[:i]
	user.Avatar = "http://test.violapioggia.cn/chatchatUsers/" + num + "%40" + "qq.com"
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		// 处理错误
		utils.ResponseFail(c, "avatar's size is beyond the max size")
		return
	}

	// 获取上传的文件
	avatar, _, err := c.Request.FormFile("image")
	if err != nil {
		// 处理错误
		utils.ResponseFail(c, err.Error())
		return
	}
	defer avatar.Close()
	err = utils.Delete(username.(string))
	err = utils.Upload(avatar, username.(string))
	if err != nil {
		utils.ResponseFail(c, fmt.Sprintf("change avatar failed,err:%s", err.Error()))
		return
	}
	b, _ := mysql.ChangeAvatar(c.Request.Context(), user)
	if !b {
		utils.ResponseFail(c, "update avatar failed")
	}
	redis.HSet(c, fmt.Sprintf("user:%s", username), "avatar", user.Avatar)
	utils.ResponseSuccess(c, "change avatar success")
	return
}
func ChangeIntroduction(c *gin.Context) {
	username, _ := c.Get("username")
	//id ,_:=c.Get("id")
	//c.Param(id)
	var user model.User
	c.ShouldBind(&user)
	introduction := user.Introduction
	u := model.User{
		Username:     fmt.Sprintf("%s", username),
		Introduction: introduction,
	}
	flag2, msg := mysql.ChangeIntroduction(c.Request.Context(), u)
	if flag2 {
	} else {
		utils.ResponseFail(c, fmt.Sprintf("introduction change failed,%s", msg))
		return
	} //更新数据库数据
	err := redis.HSet(c, fmt.Sprintf("user:%s", username), "introduction", introduction) //重新写入到redis
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
		return
	}
	utils.ResponseSuccess(c, "change introduction success")
}

//	func ChangeAvatar(c *gin.Context) {
//		username, _ := c.Get("username")
//		id := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
//		ID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
//		c.Param(id)
//		var user model.User
//		c.ShouldBind(&user)
//		avatar := user.Avatar
//		u := model.User{
//			Username: fmt.Sprintf("%s", username),
//			Avatar:   avatar,
//			ID:       ID,
//		}
//		flag2, msg := mysql.ChangeAvatar(u)
//		if flag2 {
//		} else {
//			utils.ResponseFail(c, fmt.Sprintf("avatar change failed,%s", msg))
//			return
//		} //更新数据库数据
//		global.Rdb.Del(c, fmt.Sprintf("%s:avatar", username), avatar)      //删除缓存中的键值对
//		err := redis.Set(c, fmt.Sprintf("%s:avatar", username), avatar, 0) //重新写入到redis
//		if err != nil {
//			utils.ResponseFail(c, "write into redis failed")
//			return
//		}
//		utils.ResponseSuccess(c, "change avatar success")
//	}
//
//	func getUsernameFromToken(c *gin.Context) {
//		username, _ := c.Get("username")
//		utils.ResponseSuccess(c, username.(string))
//	}
//
//	func FindPassword(c *gin.Context) {
//		username := c.PostForm("username")
//		password,_ := redis.HGet(c,"user:%s",username,"password")
//		if password=="" {
//			utils.ResponseFail(c, "user doesn't exists")
//			return
//		}
//		utils.ResponseSuccess(c, password.(string))
//	}
func GetUser(c *gin.Context) {
	username, _ := c.Get("username")
	id, _ := global.Rdb.HGet(c, fmt.Sprintf("user:%s", username), "id").Int64()
	avatar := global.Rdb.HGet(c, fmt.Sprintf("user:%s", username), "avatar").Val()
	nickname := global.Rdb.HGet(c, fmt.Sprintf("user:%s", username), "nickname").Val()
	introduction := global.Rdb.HGet(c, fmt.Sprintf("user:%s", username), "introduction").Val()
	//likeArticle := global.Rdb.HGet(c, fmt.Sprintf("user:%s", username),"id").Val()
	//uid := c.Param(id)
	//c.String(200, "%v", uid)
	var user model.User
	user.Avatar = avatar
	user.Username = username.(string)
	user.Nickname = nickname
	user.Introduction = introduction
	user.ID = id

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		//"username":     username,
		//"avatar":       avatar,
		//"nickname":     nickname,
		//"introduction": introduction,
		//"id":           id,
		"data": user,
		//"likeArticle":  likeArticle,
	})
}
