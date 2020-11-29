package app

import (
	"expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/internal/utils"
	"fmt"
)

func Run() {
	env := LoadEnvironmentVars()

	dbConn, err := repository.NewConnection(env.DBUser, env.DBPassword, env.DBHost, env.DBName)
	if err != nil {
		panic(fmt.Sprintf("couldn't establish postgres connection: %v", err))
	}

	if err := repository.Migrate(dbConn); err != nil {
		panic(fmt.Sprintf("error setting up database: %v", err))
	}

	repository := repository.New(dbConn)
	jwtService := auth.NewJWTService(env.Issuer, env.Secret)
	hasher := utils.NewPasswordHasher()

	r := router.Setup(repository, jwtService, hasher, router.DefaultConfig)
	r.Run(env.Port)
}
