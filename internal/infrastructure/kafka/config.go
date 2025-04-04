package kafka

import (
	"github.com/IBM/sarama"
)

type KafkaSettings struct {
	Consumer      sarama.Consumer
	Addres        []string
	Topic         string
	Message       <-chan *sarama.ConsumerMessage
	PartitionCons sarama.PartitionConsumer
}

func KafkaGetConfig(arr []string, s string) *KafkaSettings {
	return &KafkaSettings{
		Addres: arr,
		Topic:  s,
	}
}
