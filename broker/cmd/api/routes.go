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

	return r
}
