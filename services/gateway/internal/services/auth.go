package services

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectAuthGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.auth"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to auth service: %s", err.Error())
	}

	return connection
}