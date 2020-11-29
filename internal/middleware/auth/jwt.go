package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	CreateJWT(id uint, email string) (string, error)
	ValidateJWT(tokenString string) (*CustomClaims, error)
}

var (
	ErrorJWTClaimsInvalid = errors.New("couldn't parse claims")
	ErrorJWTExpired       = errors.New("jwt is expired")
)

type jwtClaims struct {
	CustomClaims
	jwt.StandardClaims
}

type jwtService struct {
	issuer string
	secret []byte
}

func NewJWTService(issuer, secret string) JWTService {
	return &jwtService{
		issuer: issuer,
		secret: []byte(secret),
	}
}

func (jwts *jwtService) CreateJWT(id uint, email string) (string, error) {
	claims := jwtClaims{
		CustomClaims: CustomClaims{
			ID:    id,
			Email: email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * 60 * 24 * 365).Unix(),
			Issuer:    jwts.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwts.secret)
}

func (jwts *jwtService) ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) { return jwts.secret, nil })
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, ErrorJWTClaimsInvalid
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return nil, ErrorJWTExpired
	}

	return &claims.CustomClaims, nil
}
