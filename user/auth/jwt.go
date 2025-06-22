package auth

import (
	"os"
	"strconv"
	"time"
	"user_service/proto"

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

func CreateFullJWTToken(ID int) (*proto.TokenResponse, error) {
	accesToken, err := CreateJWTToken(ID, "user", 0, 0, 1)
	if err != nil {
		return nil, err
	}
	refreshToken, err := CreateJWTToken(ID, "user", 0, 3, 0)
	if err != nil {
		return nil, err
	}

	return &proto.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 accesToken,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}, nil

}
