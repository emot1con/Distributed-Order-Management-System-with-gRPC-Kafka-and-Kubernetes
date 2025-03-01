package api

import (
	"user_service/controller"
	"user_service/repository"
	"user_service/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(DB *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(DB)
	userService := service.NewUserService(DB, userRepo)
	userController := controller.NewUserController(userService)
	userController.RegisterRoutes(r)

	return r
}
