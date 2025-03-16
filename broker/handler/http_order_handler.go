package handler

import (
	"broker/auth"
	"broker/proto"
	"broker/repository"
	"broker/transport/kafka"
	"fmt"
	"os"
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
	var totalPrices float64
	userID, ok := c.Request.Context().Value(auth.UserKey).(int)
	if !ok {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	var payload proto.CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, v := range payload.Items {
		product, err := u.productRepo.GetProduct(&proto.GetProductRequest{Id: v.ProductId})
		if err != nil {
			c.JSON(400, gin.H{"error": "Product not found"})
			return
		}
		if product.Stock < v.Quantity {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Product %s is out of stock (product id: %v)", product.Name, v.ProductId)})
			return
		}

		for i := range v.Quantity {
			var totalPrice float64
			totalPrice += float64(i+1-(i)) * product.Price
			if v.Price < totalPrice {
				c.JSON(400, gin.H{"error": fmt.Sprintf("Not enough money in %s products (product id: %v)", product.Name, v.ProductId)})
				return
			}
			totalPrices += totalPrice
		}

		if _, err := u.productRepo.UpdateProduct(&proto.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock - v.Quantity,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	payload.TotalPrice = totalPrices
	payload.UserId = int32(userID)

	order, err := u.orderRepo.CreateOrder(&payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	addr := []string{os.Getenv("KAFKA_BROKER_URL")}
	topic := os.Getenv("KAFKA_ORDER_TOPIC")
	logrus.Printf("address: %s", addr)

	partition, offset, err := kafka.SendMessage(addr, "order-topic", &proto.Order{
		Id:         order.Order.Id,
		UserId:     int32(userID),
		Status:     "Pending",
		TotalPrice: order.Order.TotalPrice,
		OrderItems: order.Order.OrderItems,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Message sent to topic: %s partition: %d, offset: %d", topic, partition, offset)

	orderResponse := &proto.OrderResponse{
		Order: &proto.Order{
			Id:         order.Order.Id,
			UserId:     int32(userID),
			Status:     "Pending",
			TotalPrice: order.Order.TotalPrice,
			CreatedAt:  order.Order.CreatedAt,
			UpdatedAt:  order.Order.CreatedAt,
			OrderItems: order.Order.OrderItems,
		},
	}

	c.JSON(200, orderResponse)
}

func (u *OrderHandler) GetOrder(c *gin.Context) {
	query := c.Query("order_id")
	orderID, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := u.orderRepo.GetOrder(&proto.GetOrderRequest{
		OrderId: int32(orderID),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}
