package router

import (
	"expense-api/handlers"
	"expense-api/repository"

	"github.com/gin-gonic/gin"
)

func Setup(r repository.Repository) *gin.Engine {
	router := gin.Default()
	transaction := router.Group("/transaction")
	{
		transaction.GET("/", handlers.GetTransactions(r))
		transaction.POST("/", handlers.CreateTransaction(r))
		transaction.GET("/:id", handlers.GetTransaction(r))
		transaction.PATCH("/:id", handlers.UpdateTransaction(r))
		transaction.DELETE("/:id", handlers.DeleteTransaction(r))
	}

	return router
}
