package kafka

import (
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
)

type KafkaSettings struct {
	ESClient *elasticsearch.Client
	Consumer      sarama.Consumer
	Addres        []string
	Topic         string
	Message       <-chan *sarama.ConsumerMessage
	PartitionCons sarama.PartitionConsumer
}

type KafkaMessage struct {
    ID            uint   `json:"id"`
    Name          string `json:"name"`
    Description   string `json:"description"`
    Price         uint   `json:"price"`
    StockQuantity uint   `json:"stock_quantity"` // Исправлено на snake_case для JSON
}

type ESProduct struct {
    ID            uint   `json:"id"`
    Name          string `json:"name"`
    Description   string `json:"description"`
    Price         uint   `json:"price"`
    StockQuantity uint   `json:"stock_quantity"`
    Metadata      struct {
        KafkaTopic     string `json:"kafka_topic"`
        KafkaPartition int32  `json:"kafka_partition"`
        KafkaOffset    int64  `json:"kafka_offset"`
    } `json:"metadata"`
}

func KafkaGetConfig(arr []string, s string) *KafkaSettings {
	return &KafkaSettings{
		Addres: arr,
		Topic:  s,
	}
}
