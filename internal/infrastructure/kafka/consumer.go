package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

func (k *KafkaSettings) NewConsumer() error {
	cons, err := sarama.NewConsumer(k.Addres, sarama.NewConfig())
	if err != nil {
		return err
	}
	k.Consumer = cons
	return nil
}

func (k *KafkaSettings) SubscribeTopic() error {
	partitionCons, err := k.Consumer.ConsumePartition(k.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("main:невозможно организовать партицию")
	}
	k.Message = partitionCons.Messages()
	k.PartitionCons = partitionCons

	return nil
}
