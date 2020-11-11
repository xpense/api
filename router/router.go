package router

import (
	"expense-api/handlers"
	"expense-api/repository"

	"github.com/gin-gonic/gin"
)

// Setup creates a new gin router
func Setup(repo repository.Repository) *gin.Engine {
	router := gin.Default()

	transaction := router.Group("/transaction")
	{
		transaction.GET("/", handlers.ListTransactions(repo))
		transaction.POST("/", handlers.CreateTransaction(repo))
		transaction.GET("/:id", handlers.GetTransaction(repo))
		transaction.PATCH("/:id", handlers.UpdateTransaction(repo))
		transaction.DELETE("/:id", handlers.DeleteTransaction(repo))
	}

	user := router.Group("/user")
	{
		user.POST("/", handlers.CreateUser(repo))
		user.GET("/:id", handlers.GetUser(repo))
		user.PATCH("/:id", handlers.UpdateUserInfo(repo))
		user.DELETE("/:id", handlers.DeleteUser(repo))
	}

	return router
}
