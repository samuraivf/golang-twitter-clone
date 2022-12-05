package main

import (
	"tag/internal/repo"
	"tag/internal/service"
	"tag/configs"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "tag/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	port := configs.GetPort()
	postgresConfig := configs.GetPostgresConfig()

	db, err := repository.NewPostgresDB(*postgresConfig)

	if err != nil {
		logrus.Fatalf("failed to connect to postgres: %s", err.Error())
	}

	tagRepo := repository.NewTagPostgres(db)
	tagService := service.NewTagService(tagRepo)
	server := grpc.NewServer()
	pb.RegisterTagServer(server, tagService)

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

	postgresDB, err := db.DB()

	if err != nil {
		logrus.Errorf("error occured while getting postgres database: %s", err.Error())
	}

	if err := postgresDB.Close(); err != nil {
		logrus.Errorf("error occured while closing postgres connection: %s", err.Error())
	}

	logrus.Print("postgres connection is closed")

	if err := lis.Close(); err != nil {
		logrus.Errorf("failed to stop a tcp server: %s", err.Error())
	}

	logrus.Print("tcp server is stopped")
}
