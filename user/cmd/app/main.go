package main

import (
	"sync"

	"github.com/sirupsen/logrus"
	"user_service/grpcuser"
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
