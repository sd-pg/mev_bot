package main

import (
	"MEV_bot/config"
	"MEV_bot/internal/analysis"
	"MEV_bot/internal/client"
	"MEV_bot/internal/mempool"
	"log"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()
	log.Println("Загружена конфигурация")

	analysis.InitializeABI()

	// Инициализация Ethereum клиента
	ethClient := client.NewEthClient(cfg.RPC_URL)
	defer ethClient.Close()

	// Запуск мониторинга мемпула
	go func() {
		err := mempool.StartMempoolMonitor(ethClient)
		if err != nil {
			log.Fatalf("Ошибка во время работы мониторинга: %v", err)
		}
	}()

	// Удерживаем приложение в рабочем состоянии
	select {}
}
