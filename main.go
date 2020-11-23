package main

import (
	"expense-api/middleware/auth"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/utils"
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
	if _, err := os.Stat(".env"); err != nil {
		panic(fmt.Sprintf("couldn't find env file: %v", err))
	}

	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("couldn't load env file: %v", err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		panic("JWT_ISSUER not set!")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set!")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	postgresConnection := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", dbUser, dbPassword, dbHost, dbName)

	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{
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
	jwtService := auth.NewJWTService(issuer, secret)
	hasher := utils.NewPasswordHasher()

	r := router.Setup(repository, jwtService, hasher, router.DefaultConfig)
	r.Run(port)
}
