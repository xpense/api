package main

import (
	"expense-api/repository"
	"expense-api/router"
	"fmt"
	"log"
	"os"
	"time"

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

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		panic(fmt.Sprintf("couldn't establish postgres connection: %v", err))
	}

	if err := repository.Migrate(db); err != nil {
		panic(fmt.Sprintf("error setting up database: %v", err))
	}

	repository := repository.New(db)
	r := router.Setup(repository)

	r.Run(port)
}
