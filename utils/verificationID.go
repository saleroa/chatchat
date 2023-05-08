package utils

import (
	"math/rand"
	"time"
)

func GetVerificationID() int {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 调用生成随机数的函数
	randomNumber := generateRandomNumber()

	return randomNumber
}

func generateRandomNumber() int {
	// 生成6位随机数
	min := 100000
	max := 999999
	return rand.Intn(max-min+1) + min
}
