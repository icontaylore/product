package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

func (k *KafkaSettings) SubscribeTopic() error {
	partitionCons, err := k.Consumer.ConsumePartition(k.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("main:невозможно организовать партицию")
	}
	k.MessageChan = partitionCons.Messages()
	k.PartitionCons = partitionCons

	return nil
}

func (k *KafkaSettings) CreateKafkaConsumer() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{k.Host}, config)
	if err != nil {
		defer log.Fatal("kafka:ошибка при создании Kafka consumer:", err)
		return err
	}
	k.Consumer = consumer
	return nil
}
