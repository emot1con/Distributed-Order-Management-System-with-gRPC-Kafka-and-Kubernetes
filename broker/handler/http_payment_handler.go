package handler

import (
	"broker/auth"
	"broker/proto"
	"broker/repository"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	repo repository.PaymentRepository
}

func NewPaymentHandler(repo repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		repo: repo,
	}
}

func (u *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	orderRoutes := r.Group("/payment")
	orderRoutes.Use(auth.ProtectedEndpoint())

	orderRoutes.POST("/transaction", u.Transaction)
}

func (u *PaymentHandler) Transaction(c *gin.Context) {
	var payload proto.PaymentTransaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := u.repo.Transaction(&payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "transaction success, your order is being processed"})
}
