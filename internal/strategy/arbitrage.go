package strategy

import (
	"MEV_bot/config"
	"MEV_bot/internal/analysis"
	"MEV_bot/internal/dex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func CalculateArbitrage(reserve1, reserve2 *big.Int, amountIn *big.Int) *big.Int {
	// Расчёт цены покупки и продажи
	price1 := new(big.Int).Div(reserve1, reserve2)
	price2 := new(big.Int).Div(reserve2, reserve1)

	profit := new(big.Int).Sub(price2, price1)
	return profit
}

// FindBestPrice ищет лучший DEX для обмена токенов
func FindBestPrice(client *ethclient.Client, dexConfigs []config.DEXConfig, amountIn *big.Int, path []common.Address) (string, *big.Int) {
	bestPrice := big.NewInt(0)
	bestDex := ""

	for _, dexConfig := range dexConfigs {
		// Используем функцию GetQuote из пакета dex
		price, err := dex.GetQuote(client, dexConfig.RouterAddr, amountIn, path, analysis.UniswapRouterParsedABI)
		if err != nil {
			log.Printf("Ошибка получения котировки для DEX %s: %v", dexConfig.Name, err)
			continue
		}

		log.Printf("Цена на %s: %s", dexConfig.Name, price.String())

		if price.Cmp(bestPrice) > 0 {
			bestPrice = price
			bestDex = dexConfig.Name
		}
	}

	return bestDex, bestPrice
}

// ExecuteMockArbitrage выполняет моковую арбитражную транзакцию
func ExecuteMockArbitrage(bestDex string, amountIn *big.Int, amountOut *big.Int, path []common.Address) {
	log.Printf("Выполняем моковую арбитражную транзакцию:")
	log.Printf("DEX: %s, Amount In: %s, Amount Out: %s, Path: %v", bestDex, amountIn.String(), amountOut.String(), path)
}
