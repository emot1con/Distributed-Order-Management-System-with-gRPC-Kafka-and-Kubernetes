package main

import (
	"net/http"
	"sync"

	"user_service/transport/grpc"
	transport "user_service/transport/http"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting Application")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logrus.Info("Calling GRPCListen()")
		grpc.GRPCListen()
		logrus.Info("GRPCListen() exited")
	}()

	srv := http.Server{
		Addr:    ":50001",
		Handler: transport.Routes(),
	}

	logrus.Info("Starting HTTP server on port 50001")
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatalf("Failed to start HTTP server: %v", err)
	}

	wg.Wait()
}
