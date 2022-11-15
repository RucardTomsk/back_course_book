package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	coursebook "github.com/RucardTomsk/course_book"
	"github.com/RucardTomsk/course_book/pkg/handler"
	"github.com/RucardTomsk/course_book/pkg/repository"
	"github.com/RucardTomsk/course_book/pkg/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
	Neo4jDriver, err := repository.NewNeo4jDriver(repository.Config{
		URI:      viper.GetString("db.uri"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize Neo4jDriver: %s", err.Error())
	}

	repos := repository.NewRepository(Neo4jDriver)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(coursebook.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error accured while running http server: %s", err.Error())
		}
	}()

	logrus.Printf("course-book-API start PORT :%s", viper.GetString("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("course-book-API Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
