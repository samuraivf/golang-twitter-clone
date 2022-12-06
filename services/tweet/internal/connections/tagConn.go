package connections

import (
	tagService "tag/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetTagClient() (tagService.TagClient, func() error) {
	connection := connectTagGrpc()
	client := tagService.NewTagClient(connection)

	return client, connection.Close
}

func connectTagGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.tag"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to tag service: %s", err.Error())
	}

	return connection
}