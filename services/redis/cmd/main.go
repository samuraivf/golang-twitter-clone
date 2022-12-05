package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"redis/internal/repo"
	"redis/internal/service"
	"redis/configs"
	"syscall"

	pb "redis/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	port := configs.GetPort()
	redisConfig := configs.GetRedisConfig()

	redisClient, err := repository.NewRedisClient(context.Background(), *redisConfig)

	if err != nil {
		logrus.Fatalf("failed to connect to redis: %s", err.Error())
	}

	redisRepo := repository.NewRedis(redisClient)
	redisService := service.NewRedisService(redisRepo)
	server := grpc.NewServer()
	pb.RegisterRedisServer(server, redisService)

	var lis net.Listener

	go func() {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%s", port))

		if err != nil {
			logrus.Fatalf("cannot listen a tcp server: %s", err.Error())
		}
		if err := server.Serve(lis); err != nil {
			logrus.Fatalf("failed to start a grpc server: %s", err.Error())
		}
	}()

	logrus.Printf("Server start at port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server is shutting down")

	server.Stop()

	logrus.Print("grpc server is stopped")

	if err := redisClient.Conn().Close(); err != nil {
		logrus.Errorf("failed to close a redis connection: %s", err.Error())
	}
	
	logrus.Print("redis connection is closed")

	if err := lis.Close(); err != nil {
		logrus.Errorf("failed to stop a tcp server: %s", err.Error())
	}

	logrus.Print("tcp server is stopped")
}
