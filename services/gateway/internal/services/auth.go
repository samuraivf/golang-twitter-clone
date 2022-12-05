package services

import (
	authService "auth/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetAuthClient() (authService.AuthorizationClient, func() error) {
	connection := connectAuthGrpc()
	client := authService.NewAuthorizationClient(connection)

	return client, connection.Close
}


func connectAuthGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.auth"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to auth service: %s", err.Error())
	}

	return connection
}