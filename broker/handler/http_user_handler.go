package handler

import (
	"broker/proto"
	"broker/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
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
	var registerPayload proto.RegisterPayload
	if err := c.ShouldBindJSON(&registerPayload); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload := &proto.RegisterRequest{
		Payload: &proto.RegisterPayload{
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

var limiter = rate.NewLimiter(2, 5)

func (u *UserHandler) Login(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(429, gin.H{"error": "Too many requests"})
		logrus.Errorf("Too many requests")
		c.Abort()
		return
	}
	c.Next()

	var loginPayload proto.LoginPayload
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payload := &proto.LoginRequest{
		Payload: &proto.LoginPayload{
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

	webResponse := &proto.TokenResponse{
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

	payload := &proto.RefreshTokenRequest{
		Payload: &proto.RefreshTokenPayload{
			RefreshToken: tokenString,
		},
	}

	token, err := u.userRepo.RefreshToken(payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	webResponse := &proto.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token.Token,
		ExpiredAt:             token.ExpiredAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiredAt: token.RefreshTokenExpiredAt,
	}

	c.JSON(200, webResponse)
}
