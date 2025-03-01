package main

import (
	// "fmt"
	// "net"
	// "net/http"
	// "user_service/cmd/api"
	// "user_service/cmd/db"

	// "github.com/gin-gonic/gin"
	"user_service/grpcuser"

	"github.com/sirupsen/logrus"
	// "google.golang.org/grpc"
)

// const (
// 	webPort  = "5000"
// 	gRPCPort = "50001"
// )

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Info("connect DB")
	// gormDB, err := db.Connect(&gin.Context{})
	// if err != nil {
	// 	logrus.Fatalf("Error connect DB: %v", err)
	// }

	go grpcuser.GRPCListen()
}
