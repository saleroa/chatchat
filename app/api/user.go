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
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func register(c *gin.Context) {
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
	uid, err := redis.Get(c, fmt.Sprintf("Rmail:%s", username))
	if err != nil {
		utils.ResponseFail(c, "verification code has expired")
		return
	}
	if uid == "" || uid != mailID {
		utils.ResponseFail(c, "wrong mailID")
	}

	user.ID = global.Rdb.ZCard(c, "userID").Val() + 1
	flag1, msg := mysql.AddUser(username, password, nickname, user.ID) //写入数据库
	if flag1 {
		err := redis.ZSetUserID(c, username)
		if err != nil {
			utils.ResponseFail(c, err.Error())
		}
		err = redis.HSet(c, fmt.Sprintf("user:%s", username), "id", user.ID, "password", user.Password, "nickname", user.Nickname, "introduction", user.Introduction, "avatar", user.Avatar)
		if err != nil {
			utils.ResponseFail(c, err.Error())
		}
		utils.ResponseSuccess(c, "register success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("register failed,%s", msg))
	}

}

func login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}
	username := user.Username
	password := user.Password

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
	return
}

func ChangePassword(c *gin.Context) {
	username, _ := c.Get("username")
	//id, _ := redis.HGet(c, fmt.Sprintf("user:%s", username), "")
	var user model.User
	c.ShouldBind(&user)
	newPassword := user.NewPassword
	mailID := user.MailID
	uid, _ := redis.Get(c, fmt.Sprintf("mail:%s", username))
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
	flag2, msg := mysql.ChangePassword(u)
	if flag2 {
		utils.ResponseSuccess(c, "password change success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("password change failed,%s", msg))
		return
	} //更新数据库数据
	err := redis.HSet(c, fmt.Sprintf("user:%s", username), "password", EncryptPassword) //重新写入到redis
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
	flag2, msg := mysql.ChangeNickname(u)
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
	// 解析表单数据，设置最大文件大小
	username, _ := c.Get("username")
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
		utils.ResponseFail(c, "wrong format of the avatar")
		return
	}
	defer avatar.Close()
	// 假设前端传来的图片数据存储在变量 imageData 中

	// 将图片数据写入临时文件
	tempFile, err := ioutil.ReadAll(avatar)

	// 获取临时文件的路径
	imagePath := string(tempFile)
	utils.Upload(imagePath, username.(string))

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
	flag2, msg := mysql.ChangeIntroduction(u)
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
