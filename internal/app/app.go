package app

import (
	"context"
	"wb-l0/internal/cache"
	"wb-l0/internal/config"
	"wb-l0/internal/controller/handler"
	"wb-l0/internal/controller/router"
	"wb-l0/internal/nats-streaming/subscribe"
	"wb-l0/internal/repository"
	"wb-l0/internal/service"
	"wb-l0/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Run() error {
	config.InitEnvConfigs()

	pool, err := database.NewPostgresDB(
		context.Background(),
		database.Config{
			Host:     config.EnvConfig.DBHost,
			Port:     config.EnvConfig.DBPort,
			Username: config.EnvConfig.DBUsername,
			Name:     config.EnvConfig.DBName,
			Password: config.EnvConfig.DBPassword,
		})
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer pool.Close()

	cache := cache.NewCache(10, 5)
	repo := repository.NewRepository(pool, cache)
	ords, err := repo.FindAll(context.Background())
	cache.RestoreCache(ords)

	subscriber := subscribe.New(repo)
	go subscriber.SubAndPub()

	app := fiber.New()
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	server := router.NewServer(app, handlers)

	server.Router()
	app.Listen(":" + config.EnvConfig.LocalServerPort)

	return nil
}
