package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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

func (k *KafkaSettings) WorkerPool(w int, indexName string) {
	wg := sync.WaitGroup{}
	for i:=0;i<w;i++{
		go func()  {
			wg.Add(1)
			defer wg.Done()
			k.worker(i,indexName)
		}()
	}
	wg.Wait()
}


func (k *KafkaSettings) worker(id int,indexName string) {
	for msg := range k.PartitionCons.Messages() {
		var kafkaMsg KafkaMessage
		if err := json.Unmarshal(msg.Value,&kafkaMsg);err != nil {
            log.Printf("worker %d:oшибка парсинга JSON: %v\nRaw data: %s", id, err, msg.Value)
            continue
        }

		esp := fastFormat(&kafkaMsg)
		esp.Metadata.KafkaTopic = msg.Topic
		esp.Metadata.KafkaOffset = msg.Offset
		esp.Metadata.KafkaPartition = msg.Partition


		// go to elastic
		jsonDoc,_ := json.Marshal(esp)
		req := esapi.IndexRequest{
			Index: indexName,
			Body: strings.NewReader(string(jsonDoc)),
		}

		res, err := req.Do(context.Background(), k.ESClient)
        if err != nil {
            log.Printf("Worker %d:oшибка отправки в ES: %v", id, err)
            continue
        }
        defer res.Body.Close()

        if res.IsError() {
            log.Printf("Worker %d:oшибка Elasticsearch: %s", id, res.String())
        } else {
            log.Printf("Worker %d:товар %d отправлен в ES", id, kafkaMsg.ID)
        }
	}
}

func fastFormat(kafkaMsg *KafkaMessage) ESProduct {
	esDoc := ESProduct{
		ID:            kafkaMsg.ID,
		Name:          kafkaMsg.Name,
		Description:   kafkaMsg.Description,
		Price:         kafkaMsg.Price,
		StockQuantity: kafkaMsg.StockQuantity,
		Metadata: struct {
			KafkaTopic     string `json:"kafka_topic"`
			KafkaPartition int32  `json:"kafka_partition"`
			KafkaOffset    int64  `json:"kafka_offset"`
		}{
		},
	}
	return esDoc
}