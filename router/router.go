package router

import (
	"expense-api/handlers"
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	transaction_middleware "expense-api/middleware/transaction"
	"expense-api/repository"
	"expense-api/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Config struct {
	withDefaultMiddleware bool
}

var DefaultConfig = &Config{withDefaultMiddleware: true}
var TestConfig = &Config{}

// Setup creates a new gin router
func Setup(
	repo repository.Repository,
	jwtService auth_middleware.JWTService,
	hasher utils.PasswordHasher,
	config *Config,
) *gin.Engine {

	var router *gin.Engine
	if config.withDefaultMiddleware {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	handler := handlers.New(repo, jwtService, hasher)

	router.POST("/", func(ctx *gin.Context) {
		var transaction struct {
			Amount decimal.Decimal `json:"amount"`
		}

		if err := ctx.Bind(&transaction); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		fmt.Println("Amount - ", transaction.Amount)
	})

	auth := router.Group("/auth")
	{
		auth.POST("/signup", handler.SignUp)
		auth.POST("/login", handler.Login)
	}

	commonM := middleware.NewCommonMiddleware()
	authM := auth_middleware.New(jwtService)

	account := router.Group("/account").Use(authM.IsAuthenticated)
	{
		account.GET("/", handler.GetAccount)
		account.PATCH("/", handler.UpdateAccount)
		account.DELETE("/", handler.DeleteAccount)
	}

	transaction := router.Group("/transaction").Use(authM.IsAuthenticated)
	{
		txM := transaction_middleware.New(repo)

		transaction.GET("/", handler.ListTransactions)
		transaction.POST("/", handler.CreateTransaction)
		transaction.GET("/:id", commonM.SetIDParamToContext, txM.ValidateOwnership, handler.GetTransaction)
		transaction.PATCH("/:id", commonM.SetIDParamToContext, txM.ValidateOwnership, handler.UpdateTransaction)
		transaction.DELETE("/:id", commonM.SetIDParamToContext, txM.ValidateOwnership, handler.DeleteTransaction)
	}

	return router
}
