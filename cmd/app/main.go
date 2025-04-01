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

	// kafka
	topic := "dbz.public.products"
	logger.Println("kafka topic:потребитель для топика:", topic)
	go kafka.StartConsumer(topic)

	// Создаем канал для получения сигналов о завершении работы программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) // Перехватываем сигналы завершения

	// Ожидаем сигнала для завершения работы
	<-stop
	log.Println("Завершение работы программы.")
}
