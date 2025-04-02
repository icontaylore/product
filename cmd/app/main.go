package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"product/internal/config"
	"product/internal/infrastructure/elasticsearch"
	"product/internal/infrastructure/kafka"
	"product/internal/infrastructure/models"
	"syscall"
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

	// elastic init
	elastic, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatal("main:addres lagging elastic")
	}
	// elastic index
	elasticService := elasticsearch.NewService(elastic, "product-index")
	// Проверить, существует ли индекс (опционально)
	res, err := elastic.Indices.Exists([]string{"product-index"})
	if err != nil {
		log.Fatal("Ошибка проверки индекса:", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println("Индекс существует!")
	} else if res.StatusCode == 404 {
		fmt.Println("Индекс не найден.")
	}

	// kafka init
	configKafka := kafka.GoConfigure("dbz.public.products", "localhost:9094", 5)
	configKafka.CreateKafkaConsumer()
	configKafka.SetElasticService(elasticService)

	if err = configKafka.SubscribeTopic(); err != nil {
		log.Fatal("main:ошибка подписки на топик")
	}
	defer configKafka.Consumer.Close()

	// Запуск воркеров
	configKafka.WorkerPoolStart()

	// Ожидаем сигнал завершения или таймер
	select {
	case <-waitForInterrupt():
		fmt.Println("Получен сигнал завершения, завершаем программу...")
		if configKafka.PartitionCons != nil {
			configKafka.PartitionCons.Close()
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
