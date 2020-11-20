package router

import (
	"expense-api/handlers"
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	party_middleware "expense-api/middleware/party"
	transaction_middleware "expense-api/middleware/transaction"
	wallet_middleware "expense-api/middleware/wallet"
	"expense-api/repository"
	"expense-api/utils"

	"github.com/gin-gonic/gin"
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
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	}

	handler := handlers.New(repo, jwtService, hasher)

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

	wallet := router.Group("/wallet").Use(authM.IsAuthenticated)
	{
		walletM := wallet_middleware.New(repo)

		wallet.GET("/", handler.ListWallets)
		wallet.POST("/", handler.CreateWallet)
		wallet.GET("/:id", commonM.SetIDParamToContext, walletM.ValidateOwnership, handler.GetWallet)
		wallet.PATCH("/:id", commonM.SetIDParamToContext, walletM.ValidateOwnership, handler.UpdateWallet)
		wallet.DELETE("/:id", commonM.SetIDParamToContext, walletM.ValidateOwnership, handler.DeleteWallet)
		wallet.GET("/:id/transaction", commonM.SetIDParamToContext, walletM.ValidateOwnership, handler.ListTransactionsByWallet)
	}

	party := router.Group("/party").Use(authM.IsAuthenticated)
	{
		partyM := party_middleware.New(repo)

		party.GET("/", handler.ListParties)
		party.POST("/", handler.CreateParty)
		party.GET("/:id", commonM.SetIDParamToContext, partyM.ValidateOwnership, handler.GetParty)
		party.PATCH("/:id", commonM.SetIDParamToContext, partyM.ValidateOwnership, handler.UpdateParty)
		party.DELETE("/:id", commonM.SetIDParamToContext, partyM.ValidateOwnership, handler.DeleteParty)
		party.GET("/:id/transaction", commonM.SetIDParamToContext, partyM.ValidateOwnership, handler.ListTransactionsByParty)
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
