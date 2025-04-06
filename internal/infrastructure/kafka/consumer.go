package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"sync"
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

func (k *KafkaSettings) WorkerPool(w int) {
	wg := sync.WaitGroup{}
	for i := 0; i < w; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			k.worker(i, k.IndexName)
		}()
	}
	wg.Wait()
}

func (k *KafkaSettings) worker(id int, indexName string) {
	for msg := range k.PartitionCons.Messages() {
		var debezMessage DebeziumMessage
		if err := json.Unmarshal(msg.Value, &debezMessage); err != nil {
			log.Printf("worker %d:удаление из бд: %v\nRaw data: %s", id, err, msg.Value)
			continue
		}

		log.Printf("worker %d работает\n", id)
		switch debezMessage.Payload.Op {
		case "c":
			product := parseProduct(debezMessage.Payload.After)
			if err := k.sendToElastic(product, "index", indexName); err != nil {
				log.Println("consumer kafka:create errr")
			}
		case "u":
			product := parseProduct(debezMessage.Payload.After)
			if err := k.sendToElastic(product, "update", indexName); err != nil {
				log.Println("consumer kafka:update err")
			}
		case "d":
			if debezMessage.Payload.Before != nil {
				k.deleteFromElastic(debezMessage.Payload.Before.ID)
			}
		}
	}
}

func (k *KafkaSettings) sendToElastic(doc map[string]interface{}, op, indexName string) error {
	ctx := context.Background()
	docID := fmt.Sprintf("%v", doc["id"])
	var req esapi.Request
	var err error

	switch op {
	case "index":
		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(doc); err != nil {
			return err
		}
		req = esapi.IndexRequest{
			Index:      indexName,
			DocumentID: docID,
			Body:       &buf,
			Refresh:    "true",
		}
	case "update": // Частичное обновление
		var buf bytes.Buffer
		updateBody := map[string]interface{}{
			"doc": doc,
		}
		if err = json.NewEncoder(&buf).Encode(updateBody); err != nil {
			return err
		}
		req = esapi.UpdateRequest{
			Index:      indexName,
			DocumentID: docID,
			Body:       &buf,
			Refresh:    "true",
		}
	default:
		return fmt.Errorf("worker:unknown operation: %s", op)
	}

	res, err := req.Do(ctx, k.ESClient)
	if err != nil {
		return fmt.Errorf("worker:request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("worker:error response: %s", res.String())
	}

	log.Printf("worker-pool:удачная операция %s ID %s", op, docID)
	return nil
}

func (k *KafkaSettings) deleteFromElastic(id uint) error {
	idInt := strconv.Itoa(int(id))

	var req esapi.Request
	var err error

	req = esapi.DeleteRequest{
		Index:      k.IndexName,
		DocumentID: idInt,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), k.ESClient)
	if err != nil {
		return fmt.Errorf("worker delete:request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("worker delete:error response: %s", res.String())
	}

	log.Printf("worker delete:performed delete operation for document ID %s", k.IndexName)
	return nil
}
