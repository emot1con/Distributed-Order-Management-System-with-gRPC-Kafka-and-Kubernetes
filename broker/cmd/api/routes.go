package routes

import (
	"broker/handler"
	"broker/repository"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository()
	userHandler := handler.NewUserHandler(userRepo)
	userHandler.RegisterRoutes(r)

	productRepo := repository.NewProductRepository()
	productHandler := handler.NewProductHandler(productRepo)
	productHandler.RegisterRoutes(r)

	orderRepo := repository.NewOrderRepository()
	orderHandler := handler.NewOrderHandler(userRepo, productRepo, orderRepo)
	orderHandler.RegisterRoutes(r)

	paymentRepo := repository.NewPaymentRepository()
	paymentHandler := handler.NewPaymentHandler(paymentRepo)
	paymentHandler.RegisterRoutes(r)

	return r
}
