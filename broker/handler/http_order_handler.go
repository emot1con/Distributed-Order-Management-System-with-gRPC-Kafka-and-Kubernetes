package handler

import (
	"broker/auth"
	"broker/proto"
	"broker/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
}

func NewOrderHandler(userRepo repository.UserRepository, productRepo repository.ProductRepository, orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{
		userRepo:    userRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

func (u *OrderHandler) RegisterRoutes(r *gin.Engine) {
	orderRoutes := r.Group("/order")
	orderRoutes.Use(auth.ProtectedEndpoint())

	orderRoutes.POST("/", u.CreateOrder)
	orderRoutes.GET("/", u.GetOrder)
}

func (u *OrderHandler) CreateOrder(c *gin.Context) {
	userID, ok := c.Request.Context().Value(auth.UserKey).(int)
	if !ok {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	logrus.Infof("Creating payload for user ID: %d", userID)
	var payload proto.CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	payload.UserId = int32(userID)

	logrus.Infof("Creating order for user ID: %d", userID)
	order, err := u.orderRepo.CreateOrder(&payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	orderResponse := &proto.OrderResponse{
		Order: &proto.Order{
			Id:         order.Order.Id,
			UserId:     int32(userID),
			Status:     "Pending",
			TotalPrice: order.Order.TotalPrice,
			CreatedAt:  order.Order.CreatedAt,
			UpdatedAt:  order.Order.CreatedAt,
		},
	}

	c.JSON(200, orderResponse)
}

func (u *OrderHandler) GetOrder(c *gin.Context) {
	logrus.Info("Fetching order by ID")
	query := c.Query("order_id")
	orderID, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	logrus.Infof("calling repo")
	order, err := u.orderRepo.GetOrder(&proto.GetOrderRequest{
		OrderId: int32(orderID),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}
