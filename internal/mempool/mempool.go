package mempool

import (
	"MEV_bot/config"
	"MEV_bot/internal/analysis"
	"MEV_bot/internal/strategy"
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"

	"MEV_bot/internal/client"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func StartMempoolMonitor(ethClient *client.EthClient) error {
	pendingTxs := make(chan string)
	sub, err := ethClient.Client.EthSubscribe(context.Background(), pendingTxs, "newPendingTransactions")
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	log.Println("Подписка на мемпул установлена")

	for txHash := range pendingTxs {
		//log.Printf("Обнаружена транзакция: %s", txHash)

		// Получаем детали транзакции
		tx, isValid := fetchTransactionDetails(ethClient, txHash)
		if !isValid {
			continue
		}

		// Фильтруем по интересующим контрактам (например, Uniswap)
		if !isTransactionTargetingDEX(tx) {
			continue
		}

		log.Printf("Расшифровываем данные транзакции")
		analysis.DecodeTransactionData("0x" + hex.EncodeToString(tx.Data()))

		// Сравниваем цены токенов между DEX
		amountIn := big.NewInt(1000000000000000000) // 1 токен (в Wei)
		path := []common.Address{
			common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), //USDT // Замените на реальные адреса токенов
			common.HexToAddress("0x514910771AF9Ca656af840dff83E8264EcF986CA"), //LINK
		}

		ethClientWrapper := ethclient.NewClient(ethClient.Client) // Обертка
		bestDex, bestPrice := strategy.FindBestPrice(ethClientWrapper, config.DexConfigs, amountIn, path)

		log.Printf("Лучший DEX: %s с ценой: %s", bestDex, bestPrice.String())

		// Исполняем моковую арбитражную транзакцию
		strategy.ExecuteMockArbitrage(bestDex, amountIn, bestPrice, path)
	}
	return nil
}

// Проверяет, нацелена ли транзакция на DEX (Uniswap или Sushiswap)
func isTransactionTargetingDEX(tx *types.Transaction) bool {
	targetDEX := []common.Address{
		common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"), // Uniswap
		common.HexToAddress("0xd9e1ce17f2641f24ae83637ab66a2cca9c378b9f"), // Sushiswap
	}

	for _, dexAddr := range targetDEX {
		log.Printf(*tx.To(), dexAddr)
		if tx.To() != nil && *tx.To() == dexAddr {
			return true
		}
	}
	return false
}

// Получает детали транзакции
func fetchTransactionDetails(client *client.EthClient, txHash string) (*types.Transaction, bool) {
	var tx *types.Transaction
	for i := 0; i < 5; i++ { // Пытаемся 5 раз
		err := client.Client.CallContext(context.Background(), &tx, "eth_getTransactionByHash", common.HexToHash(txHash))
		if err != nil {
			log.Printf("Ошибка вызова eth_getTransactionByHash для транзакции %s: %v", txHash, err)
			return nil, false
		}
		if tx != nil {
			return tx, true
		}
		log.Printf("Транзакция %s пока недоступна, повтор через 100 мс", txHash)
		time.Sleep(100 * time.Millisecond)
	}
	log.Printf("Транзакция %s не найдена после нескольких попыток", txHash)
	return nil, false
}

func getSenderAddress(tx *types.Transaction) (string, error) {
	// Выбираем тип подписанта (обычно EIP-155 используется в большинстве сетей Ethereum)
	signer := types.LatestSignerForChainID(tx.ChainId())

	// Извлекаем подпись и проверяем её
	sender, err := signer.Sender(tx)
	if err != nil {
		return "", err
	}

	return sender.Hex(), nil
}
