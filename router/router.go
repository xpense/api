package router

import (
	"expense-api/handlers"
	"expense-api/repository"
	"expense-api/utils"

	"github.com/gin-gonic/gin"
)

// Setup creates a new gin router
func Setup(repo repository.Repository, hasher utils.PasswordHasher) *gin.Engine {
	router := gin.Default()
	handler := handlers.New(repo, hasher)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", handler.SignUp)
		auth.POST("/login", handler.Login)
		// auth.POST("/change-password", handlers.ChangePassword(repo, hasher))
	}

	transaction := router.Group("/transaction")
	{
		transaction.GET("/", handler.ListTransactions)
		transaction.POST("/", handler.CreateTransaction)
		transaction.GET("/:id", handler.GetTransaction)
		transaction.PATCH("/:id", handler.UpdateTransaction)
		transaction.DELETE("/:id", handler.DeleteTransaction)
	}

	user := router.Group("/user")
	{
		user.GET("/:id", handler.GetUser)
		user.PATCH("/:id", handler.UpdateUserInfo)
		user.DELETE("/:id", handler.DeleteUser)
	}

	return router
}
