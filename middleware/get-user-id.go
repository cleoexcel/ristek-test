package middleware

import (

	"net/http"

	"github.com/cleoexcel/ristek-test/app/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

var userIDContextKey = ContextKey("user-id")

func ExtractUserID(c *gin.Context) string {
	if userID, exists := c.Get(string(userIDContextKey)); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func GetUserIDFromToken(c *gin.Context) (string, error) {
	tokenString, err := auth.RetrieveTokenFromHeader(c.Request)
	if err != nil {
		return "", err
	}


	token, err := auth.DecodeJWT(tokenString)
	if err != nil {
		return "", err
	}


	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", http.ErrNoCookie
	}


	userID, ok := claims["sub"].(string)
	if !ok {
		return "", http.ErrNoCookie
	}


	return userID, nil
}