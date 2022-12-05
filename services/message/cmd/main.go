package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"message/internal/repo"
	"message/internal/service"
	"syscall"

	pb "message/proto"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	port := viper.GetString("port")

	config := repository.PostgresConfig{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	}

	db, err := repository.NewPostgresDB(config)

	if err != nil {
		logrus.Fatalf("failed to connect to postgres: %s", err.Error())
	}

	messageRepo := repository.NewMessagePostgres(db)
	messageService := service.NewMessageService(messageRepo)
	server := grpc.NewServer()
	pb.RegisterMessageServer(server, messageService)

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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}