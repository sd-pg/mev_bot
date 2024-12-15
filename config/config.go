package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурации
type Config struct {
	RPC_URL string
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() *Config {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Читаем URL для RPC
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Fatalf("Не указан RPC_URL в .env")
	}

	return &Config{
		RPC_URL: rpcURL,
	}
}
