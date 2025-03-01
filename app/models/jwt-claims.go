package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	TokenType string `json:"type"`
	jwt.RegisteredClaims
}

