package app

import (
	"expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/internal/utils"
	"fmt"

	"github.com/joho/godotenv"
)

func Run() {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("couldn't load env file: %v", err))
	}

	env := NewDefaultEnviroment()
	env.LoadVariables()

	dbConn, err := repository.NewConnection(
		env.DBUser.Value,
		env.DBPassword.Value,
		env.DBHost.Value,
		env.DBName.Value,
		repository.DefaultConfig,
	)
	if err != nil {
		panic(fmt.Sprintf("couldn't establish postgres connection: %v", err))
	}

	if err := repository.Migrate(dbConn); err != nil {
		panic(fmt.Sprintf("error setting up database: %v", err))
	}

	repository := repository.New(dbConn)
	jwtService := auth.NewJWTService(env.Issuer.Value, env.Secret.Value)
	hasher := utils.NewPasswordHasher()

	r := router.Setup(repository, jwtService, hasher, router.DefaultConfig)
	r.Run(env.Port.Value)
}
