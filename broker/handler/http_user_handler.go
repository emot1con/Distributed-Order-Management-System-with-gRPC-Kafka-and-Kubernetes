package handler

import (
	"broker/proto"
	"broker/repository"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type UserHandler struct {
	userRepo     repository.UserRepository
	userLimiters map[string]*rate.Limiter
	mu           sync.Mutex
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		userLimiters: make(map[string]*rate.Limiter),
	}
}

func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1/auth")
	{
		apiV1.POST("/register", u.Register)
		apiV1.POST("/login", u.Login)
		apiV1.GET("/refresh-token", u.RefreshToken)
		apiV1.GET("/google", u.GoogleOauth)
		apiV1.GET("/facebook", u.FacebookOauth)
		apiV1.GET("/github", u.GithubOauth)
	}
	apiV1Oauth := r.Group("/api/v1/oauth")
	{
		apiV1Oauth.GET("/google/callback", u.GoogleOauthCallback)
		apiV1Oauth.GET("/facebook/callback", u.FacebookOauthCallback)
		apiV1Oauth.GET("/github/callback", u.GithubOauthCallback)
	}
}

func (u *UserHandler) getLimiter(userID string) *rate.Limiter {
	u.mu.Lock()
	defer u.mu.Unlock()

	limiter, exists := u.userLimiters[userID]
	if !exists {
		limiter = rate.NewLimiter(2, 5)
		u.userLimiters[userID] = limiter
	}
	return limiter
}

func (u *UserHandler) Register(c *gin.Context) {
	var registerPayload proto.RegisterPayload
	if err := c.ShouldBindJSON(&registerPayload); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userLimiter := u.getLimiter(registerPayload.Email)
	if !userLimiter.Allow() {
		c.JSON(429, gin.H{"error": "Too many register attempts"})
		logrus.Warnf("Rate limit exceeded for user %s", registerPayload.Email)
		return
	}

	payload := &proto.RegisterRequest{
		Payload: &proto.RegisterPayload{
			FullName: registerPayload.FullName,
			Email:    registerPayload.Email,
			Password: registerPayload.Password,
		},
	}

	logrus.Infof("Registering user: %s", payload.Payload.Email)

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
	var loginPayload proto.LoginPayload
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userLimiter := u.getLimiter(loginPayload.Email)
	if !userLimiter.Allow() {
		c.JSON(429, gin.H{"error": "Too many login attempts"})
		logrus.Warnf("Rate limit exceeded for user %s", loginPayload.Email)
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

func (u *UserHandler) GoogleOauth(c *gin.Context) {
	url, err := u.userRepo.GoogleOauth()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(302, url.Url)
}

func (u *UserHandler) FacebookOauth(c *gin.Context) {
	url, err := u.userRepo.FacebookOauth()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(302, url.Url)
}

func (u *UserHandler) GithubOauth(c *gin.Context) {
	url, err := u.userRepo.GithubOauth()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(302, url.Url)
}

func (u *UserHandler) GoogleOauthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "code is required"})
		return
	}

	token, err := u.userRepo.GoogleOauthCallback(code)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}

func (u *UserHandler) FacebookOauthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "code is required"})
		return
	}

	token, err := u.userRepo.FacebookOauthCallback(code)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}

func (u *UserHandler) GithubOauthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "code is required"})
		return
	}

	token, err := u.userRepo.GithubOauthCallback(code)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}
