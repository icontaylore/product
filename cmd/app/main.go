package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"product/internal/config"
	"product/internal/infrastructure/models"
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
	// если нету таблицы, создаём её
	if err = db.AutoMigrate(&models.Product{}); err != nil {
		logger.Println("open db: ошибка с применением миграции")
	}
	fmt.Println("успешно")
