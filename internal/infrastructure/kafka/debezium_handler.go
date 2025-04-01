package kafka

import (
	"encoding/json"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
)

type DebeziumEvent struct {
	Before *Product `json:"before"`
	After  *Product `json:"after"`
	Op     string   `json:"op"`
	Source struct {
		Connector string `json:"connector"`
		Table     string `json:"table"`
	} `json:"source"`
}

type Product struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stockQuantity"`
}

func NewDebeziumHandler() func(*kafka.Message) error {
	return func(msg *kafka.Message) error {
		var event DebeziumEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		switch event.Op {
		case "c":
			log.Printf("INSERT: %+v", event.After)
		case "u":
			log.Printf("UPDATE: %+v", event.After)
		case "d":
			log.Printf("DELETE: %d", event.Before.Id)
		default:
			log.Printf("Unknown operation: %s", event.Op)
		}

		return nil
	}
}
