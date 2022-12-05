package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gateway/internal/handler"
	"gateway/internal/server"
	"gateway/configs"
	"github.com/sirupsen/logrus"
)

// @title           Twitter Clone
// @version         1.0
// @description     Application with basic functionality of twitter

// @host      localhost:7000
// @BasePath  /api/

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	handler := initApplication(ctx)

	server := new(server.Server)
	go func() {
		if err := server.Run(configs.GetPort(), handler.InitServer()); err != http.ErrServerClosed {
			logrus.Fatalf("failed to start the server: %s", err.Error())
		}
	}()
	logrus.Print("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server is shutting down")

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initApplication(ctx context.Context) (*handler.Handler) {
	logrus.SetFormatter(&logrus.TextFormatter{})

	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	handler := handler.NewHandler()

	return handler
}
