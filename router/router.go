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

	account := router.Group("/account")
	{
		account.GET("/:id", authMiddleware.Handler, handler.GetAccount)
		account.PATCH("/:id", authMiddleware.Handler, handler.UpdateAccount)
		account.DELETE("/:id", authMiddleware.Handler, handler.DeleteAccount)
	}

	transaction := router.Group("/transaction")
	{
		transaction.GET("/", authMiddleware.Handler, handler.ListTransactions)
		transaction.POST("/", authMiddleware.Handler, handler.CreateTransaction)
		transaction.GET("/:id", authMiddleware.Handler, handler.GetTransaction)
		transaction.PATCH("/:id", authMiddleware.Handler, handler.UpdateTransaction)
		transaction.DELETE("/:id", authMiddleware.Handler, handler.DeleteTransaction)
	}

	return router
}
