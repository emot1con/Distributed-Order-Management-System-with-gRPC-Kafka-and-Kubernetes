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
	topic := []string{os.Getenv("KAFKA_ORDER_TOPIC")}
	groupID := os.Getenv("KAFKA_GROUP_ID")

	go grpc.GRPCListen(addr, topic, groupID)

	logrus.Info("Application started")
	select {}
}
