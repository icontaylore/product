package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"product/internal/config"
	"product/internal/ggrpc"
	"product/internal/infrastructure/elasticsearch"
	"product/internal/infrastructure/kafka"
	"product/internal/infrastructure/models"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Lоgger
	logger := log.New(os.Stderr, "prefix: ", log.LstdFlags)

	// Kонфиг запуска
	connStr, err := config.ParseConfig()
	if err != nil {
		logger.Fatal("parse err: ошибка в парсинге, невозможно запарсить конфиг")
	}

	// Open db
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logger.Fatal("open db: ошибка при открытии соединения с бд")
	}
	logger.Println("open db: успешное подключение к бд")
	// если нет таблицы, создаём её
	if err = db.AutoMigrate(&models.Product{}); err != nil {
		log.Printf("open db: ошибка с применением миграции: %v", err)
	}

	// Elastic
	addresElastic := []string{"http://localhost:9200"}
	// Elastic Client
	clientElastic := elasticsearch.NewClientElastic(addresElastic)
	if err != nil {
		log.Fatal("main:elastic conn err")
	}
	// Elastic Create Index
	indexName := "product_index"
	clientElastic.CreateIndex(indexName)

	// Kafka setup
	brokers := []string{"localhost:9094"}
	topic := "dbz.public.products"
	kf := kafka.KafkaGetConfig(brokers, topic)
	// Consumer
	if err = kf.NewConsumer(); err != nil {
		log.Fatal("main:трабл с созданием консьюмера")
	}
	// Грузим внешние данные
	kf.ESClient = clientElastic.Client
	kf.IndexName = indexName
	defer kf.Consumer.Close()

	// Subscribe
	if err = kf.SubscribeTopic(); err != nil {
		log.Fatal("main:subs make err")
	}
	// Read kafka
	worker := 5
	kf.WorkerPool(worker)

	// gRPC
	srv := ggrpc.NewServer("50051")
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()

	// Ожидаем сигнал завершения или таймер
	select {
	case <-waitForInterrupt():
		fmt.Println("Получен сигнал завершения, завершаем программу...")
		if kf.PartitionCons != nil {
			kf.PartitionCons.Close()
		}
	}

	fmt.Println("Программа завершена")
}

// Функция для обработки сигнала завершения (Ctrl+C)
func waitForInterrupt() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}
