package router

import (
	"expense-api/handlers"
	"expense-api/middleware"
	"expense-api/repository"
	"expense-api/utils"

	"github.com/gin-gonic/gin"
)

// Setup creates a new gin router
func Setup(
	repo repository.Repository,
	jwtService utils.JWTService,
	hasher utils.PasswordHasher,
) *gin.Engine {
	router := gin.Default()

	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	handler := handlers.New(repo, jwtService, hasher)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", handler.SignUp)
		auth.POST("/login", handler.Login)
	}

	account := router.Group("/account").Use(authMiddleware.Handler)
	{
		account.GET("", handler.GetAccount)
		account.GET("/", handler.GetAccount)

		account.PATCH("", handler.UpdateAccount)
		account.PATCH("/", handler.UpdateAccount)

		account.DELETE("", handler.DeleteAccount)
		account.DELETE("/", handler.DeleteAccount)
	}

	transaction := router.Group("/transaction").Use(authMiddleware.Handler)
	{
		transaction.GET("/", handler.ListTransactions)
		transaction.POST("/", handler.CreateTransaction)
		transaction.GET("/:id", handler.GetTransaction)
		transaction.PATCH("/:id", handler.UpdateTransaction)
		transaction.DELETE("/:id", handler.DeleteTransaction)
	}

	return router
}
