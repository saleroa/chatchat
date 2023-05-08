package api

import (
	"chatchat/dao"
	"chatchat/dao/redis"
	"chatchat/utils"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"strconv"
	"time"
)

func SendMail(c *gin.Context) {
	username, f := c.GetPostForm("username")
	if !f {
		utils.ResponseFail(c, "get username failed")
		return
	}
	flag, _ := redis.HGet(c, fmt.Sprintf("user:%s", username), "password")
	if flag != "" {
		utils.ResponseFail(c, "user already exists")
		return
	}
	flag1, msg := dao.AddUserCheck(username, "wzywzywzywzy", "紫雨")
	if !flag1 {
		utils.ResponseFail(c, msg)
		return
	}
	uid := utils.GetVerificationID()
	Mail(username, uid)
	err := redis.Set(c, fmt.Sprintf("mail:%s", username), uid, 5*time.Minute)
	if err != nil {
		utils.ResponseFail(c, "write mailID into redis failed")
		return
	}
	utils.ResponseSuccess(c, "send email to the user success")
	return
}
func Mail(username string, uid int) {
	m := gomail.NewMessage()
	m.SetHeader("From", "violapioggia@qq.com")
	m.SetHeader("To", username)
	m.SetHeader("Subject", "Verify to login into the chatchat")
	m.SetBody("text/plain", "Hello! 你的验证码是"+strconv.Itoa(uid)+"，不要告诉别人哦~")
	m.Attach("./images/mail.png")

	host := "smtp.qq.com"
	port := 25
	userName := "violapioggia@qq.com"
	password := "htuhlncfmjdqdicj" // qq邮箱填授权码
	d := gomail.NewDialer(
		host,
		port,
		userName,
		password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
