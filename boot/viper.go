package boot

import (
	"chatchat/app/global"
	"github.com/spf13/viper"
)

func ViperSetup(ConfigPath string) {
	v := viper.New()
	v.SetConfigFile(ConfigPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if err = v.Unmarshal(&global.Config); err != nil {
		panic(err)
	}

}
