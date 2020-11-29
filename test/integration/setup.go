package integration

import (
	"expense-api/internal/app"
	"expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/internal/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	TEST_JWT_ISSUER  = "TEST_JWT_ISSUER"
	TEST_JWT_SECRET  = "TEST_JWT_SECRET"
	TEST_DB_USER     = "TEST_DB_USER"
	TEST_DB_PASSWORD = "TEST_DB_PASSWORD"
	TEST_DB_HOST     = "TEST_DB_HOST"
	TEST_DB_NAME     = "TEST_DB_NAME"
)

func NewTestingEnvironment() *app.Environment {
	return &app.Environment{
		Issuer:     app.EnvironmentVariable{Name: TEST_JWT_ISSUER},
		Secret:     app.EnvironmentVariable{Name: TEST_JWT_SECRET},
		DBUser:     app.EnvironmentVariable{Name: TEST_DB_USER},
		DBPassword: app.EnvironmentVariable{Name: TEST_DB_PASSWORD},
		DBHost:     app.EnvironmentVariable{Name: TEST_DB_HOST},
		DBName:     app.EnvironmentVariable{Name: TEST_DB_NAME},
	}
}

func Setup() *gin.Engine {
	env := NewTestingEnvironment()
	env.LoadVariables()

	dbConn, err := repository.NewConnection(env.DBUser.Value, env.DBPassword.Value, env.DBHost.Value, env.DBName.Value, nil)
	if err != nil {
		panic(fmt.Sprintf("couldn't establish postgres connection: %v", err))
	}

	if err := repository.Cleanup(dbConn); err != nil {
		panic(fmt.Sprintf("error cleaning up database: %v", err))
	}

	if err := repository.Migrate(dbConn); err != nil {
		panic(fmt.Sprintf("error setting up database: %v", err))
	}

	repository := repository.New(dbConn)
	jwtService := auth.NewJWTService(env.Issuer.Value, env.Secret.Value)
	hasher := utils.NewPasswordHasher()

	return router.Setup(repository, jwtService, hasher, router.TestConfig)
}
