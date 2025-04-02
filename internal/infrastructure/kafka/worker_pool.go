package kafka

import (
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
		fmt.Printf("worker-pool:worker %d обрабатывает сообщение\n", id)
		log.Printf(
			"Получено сообщение: ключ=%s, значение=%s, партиция=%d, оффсет=%d",
			string(message.Key),
			string(message.Value),
			message.Partition,
			message.Offset,
		)
	}
}
