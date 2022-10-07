package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/samuraivf/twitter-clone/internal/handler"
	"github.com/samuraivf/twitter-clone/internal/repository"
	"github.com/samuraivf/twitter-clone/internal/server"
	"github.com/samuraivf/twitter-clone/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	handler, closeRepos := initApplication(ctx)
	
	server := new(server.Server)
	go func() {
		if err := server.Run(viper.GetString("port"), handler.InitServer()); err != http.ErrServerClosed {
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

	closeRepos()
}

func initApplication(ctx context.Context) (*handler.Handler, func()) {
	logrus.SetFormatter(&logrus.TextFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("failed to load env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	redis, err := repository.NewRedisClient(ctx, repository.RedisConfig{
		Host: viper.GetString("redis.host"),
		Port: viper.GetString("redis.port"),
	})

	if err != nil {
		logrus.Fatalf("failed to connect to redis: %s", err.Error())
	}

	repos := repository.NewRepository(db, redis)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

	return handler, func() {
		if err := redis.Close(); err != nil {
			logrus.Errorf("error occured while closing redis connection: %s", err.Error())
		}
		postgresDB, err := db.DB()

		if err != nil {
			logrus.Errorf("error occured while getting postgres database: %s", err.Error())
		}

		if err := postgresDB.Close(); err != nil {
			logrus.Errorf("error occured while closing postgres connection: %s", err.Error())
		}

		logrus.Print("repositories are closed")
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
