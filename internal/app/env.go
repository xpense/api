package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	PORT        = "PORT"
	JWT_ISSUER  = "JWT_ISSUER"
	JWT_SECRET  = "JWT_SECRET"
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_HOST     = "DB_HOST"
	DB_NAME     = "DB_NAME"
)

type environment struct {
	Port       string
	Issuer     string
	Secret     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
}

func LoadEnvironmentVars() environment {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("couldn't load env file: %v", err))
	}

	port := os.Getenv(PORT)
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	issuer := os.Getenv(JWT_ISSUER)
	if issuer == "" {
		panic("JWT_ISSUER not set!")
	}

	secret := os.Getenv(JWT_SECRET)
	if secret == "" {
		panic("JWT_SECRET not set!")
	}

	dbUser := os.Getenv(DB_USER)
	dbPassword := os.Getenv(DB_PASSWORD)
	dbHost := os.Getenv(DB_HOST)
	dbName := os.Getenv(DB_NAME)

	return environment{
		Port:       port,
		Issuer:     issuer,
		Secret:     secret,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBHost:     dbHost,
		DBName:     dbName,
	}
}
