package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWTToken(userID int, role string, y, m, d int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(userID),
		"exp":    time.Now().AddDate(y, m, d).Unix(),
		"role":   role,
	})

	tokenEnv := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(tokenEnv))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenString, nil
}
