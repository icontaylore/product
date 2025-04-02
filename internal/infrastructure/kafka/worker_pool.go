package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
)

func (k *KafkaSettings) WorkerPoolStart() {
	n := k.Worker
	wg := sync.WaitGroup{}

	wg.Add(int(n))
	for i := 0; i < int(n); i++ {
		go func(i int) {
			k.worker(i, &wg)
		}(i)
	}
}

func (k *KafkaSettings) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for message := range k.MessageChan {
		doc := map[string]interface{}{
			"key":       string(message.Key),
			"value":     string(message.Value),
			"partition": message.Partition,
			"offset":    message.Offset,
			"timestamp": message.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
		if err := k.ElasticService.IndexDocument(context.Background(), doc); err != nil {
			log.Printf("kafka+elastic:workerpool warn")
		}
		fmt.Printf("worker-pool:worker %d обрабатывает сообщение\n", id)
	}
}
