package client

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

// EthClient определяет структуру для взаимодействия с Ethereum
type EthClient struct {
	Client *rpc.Client
}

// NewEthClient создает новый клиент для подключения к Ethereum
func NewEthClient(rpcURL string) *EthClient {
	client, err := rpc.DialContext(context.Background(), rpcURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к Ethereum RPC: %v", err)
	}
	return &EthClient{Client: client}
}

// Close закрывает соединение клиента
func (e *EthClient) Close() {
	e.Client.Close()
}
