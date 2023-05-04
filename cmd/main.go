package main

import (
	"chatchat/boot"
	"chatchat/utils"
)

func main() {
	boot.ViperSetup(utils.Path())
	boot.Loggersetup()
	boot.MysqlDBSetup()
	boot.RedisSetup()
	boot.ServerSetup()
}
