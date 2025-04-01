package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"sync"
)

func CreateKafkaConsumer() sarama.Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9094"}, config)
	if err != nil {
		log.Fatal("Ошибка при создании Kafka consumer:", err)
	}

	return consumer
}

func Worker(id int, messageChan chan *sarama.ConsumerMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	for message := range messageChan {
		fmt.Printf("Worker %d обрабатывает сообщение: %s\n", id, string(message.Value))
	}
}
