package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(dbUser, dbPassword, dbHost, dbName string) (*gorm.DB, error) {
	conn := createConnectionString(dbUser, dbPassword, dbHost, dbName)
	return gorm.Open(postgres.Open(conn), createDefaultConfig())
}

func createConnectionString(dbUser, dbPassword, dbHost, dbName string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", dbUser, dbPassword, dbHost, dbName)
}

func createDefaultConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	}
}
