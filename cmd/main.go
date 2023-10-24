package main

import (
	"FioapiKafka/internal/consumer"
	"FioapiKafka/internal/graphqlhandler"
	"FioapiKafka/internal/models"
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
	log.SetFormatter(new(log.JSONFormatter))
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env file: %s", err.Error())
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
		log.Fatalf("failed to connect to the database: %s", err.Error())
	}
	defer db.Close()

	// Инициализация Redis
	cacheClient, err := cache.InitRedis(&cache.RedisConfig{
		Addr: os.Getenv("REDIS_ADDR"),
		Pass: os.Getenv("REDIS_PASS"),
	})
	if err != nil {
		log.Fatalf("failed to connect to the Redis: %s", err.Error())
	}
	defer cacheClient.Close()
	// Инициализация репозитория
	log.Infoln("Init repository...")
	repository := repository2.NewRepository(db, cacheClient)
	// Инициализация сервиса
	log.Infoln("Init service...")
	service := services.NewService(repository)

	// Инициализация GoFiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	app.Use(recover.New())

	// GraphQL handlers
	log.Infoln("Init GraphQL handler")
	gh := graphqlhandler.NewGraphQLHandler(service)
	app.All("/graphql", func(c *fiber.Ctx) error {
		result := gh.ExecuteQuery(string(c.Body()))
		return c.JSON(result)
	})

	// Инициализация API хэндлеров
	log.Infoln("Init handlers...")
	apiHandlers := handler.NewHandlers(service)
	app.Get("/handler/people", apiHandlers.GetPersons)
	app.Get("/handler/people/:id", apiHandlers.GetPersonByID)
	app.Get("/handler/people/name/:name", apiHandlers.GetPersonsByName)
	app.Get("/handler/people/age/:age", apiHandlers.GetPersonsByAge)
	app.Get("/handler/people/gender/:gender", apiHandlers.GetPersonsByGender)
	app.Post("/handler/people", apiHandlers.CreatePerson)
	app.Put("/handler/people/:id", apiHandlers.UpdatePerson)
	app.Delete("/handler/people/:id", apiHandlers.DeletePerson)

	// Вывод роутов в консоль
	log.Infoln("GET /handler/people")
	log.Infoln("GET /handler/people/:id")
	log.Infoln("GET handler/people/name/:name")
	log.Infoln("GET handler/people/age/:age")
	log.Infoln("GET handler/people/gender/:gender")
	log.Infoln("POST /handler/people")
	log.Infoln("PUT /handler/people/:id")
	log.Infoln("DELETE /handler/people/:id")
	log.Infoln("/graphql")

	// Канал который будет хранить данные о прошедших проверку FIO
	aproovedFIO := make(chan *models.PersonOut, 100)

	// Запуск воркера для чтения из канала и запись в бд прошедших проверку FIO
	log.Infoln("Star worker CreatePeople...")
	go service.CreatePeople(aproovedFIO)

	var cancels sync.Map // Сохраняем все cancel-функции
	// Запуск Consumer FIO
	log.Infoln("Start FIO consumer...")
	cancels.Store("fio_consumer", consumer.InitConsumerFIO([]string{os.Getenv("KAFKA_BROKER")}, os.Getenv("KAFKA_FIO_TOPIC"), aproovedFIO))
	// Запуск Consumer FIO_FAILED
	log.Infoln("Start FIO_FAILED consumer")
	cancels.Store("failed_consumer", consumer.InitConsumerFailed([]string{os.Getenv("KAFKA_BROKER")}, os.Getenv("KAFKA_FIO_FAILED_TOPIC"))) // Запуск сервера в отдельной goroutine
	go func() {
		log.Println("Starting HTTP server...")
		if err := app.Listen(os.Getenv("SRV_PORT")); err != nil {
			log.Fatalf("Failed to start HTTP server: %s", err.Error())
		}
	}()

	// Ожидание сигнала завершения приложения
	waitForShutdown(&cancels)

	// shutdown
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
