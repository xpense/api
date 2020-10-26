package main

import (
	"expense-api/handlers"
	"expense-api/model"
	"expense-api/repository"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if _, err := os.Stat(".env.yml"); err != nil {
		panic(fmt.Sprintf("couldn't find env file: %v", err))
	}

	if err := godotenv.Load(".env.yml"); err != nil {
		panic(fmt.Sprintf("couldn't load env file: %v", err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	postgresURI := os.Getenv("POSTGRES_CONNECTION")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("couldn't establish postgres connection: %v", err))
	}

	if err := db.AutoMigrate(model.Transaction{}); err != nil {
		panic(fmt.Sprintf("error setting up database: %v", err))
	}

	repository := repository.New(db)

	r := gin.Default()

	transaction := r.Group("/transaction")
	{
		transaction.GET("/", handlers.GetTransactions(repository))
		transaction.POST("/", handlers.CreateTransaction(repository))
		transaction.GET("/:id", handlers.GetTransaction(repository))
		transaction.PATCH("/:id", handlers.UpdateTransaction(repository))
		transaction.DELETE("/:id", handlers.DeleteTransaction(repository))
	}

	r.Run(port)
}
