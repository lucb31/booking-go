package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

// Todo Should put this key somewhere else
var jwtSecretKey = []byte("generateMeSecretlyAndStoreInEnv")

func LoginRequest(username string, password string) (string, error) {
	if username != "root" || password != "root" {
		return "", errors.New("Invalid credentials")
	}
	return GenerateJWT(username)
}

func GenerateJWT(username string) (string, error) {
	// Define claims
	claims := JwtClaims{"foo",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   username,
		}}

	// Initialize token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token string
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	// Empty token provided
	if len(tokenString) == 0 {
		return nil, errors.New("No token provided")
	}
	// Parse token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Incompatible signing method")
		}

		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
