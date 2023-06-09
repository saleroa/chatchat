package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	NewPassword string `form:"newPassword" json:"newPassword"`
	Nickname    string `form:"nickname" json:"nickname" `
	ID          int64  `form:"id" json:"id" `
	MailID      string `form:"mailID" json:"mailID" `
	//VIP          int       `form:"vip" json:"vip" `
	Avatar       string    `form:"avatar" json:"avatar" `
	FriendID     int64     `form:"friendID" json:"friendID" `
	GroupID      int64     `form:"groupID" json:"groupID" `
	Introduction string    `form:"introduction" json:"introduction" `
	CreateTime   time.Time `form:"create_time" json:"create_time" `
	UpdateTime   time.Time `form:"update_time" json:"update_time" `
}
type OauthUser struct {
	ExpiresIn      int64  `form:"expires_in" json:"expires_in" `
	ClientID       string `form:"client_id" json:"client_id" `
	Avatar         string `form:"avatar" json:"avatar" `
	Nickname       string `form:"nickname" json:"nickname" `
	Oauth2Username int64  `form:"oauth2Username" json:"oauth2Username" `
}
type MyClaims struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
	jwt.StandardClaims
}
