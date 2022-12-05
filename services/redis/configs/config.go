package configs

import (
	repository "redis/internal/repo"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func GetPort() string {
	return viper.GetString("port")
}

func GetRedisConfig() *repository.RedisConfig {
	return &repository.RedisConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
	}
}