package services

import (
	tweetService "tweet/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetTweetClient() (tweetService.TweetClient, func() error) {
	connection := connectTweetGrpc()
	client := tweetService.NewTweetClient(connection)

	return client, connection.Close
}

func connectTweetGrpc() *grpc.ClientConn {
	var connection *grpc.ClientConn
	connection, err := grpc.Dial(viper.GetString("services.tweet"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("cannot connect to tweet service: %s", err.Error())
	}

	return connection
}