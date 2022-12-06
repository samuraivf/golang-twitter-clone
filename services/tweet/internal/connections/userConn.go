package connections

import (
	userService "user/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetUserClient() (userService.UserClient, func() error) {
	connection := connectUserGrpc()
	client := userService.NewUserClient(connection)

	return client, connection.Close
}

func connectUserGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.user"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to user service: %s", err.Error())
	}

	return connection
}