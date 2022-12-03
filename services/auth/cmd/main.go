package main

import (
	"auth/internal/service"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	pb "auth/proto"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	port := viper.GetString("port")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logrus.Fatalf("cannot listen a tcp server: %s", err.Error())
	}

	authService := service.NewAuthService()
	server := grpc.NewServer()
	pb.RegisterAuthorizationServer(server, authService)

	logrus.Printf("Server start at port %s", port)

	server.Serve(lis)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}