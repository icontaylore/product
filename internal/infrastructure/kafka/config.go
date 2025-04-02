package kafka

import (
	"github.com/IBM/sarama"
	"product/internal/infrastructure/elasticsearch"
)

type KafkaSettings struct {
	ElasticService *elasticsearch.Service
	Consumer       sarama.Consumer
	Worker         uint
	Topic          string
	Host           string
	MessageChan    <-chan *sarama.ConsumerMessage
	PartitionCons  sarama.PartitionConsumer // контроль закрытия
}

func (k *KafkaSettings) SetElasticService(service *elasticsearch.Service) {
	k.ElasticService = service
}
func GoConfigure(top, host string, worker uint) *KafkaSettings {
	return &KafkaSettings{
		Worker:      worker,
		Topic:       top,
		Host:        host,
		MessageChan: make(chan *sarama.ConsumerMessage),
	}
}
