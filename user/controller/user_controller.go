package controller

import (
	"time"
	"user_service/service"
	"user_service/types"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (u *UserController) RegisterRoutes(router *gin.Engine) {
	router.POST("/auth/register", u.Register)
	router.POST("/auth/login", u.Login)
	router.GET("/refresh-token", u.RefreshToken)
}

func (u *UserController) Register(c *gin.Context) {
	var payload types.RegisterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := u.service.Register(payload); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}

func (u *UserController) Login(c *gin.Context) {
	var payload types.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		logrus.Error(err)
		return
	}

	token, refreshToken, err := u.service.Login(payload)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		logrus.Error(err)
		return
	}

	response := types.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}

	c.JSON(200, response)
}

func (u *UserController) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(401, gin.H{"error": "refresh token missing"})
		return
	}

	newToken, newRefreshToken, err := u.service.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := types.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 newToken,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}

	c.JSON(200, response)
}
