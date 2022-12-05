package services

import (
	redisService "redis/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetRedisClient() (redisService.RedisClient, func() error) {
	connection := connectRedisGrpc()
	client := redisService.NewRedisClient(connection)

	return client, connection.Close
}

func connectRedisGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.redis"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to redis service: %s", err.Error())
	}

	return connection
}