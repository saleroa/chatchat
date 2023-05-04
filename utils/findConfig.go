package utils

import (
	"gopkg.in/ini.v1"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAppPath() string {
	file, _ := exec.LookPath("config.yaml")
	path, _ := filepath.Abs(file)
	//index := strings.LastIndex(path, string(os.PathSeparator))
	//path = path[:index]

	_, err := ini.Load(path + "/manifest/config/config.yaml")
	if err != nil {
		log.Fatal("配置文件读取失败, err = ", err)
	}
	return path
}
func Path() string {
	var build strings.Builder
	build.WriteString(GetAppPath())
	build.WriteString("/manifest/config/config.yaml")
	path := build.String()
	return path
}
