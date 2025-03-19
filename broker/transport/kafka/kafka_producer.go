package kafka

// import (
// 	"broker/proto"
// 	"encoding/json"

// 	"github.com/IBM/sarama"
// )

// func ConnectProducer(addr []string) (sarama.SyncProducer, error) {
// 	config := sarama.NewConfig()
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Retry.Max = 5
// 	config.Producer.Return.Successes = true

// 	return sarama.NewSyncProducer(addr, config)
// }

// func SendMessage(addr []string, topic string, msg *proto.Order) (int32, int64, error) {
// 	producer, err := ConnectProducer(addr)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	defer producer.Close()

// 	data, err := json.Marshal(msg)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	return producer.SendMessage(&sarama.ProducerMessage{
// 		Topic: topic,
// 		Value: sarama.ByteEncoder(data),
// 	})
// }
