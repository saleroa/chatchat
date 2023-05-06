package dao

import (
	"chatchat/model"
	"chatchat/utils"
)

var user *model.User

func AddUserCheck(username, password, nickname string) (bool, string) {
	if len(username) < 6 || len(username) > 10 {
		return false, "length of the username is wrong "
	}
	if !utils.IsNum(username) {
		return false, "the username has covered characters"
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
