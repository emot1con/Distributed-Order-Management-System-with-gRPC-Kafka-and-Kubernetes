package handler

import (
	"broker/auth"
	"broker/proto"
	"broker/repository"
	"fmt"
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
			c.JSON(400, gin.H{"error": fmt.Sprintf("Product %s is out of stock", product.Name)})
			return
		}

		for i := range v.Quantity {
			var totalPrice float64
			totalPrice += float64(i+1) * product.Price
			logrus.Infof("product price: %v", product.Price)
			logrus.Infof("total price: %v", totalPrice)
			logrus.Infof("order money: %v", v.Price)
			if v.Price < totalPrice {
				c.JSON(400, gin.H{"error": fmt.Sprintf("Not enough money in %s products (product id: %v)", product.Name, v.ProductId)})
				return
			}
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

	order, err := u.orderRepo.CreateOrder(&payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	orderResponse := &proto.OrderResponse{
		Order: &proto.Order{
			Id:         order.Order.Id,
			UserId:     int32(userID),
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
