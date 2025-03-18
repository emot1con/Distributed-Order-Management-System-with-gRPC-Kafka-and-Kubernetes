package kafka

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"payment/proto"
	"payment/service"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func connectKafka(addr []string) (sarama.Consumer, error) {
	topic := sarama.NewConfig()
	topic.Consumer.Return.Errors = true

	for i := 0; i < 5; i++ {
		worker, err := sarama.NewConsumer(addr, topic)
		if err == nil {
			logrus.Info("Connected to kafka")
			return worker, nil
		}
		logrus.Warnf("failed connect to kafka, retrying...(%d/5)", i+1)
		time.Sleep(5 * time.Second)
		continue
	}

	return nil, errors.New("failed to connect to kafka")
}

func ProcessMessage(addr []string, topic string, service *service.PaymentService) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	msgCnt := 0

	donech := make(chan struct{})

	worker, err := connectKafka(addr)
	if err != nil {
		logrus.Fatalf("error when connect to kafka: %v", err)
	}

	defer worker.Close()

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logrus.Fatalf("error when consume partition: %v", err)
	}
	defer consumer.Close()

	order := new(proto.Order)
	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					logrus.Errorf("Error parsing message: %s. Raw data: %s", err, string(msg.Value))
					continue
				}
				msgCnt++
				logrus.Infof("Received message, UserID: %d with OrderId: %d \n", order.UserId, order.Id)
				response, err := service.AddPayment(&proto.CreatePaymentRequest{
					OrderId:    order.Id,
					UserId:     order.UserId,
					TotalPrice: order.TotalPrice,
				})
				if err != nil {
					logrus.Error(err)
					continue
				}
				logrus.Infof("Payment created with ID: %d\n", response.Id)
			case err := <-consumer.Errors():
				logrus.Error(err)
			case <-sigchan:
				logrus.Info("Consumer stopped")
				donech <- struct{}{}
				return
			}
		}
	}()
	<-donech
	logrus.Infof("Processed %d messages\n", msgCnt)
}
