package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const issuer = "xpense"

var secret = []byte(os.Getenv("JWT_SECRET"))

var (
	ErrorJWTClaimsInvalid = errors.New("couldn't parse claims")
	ErrorJWTExpired       = errors.New("jwt is expired")
)

type CustomClaims struct {
	Email string `json:"email"`
}

type customClaims struct {
	CustomClaims
	jwt.StandardClaims
}

func CreateJWT(email string) (string, error) {
	claims := customClaims{
		CustomClaims: CustomClaims{
			Email: email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * 2).Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, ErrorJWTClaimsInvalid
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return nil, ErrorJWTExpired
	}

	return &claims.CustomClaims, nil
}
