package main

import (
	"os"
	"payment/transport/grpc"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting Payment Application")

	addr := []string{os.Getenv("KAFKA_BROKER_URL")}
	topic := os.Getenv("KAFKA_ORDER_TOPIC")

	go grpc.GRPCListen(addr, topic)

	logrus.Info("Application started")
	select {}
}
