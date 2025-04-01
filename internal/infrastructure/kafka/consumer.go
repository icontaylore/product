package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"product/internal/infrastructure/models"
	"syscall"
)

type KafkaMessage struct {
	Payload struct {
		Op     string          `json:"op"`
		Before *models.Product `json:"before"`
		After  *models.Product `json:"after"`
	} `json:"payload"`
}

type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"` // mb null
	Price         string  `json:"price"`
	StockQuantity int64   `json:"stock_quantity"`
}

func StartConsumer(topic string) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9093",
		"group.id":          "product-consumer",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(conf)
	if err != nil {
		log.Fatal("kafka:not create cons..")
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{topic}, nil)

	log.Println("kafka:consumer запущен")

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			fmt.Printf("kafka:signal %v операционной системы off", sig)
			return
		default:
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				handleMessage(msg.Value)
			} else {
				log.Printf("kafka:error reading message: %v\n", err)
			}
		}
	}
}
func handleMessage(msg []byte) {
	fmt.Printf("kafka:принято - %s\n", string(msg))
}
