package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cleoexcel/ristek-test/app/models"
	"github.com/cleoexcel/ristek-test/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAuthTokens(user *models.User) (string, error) {
	now := time.Now()

	accessToken, err := createJWT(user.ID, now, config.JWT_EXPIRY_IN_DAY)
	if err != nil {
		return "", fmt.Errorf("error creating access token: %w", err)
	}

	return accessToken, nil
}

func createJWT(userId int, issuedAt time.Time, expiryDays int) (string, error) {
	expiresAt := issuedAt.Add(time.Hour * 24 * time.Duration(expiryDays))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   strconv.Itoa(int(userId)),
		},
	})
	return token.SignedString([]byte(config.JWT_SECRET_KEY))
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func RetrieveTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	sections := strings.Split(authHeader, " ")
	if len(sections) != 2 || sections[0] != "Bearer" {
		return "", errors.New("invalid authorization header structure")
	}
	return sections[1], nil
}

func DecodeJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method)
		}
		return []byte(config.JWT_SECRET_KEY), nil
	})
}
