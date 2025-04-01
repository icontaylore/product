package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"product/internal/config"
	"product/internal/infrastructure/kafka"
	"product/internal/infrastructure/models"
	"sync"
	"syscall"
	"time"
)

func main() {
	// логгер
	logger := log.New(os.Stderr, "prefix: ", log.LstdFlags)
	// конфиг
	connStr, err := config.ParseConfig()
	if err != nil {
		logger.Fatal("parse err: ошибка в парсинге, невозможно запарсить конфиг")
	}

	// открываем бд
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logger.Fatal("open db: ошибка при открытии соединения с бд")
	}
	logger.Println("open db: успешное подключение к бд")
	// если нет таблицы, создаём её
	if err = db.AutoMigrate(&models.Product{}); err != nil {
		log.Printf("open db: ошибка с применением миграции: %v", err)
	}

	// kafka
	consumer := kafka.CreateKafkaConsumer()
	defer consumer.Close()

	// Подписываемся на топик
	partitionConsumer, err := consumer.ConsumePartition("dbz.public.products", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Ошибка при подписке на топик:", err)
	}
	defer partitionConsumer.Close()
	// Канал для сообщений
	messageChan := make(chan *sarama.ConsumerMessage)

	// Запуск воркеров
	var wg sync.WaitGroup
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go kafka.Worker(i, messageChan, &wg)
	}

	go func() {
		for message := range partitionConsumer.Messages() {
			fmt.Printf("Получено сообщение: %s\n", string(message.Value))
			messageChan <- message
		}
	}()

	// Таймер на 5 секунд
	timer := time.NewTimer(50 * time.Second)

	// Ожидаем сигнал завершения или таймер
	select {
	case <-timer.C:
		fmt.Println("Время вышло, завершение программы...")
		close(messageChan) // Закрываем канал для воркеров
	case <-waitForInterrupt():
		fmt.Println("Получен сигнал завершения, завершаем программу...")
		close(messageChan) // Закрываем канал для воркеров
	}

	// Ожидаем завершения всех воркеров
	wg.Wait()

	fmt.Println("Программа завершена")
}

// Функция для обработки сигнала завершения (Ctrl+C)
func waitForInterrupt() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}
