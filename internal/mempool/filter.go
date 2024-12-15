package mempool

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func IsTargetTransaction(tx *types.Transaction) bool {
	// Проверяем адрес получателя
	targetContract := common.HexToAddress("0xUniswapRouterAddress")
	if tx.To() == nil || *tx.To() != targetContract {
		return false // Не транзакция, связанная с Uniswap
	}

	// Проверяем данные вызова
	data := tx.Data()
	if len(data) < 4 {
		return false // Неверные данные
	}

	// Проверяем метод вызова (например, swapExactTokensForTokens)
	methodID := data[:4]
	targetMethodID := []byte{0x38, 0xed, 0x17, 0x39} // Метод ID swapExactTokensForTokens
	if !strings.EqualFold(string(methodID), string(targetMethodID)) {
		return false
	}

	// Пример: можно дополнительно фильтровать по объёму транзакции
	if tx.Value().Int64() < 1000000000000000000 { // < 1 ETH
		return false
	}

	return true
}
