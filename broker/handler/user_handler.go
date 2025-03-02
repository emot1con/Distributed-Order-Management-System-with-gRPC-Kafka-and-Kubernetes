package handler

import (
	"broker/repository"
	"broker/usergrpc"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/auth/register", u.Register)
	r.POST("/auth/login", u.Login)
	r.GET("/refresh-token", u.RefreshToken)
}

func (u *UserHandler) Register(c *gin.Context) {
	var registerPayload usergrpc.RegisterPayload
	if err := c.ShouldBindJSON(&registerPayload); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload := &usergrpc.RegisterRequest{
		Payload: &usergrpc.RegisterPayload{
			FullName: registerPayload.FullName,
			Email:    registerPayload.Email,
			Password: registerPayload.Password,
		},
	}

	logrus.Info(payload.Payload.Email)

	logrus.Info("Registering user")
	if err := u.userRepo.Register(payload); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User registered successfully",
	})
}

func (u *UserHandler) Login(c *gin.Context) {
	var loginPayload usergrpc.LoginPayload
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payload := &usergrpc.LoginRequest{
		Payload: &usergrpc.LoginPayload{
			Email:    loginPayload.Email,
			Password: loginPayload.Password,
		},
	}

	token, err := u.userRepo.Login(payload)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	webResponse := &usergrpc.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token.Token,
		ExpiredAt:             token.ExpiredAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiredAt: token.RefreshTokenExpiredAt,
	}

	c.JSON(200, webResponse)
}

func (u *UserHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(400, gin.H{
			"error": "token is required",
		})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	payload := &usergrpc.RefreshTokenRequest{
		Payload: &usergrpc.RefreshTokenPayload{
			RefreshToken: tokenString,
		},
	}

	token, err := u.userRepo.RefreshToken(payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	webResponse := &usergrpc.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token.Token,
		ExpiredAt:             token.ExpiredAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiredAt: token.RefreshTokenExpiredAt,
	}

	c.JSON(200, webResponse)
}
