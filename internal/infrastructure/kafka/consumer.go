package kafka

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"sync"
	"time"
)

// Consumer обрабатывает сообщения из Kafka
type Consumer struct {
	consumer   *kafka.Consumer
	workers    int
	messageCh  chan *kafka.Message
	shutdownCh chan struct{}
	wg         sync.WaitGroup
}

// NewConsumer создает новый экземпляр Consumer
func NewConsumer(cfg *Config, workers int) (*Consumer, error) {
	kafkaConfig := cfg.CreateConsumerConfig()

	c, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:   c,
		workers:    workers,
		messageCh:  make(chan *kafka.Message, workers*3),
		shutdownCh: make(chan struct{}),
	}, nil
}

// Subscribe подписывается на топик и запускает обработчики
func (c *Consumer) Subscribe(topic string, handler func(*kafka.Message) error) error {
	if err := c.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}

	// Запуск воркеров
	for i := 0; i < c.workers; i++ {
		c.wg.Add(1)
		go c.worker(handler)
	}

	// Чтение сообщений
	go c.readMessages()

	return nil
}

func (c *Consumer) readMessages() {
	for {
		select {
		case <-c.shutdownCh:
			return
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("Consumer error: %v", err)
				continue
			}
			c.messageCh <- msg
		}
	}
}

func (c *Consumer) worker(handler func(*kafka.Message) error) {
	defer c.wg.Done()

	for msg := range c.messageCh {
		if err := handler(msg); err != nil {
			log.Printf("Failed to handle message: %v", err)
		}
	}
}

// Close останавливает потребителя
func (c *Consumer) Close() {
	close(c.shutdownCh)
	close(c.messageCh)
	c.wg.Wait()
	_ = c.consumer.Close()
}
