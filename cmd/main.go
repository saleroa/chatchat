package main

import (
	"chatchat/boot"
)

func main() {
	boot.ViperSetup("./config.yaml")
	boot.Loggersetup()
	boot.MysqlDBSetup()
	boot.RedisSetup()
	boot.ServerSetup()
}
