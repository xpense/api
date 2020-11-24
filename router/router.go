package router

import (
	"expense-api/handlers"
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	parties_middleware "expense-api/middleware/parties"
	transactions_middleware "expense-api/middleware/transactions"
	wallets_middleware "expense-api/middleware/wallets"
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

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/signup", handler.SignUp)
		auth.POST("/login", handler.Login)
	}

	commonM := middleware.NewCommonMiddleware()
	authM := auth_middleware.New(jwtService)

	account := v1.Group("/account").Use(authM.IsAuthenticated)
	{
		account.GET("/", handler.GetAccount)
		account.PATCH("/", handler.UpdateAccount)
		account.DELETE("/", handler.DeleteAccount)
	}

	wallets := v1.Group("/wallets").Use(authM.IsAuthenticated)
	{
		walletsM := wallets_middleware.New(repo)

		wallets.GET("/", handler.ListWallets)
		wallets.POST("/", handler.CreateWallet)
		wallets.GET("/:id", commonM.SetIDParamToContext, walletsM.ValidateOwnership, handler.GetWallet)
		wallets.PATCH("/:id", commonM.SetIDParamToContext, walletsM.ValidateOwnership, handler.UpdateWallet)
		wallets.DELETE("/:id", commonM.SetIDParamToContext, walletsM.ValidateOwnership, handler.DeleteWallet)
		wallets.GET("/:id/transactions", commonM.SetIDParamToContext, walletsM.ValidateOwnership, handler.ListTransactionsByWallet)
	}

	parties := v1.Group("/parties").Use(authM.IsAuthenticated)
	{
		partiesM := parties_middleware.New(repo)

		parties.GET("/", handler.ListParties)
		parties.POST("/", handler.CreateParty)
		parties.GET("/:id", commonM.SetIDParamToContext, partiesM.ValidateOwnership, handler.GetParty)
		parties.PATCH("/:id", commonM.SetIDParamToContext, partiesM.ValidateOwnership, handler.UpdateParty)
		parties.DELETE("/:id", commonM.SetIDParamToContext, partiesM.ValidateOwnership, handler.DeleteParty)
		parties.GET("/:id/transactions", commonM.SetIDParamToContext, partiesM.ValidateOwnership, handler.ListTransactionsByParty)
	}

	transactions := v1.Group("/transactions").Use(authM.IsAuthenticated)
	{
		txM := transactions_middleware.New(repo)

		transactions.GET("/", handler.ListTransactions)
		transactions.POST("/", handler.CreateTransaction)
		transactions.GET("/:id", commonM.SetIDParamToContext, txM.ValidateOwnership, handler.GetTransaction)
		transactions.PATCH("/:id", commonM.SetIDParamToContext, txM.ValidateOwnership, handler.UpdateTransaction)
		transactions.DELETE("/:id", commonM.SetIDParamToContext, txM.ValidateOwnership, handler.DeleteTransaction)
	}

	return router
}
