package main

import (
	"context"
	"os"

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
	handler := initApplication(ctx)
	server := new(server.Server)

	if err := server.Run(viper.GetString("port"), handler.InitServer()); err != nil {
		logrus.Fatalf("failed to start the server: %s", err.Error())
	}

	logrus.Print("Server started")
}

func initApplication(ctx context.Context) *handler.Handler {
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

	return handler
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
