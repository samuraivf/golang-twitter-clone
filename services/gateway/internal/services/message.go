package services

import (
	messageService "message/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetMessageClient() (messageService.MessageClient, func() error) {
	connection := connectMessageGrpc()
	client := messageService.NewMessageClient(connection)

	return client, connection.Close
}


func connectMessageGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.message"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to message service: %s", err.Error())
	}

	return connection
}