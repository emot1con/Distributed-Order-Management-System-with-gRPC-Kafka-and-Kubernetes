package main

import (
	"sync"
	// "user_service/cmd/db"
	"user_service/grpcuser"
	// "user_service/repository"
	// "user_service/service"

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
		grpcuser.GRPCListen()
		logrus.Info("GRPCListen() exited")
	}()

	wg.Wait()
}
