package main

import (
	"auth/internal/service"
	"auth/configs"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "auth/proto"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	port := configs.GetPort()

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
