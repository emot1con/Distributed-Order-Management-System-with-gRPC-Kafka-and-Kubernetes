package main

import (
	"os"
	"payment/transport/kafka"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting Payment Application")

	addr := []string{os.Getenv("KAFKA_BROKER_URL")}
	topic := os.Getenv("KAFKA_ORDER_TOPIC")

	consumer, err := kafka.NewConsumer(addr, topic)
	if err != nil {
		logrus.Fatal(err)
	}
	defer consumer.Close()

	go kafka.ProcessMessage(consumer)
	logrus.Info("Application started")
	select {}
}
