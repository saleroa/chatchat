package dao

import (
	"chatchat/model"
	"chatchat/utils"
)

var user *model.User

func AddUserCheck(username, password, nickname string) (bool, string) {
	var i int
	for k, v := range username {
		if v == '@' {
			i = k
			break
		}
	}
	num := username[:i]
	mail := username[i:]
	if !utils.IsNum(num) || mail != "@qq.com" {
		return false, "it is not the format of qq's email"
	}
	if len(password) < 10 || len(password) > 20 {
		return false, "length of the password is wrong "
	}
	if utils.IsNum(password) {
		return false, "the password can't have only nums"
	}
	if len(nickname) < 2 || len(nickname) > 8 {
		return false, "length of the nickname is wrong"
	}
	return true, ""
}
