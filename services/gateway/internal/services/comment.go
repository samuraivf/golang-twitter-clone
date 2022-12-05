package services

import (
	commentService "comment/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetCommentClient() (commentService.CommentClient, func() error) {
	connection := connectCommentGrpc()
	client := commentService.NewCommentClient(connection)

	return client, connection.Close
}


func connectCommentGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.comment"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to comment service: %s", err.Error())
	}

	return connection
}