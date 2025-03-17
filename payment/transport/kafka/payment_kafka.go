package kafka

import (
	"encoding/json"
	"os"
	"os/signal"
	"payment/proto"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func connectKafka(addr []string) (sarama.Consumer, error) {
	topic := sarama.NewConfig()
	topic.Consumer.Return.Errors = true

	return sarama.NewConsumer(addr, topic)
}

func NewConsumer(addr []string, topic string) (sarama.PartitionConsumer, error) {
	worker, err := connectKafka(addr)
	if err != nil {
		return nil, err
	}

	return worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
}

func ProcessMessage(consumer sarama.PartitionConsumer) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	msgCnt := 0

	donech := make(chan struct{})

	payment := new(proto.OrderPayment)
	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				if err := json.Unmarshal(msg.Value, &payment); err != nil {
					logrus.Errorf("Error parsing message: %s. Raw data: %s", err, string(msg.Value))
					continue
				}
				msgCnt++
				logrus.Infof("Received message, UserID: %d with OrderId: %d \n", payment.UserId, payment.Id)
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
