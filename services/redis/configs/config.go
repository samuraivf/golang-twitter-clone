package configs

import (
	repository "redis/internal/repo"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("services/redis/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func GetPort() string {
	return viper.GetString("port")
}

func GetRedisConfig() *repository.RedisConfig {
	return &repository.RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
	}
}