package kafka

import "github.com/IBM/sarama"

type KafkaSettings struct {
	Consumer      sarama.Consumer
	Worker        uint
	Topic         string
	Host          string
	MessageChan   <-chan *sarama.ConsumerMessage
	PartitionCons sarama.PartitionConsumer // контроль закрытия
}

func GoConfigure(top, host string, worker uint) *KafkaSettings {
	return &KafkaSettings{
		Worker:      worker,
		Topic:       top,
		Host:        host,
		MessageChan: make(chan *sarama.ConsumerMessage),
	}
}
