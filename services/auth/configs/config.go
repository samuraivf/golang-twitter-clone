package configs

import (
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("services/auth/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func GetPort() string {
	return viper.GetString("port")
}
