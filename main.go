package main

import (
	"expense-api/model"
	"expense-api/repository"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

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

	if err := db.AutoMigrate(); err != nil {
		panic(fmt.Sprintf("error setting up database: %v", err))
	}

	r := gin.Default()

	expenses := repository.NewDB()

	r.GET("/admin/overview/expense", func(c *gin.Context) {
		c.JSON(http.StatusOK, expenses)
	})

	r.POST("/:user/expense", func(c *gin.Context) {
		user := c.Param("user")

		var json model.Transaction
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		t, err := model.NewTransaction(json.Timestamp, json.Amount, json.Type)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		expenses.AddTransaction(user, t)

		c.JSON(http.StatusOK, gin.H{
			"id": t.ID,
		})
	})

	r.Run(port)
}
