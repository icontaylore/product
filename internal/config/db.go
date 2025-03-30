package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB_USER       string
	DB_NAME       string
	HOST          string
	DB_PAS        int
	POSTGRES_PORT int
	DB_PORT       int
}

func ParseConfig() (string, error) {
	if err := godotenv.Load("../../configs/.env"); err != nil {
		log.Printf("err: ошибка загрузки кофига")
		return "", err
	}
	out := fmt.Sprintf("host=%d user=%s password=%s dbname=%s port=%d sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PAS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	return out, nil
}
