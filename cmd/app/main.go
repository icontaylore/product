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

	// kafka setup
	brokers := []string{"localhost:9094"}
	topic := "dbz.public.products"
	// Consumer
	kf := kafka.KafkaGetConfig(brokers, topic)
	if err = kf.NewConsumer(); err != nil {
		log.Fatal("main:трабл с созданием консьюмера")
	}
	defer kf.Consumer.Close()
	// Subscribe
	if err = kf.SubscribeTopic(); err != nil {
		log.Fatal("main:subs make err")
	}


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
