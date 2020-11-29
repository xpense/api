package app

import (
	"fmt"
	"os"
)

const (
	PORT        = "PORT"
	JWT_ISSUER  = "JWT_ISSUER"
	JWT_SECRET  = "JWT_SECRET"
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_HOST     = "DB_HOST"
	DB_NAME     = "DB_NAME"
)

type (
	EnvironmentVariable struct {
		Name  string
		Value string
	}

	Environment struct {
		Port       EnvironmentVariable
		Issuer     EnvironmentVariable
		Secret     EnvironmentVariable
		DBUser     EnvironmentVariable
		DBPassword EnvironmentVariable
		DBHost     EnvironmentVariable
		DBName     EnvironmentVariable
	}
)

func NewDefaultEnviroment() *Environment {
	return &Environment{
		Port:       EnvironmentVariable{Name: PORT},
		Issuer:     EnvironmentVariable{Name: JWT_ISSUER},
		Secret:     EnvironmentVariable{Name: JWT_SECRET},
		DBUser:     EnvironmentVariable{Name: DB_USER},
		DBPassword: EnvironmentVariable{Name: DB_PASSWORD},
		DBHost:     EnvironmentVariable{Name: DB_HOST},
		DBName:     EnvironmentVariable{Name: DB_NAME},
	}
}

func (e *Environment) LoadVariables() {
	if port := os.Getenv(e.Port.Name); port == "" {
		e.Port.Value = ":8080"
	} else {
		e.Port.Value = ":" + port
	}

	e.Issuer.Value = os.Getenv(e.Issuer.Name)
	assertEnvVarSet(e.Issuer)

	e.Secret.Value = os.Getenv(e.Secret.Name)
	assertEnvVarSet(e.Secret)

	e.DBUser.Value = os.Getenv(e.DBUser.Name)
	assertEnvVarSet(e.DBUser)

	e.DBPassword.Value = os.Getenv(e.DBPassword.Name)
	assertEnvVarSet(e.DBPassword)

	e.DBHost.Value = os.Getenv(e.DBHost.Name)
	assertEnvVarSet(e.DBHost)

	e.DBName.Value = os.Getenv(e.DBName.Name)
	assertEnvVarSet(e.DBName)

}

func assertEnvVarSet(envVar EnvironmentVariable) {
	if envVar.Value == "" {
		panic(fmt.Sprintf("%s not set!", envVar.Name))
	}
}
