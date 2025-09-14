package kafka

import (
	"encoding/json"
	"errors"
	"order/proto"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var producer sarama.SyncProducer

func ConnectProducer(addr []string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		return err
	}
	producer = p

	logrus.Info("success connect producer")

	return nil
}

func SendMessage(addr []string, topic string, msg *proto.Order) (int32, int64, error) {
	if producer == nil {
		return 0, 0, errors.New("kafka producer is not initialized")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return 0, 0, err
	}

	return producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
}
