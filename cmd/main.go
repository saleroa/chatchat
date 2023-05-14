package main

import (
	"chatchat/app/api/middleware"
	"chatchat/boot"
)

func main() {
	boot.ViperSetup("./config.yaml")
	boot.Loggersetup()
	boot.MysqlDBSetup()
	boot.RedisSetup()
	middleware.SessionSetup()
	boot.ServerSetup()
}
