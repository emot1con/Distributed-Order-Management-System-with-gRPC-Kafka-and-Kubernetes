package main

import (
	"fmt"
	"net"
	"net/http"
	"user_service/cmd/api"
	"user_service/cmd/db"
	"user_service/user/usergrpc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	webPort  = "5000"
	gRPCPort = "50001"
)

type UserGRPCServer struct {
	usergrpc.UnimplementedAuthServiceServer
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Info("connect DB")
	gormDB, err := db.Connect(&gin.Context{})
	if err != nil {
		logrus.Fatalf("Error connect DB: %v", err)
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api.Router(gormDB),
	}
	logrus.Info("Starting server on port 5000")

	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatalf("Error listen and serve: %v", err)
	}
}

func gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		logrus.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	usergrpc.RegisterAuthServiceServer(srv, &UserGRPCServer{})
	logrus.Infof("gRPC Server started on port %s", gRPCPort)

	if err := srv.Serve(lis); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}
