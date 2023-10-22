package main

import (
	"FioapiKafka/internal/consumer"
	repository2 "FioapiKafka/internal/repository"
	"FioapiKafka/internal/repository/database"
	"FioapiKafka/internal/repository/database/migrations"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"FioapiKafka/internal/cache"
	"FioapiKafka/internal/handler"
	"FioapiKafka/internal/services"
	"context"
)

func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %s", err.Error())
	}

	// Инициализация базы данных
	db, err := migrations.Migrate(&database.Config{
		Uname:      os.Getenv("DB_UNAME"),
		Pass:       os.Getenv("DB_PASS"),
		NameDB:     os.Getenv("DB_NAMEDB"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		SSL:        os.Getenv("DB_SSL"),
		DriverName: os.Getenv("DB_DriverName"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err.Error())
	}
	defer db.Close()

	// Инициализация Redis
	cacheClient, err := cache.InitRedis()
	if err != nil {
		log.Fatalf("Failed to connect to the Redis: %s", err.Error())
	}
	defer cacheClient.Close()
	// Инициализация репозитория
	repository := repository2.NewRepository(db, cacheClient)
	// Инициализация сервиса
	service := services.NewService(repository)

	// Инициализация GoFiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	app.Use(recover.New())

	// Инициализация API хэндлеров
	apiHandlers := handler.NewHandlers(service)
	app.Get("/handler/people", apiHandlers.GetPeople)
	app.Get("/handler/people/:id", apiHandlers.GetPersonByID)
	app.Post("/handler/people", apiHandlers.CreatePerson)
	app.Put("/handler/people/:id", apiHandlers.UpdatePerson)
	app.Delete("/handler/people/:id", apiHandlers.DeletePerson)

	var cancels sync.Map // Сохраняем все cancel-функции
	// Запуск Consumer FIO
	cancels.Store("fio_consumer", consumer.InitConsumerFIO([]string{os.Getenv("KAFKA_BROKER")}, os.Getenv("KAFKA_FIO_TOPIC")))
	// Запуск Consumer FIO_FAILED
	cancels.Store("failed_consumer", consumer.InitConsumerFailed([]string{os.Getenv("KAFKA_BROKER")}, os.Getenv("KAFKA_FIO_FAILED_TOPIC"))) // Запуск сервера в отдельной goroutine
	go func() {
		log.Println("Starting HTTP server...")
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Failed to start HTTP server: %s", err.Error())
		}
	}()

	// Ожидание сигнала завершения приложения (CTRL+C)
	waitForShutdown(&cancels)

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		log.Printf("Graceful shutdown failed: %s", err.Error())
	} else {
		log.Println("Server stopped")
	}
}

// Ожидание сигнала завершения приложения
func waitForShutdown(cancels *sync.Map) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-stopChan:
		log.Println("Shutting down...")

		cancels.Range(func(k, v interface{}) bool {
			if cancel, ok := v.(context.CancelFunc); ok {
				cancel()
			}
			return true
		})
	}
}
