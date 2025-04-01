package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"product/internal/config"
	"product/internal/infrastructure/kafka"
	"product/internal/infrastructure/models"
	"syscall"
)

func main() {
	// log
	logger := log.Logger{}
	// conf
	connStr, err := config.ParseConfig()
	if err != nil {
		logger.Fatal("parse err: ошибка в парсинге, невозможно запарсить конфиг")
	}
	// db init open
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logger.Fatal("open db: ошибка при открытии соединения с бд")
	}
	logger.Println("open db: успешное подключение к бд")
	// если нет таблицы, создаём её
	if err = db.AutoMigrate(&models.Product{}); err != nil {
		logger.Println("open db: ошибка с применением миграции")
	}
	fmt.Println("успешно")

	// kafka cfg
	configKafka := kafka.NewConfig("localhost:9093", "product-service-group")
	// 3 workers
	consumer, err := kafka.NewConsumer(configKafka, 3)
	if err != nil {
		log.Fatal("main:ошибка в загрузке консьюмера")
	}
	defer consumer.Close()
	// debezium handle
	handler := kafka.NewDebeziumHandler()
	// subscribe topic
	if err := consumer.Subscribe("dbz.public.products", handler); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// создаем канал для получения сигналов о завершении работы программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) // Перехватываем сигналы завершения

	// ожидаем сигнала для завершения работы
	<-stop
	log.Println("Завершение работы программы.")
}
