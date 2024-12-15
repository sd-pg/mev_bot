package analysis

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

// Мапа обработчиков методов
var methodHandlers = map[string]func([]interface{}){
	"swapExactTokensForTokens":                           handleSwapExactTokensForTokens,
	"swapExactETHForTokens":                              handleSwapExactETHForTokens,
	"swapExactETHForTokensSupportingFeeOnTransferTokens": handleSwapExactETHForTokensSupportingFeeOnTransferTokens,
	"multicall":       handleMulticall,
	"addLiquidity":    handleAddLiquidity,
	"removeLiquidity": handleRemoveLiquidity,
	// Добавляйте другие методы
}

// DecodeTransactionData декодирует данные транзакции и вызывает соответствующий обработчик
func DecodeTransactionData(inputData string) {
	// Декодирование входящих данных транзакции
	data, err := hex.DecodeString(inputData[2:])
	if err != nil {
		log.Printf("Ошибка декодирования данных: %v", err)
		return
	}

	// Определяем метод вызова по ID
	method, err := UniswapRouterParsedABI.MethodById(data[:4])
	if err != nil {
		log.Printf("Неподдерживаемая транзакция с ID: %x", data[:4])
		return
	}

	// Распаковка параметров метода
	params, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		log.Printf("Ошибка распаковки параметров метода %s: %v", method.Name, err)
		return
	}

	log.Printf("Обнаружен вызов метода: %s", method.Name)

	// Ищем обработчик в мапе
	if handler, exists := methodHandlers[method.Name]; exists {
		handler(params)
	} else {
		log.Printf("Метод %s не поддерживается", method.Name)
	}
}

// Пример обработчиков
func handleSwapExactTokensForTokens(params []interface{}) {
	amountIn := params[0].(*big.Int)
	amountOutMin := params[1].(*big.Int)
	path := params[2].([]common.Address)
	to := params[3].(common.Address)
	deadline := params[4].(*big.Int)

	log.Printf("Обработан swapExactTokensForTokens:")
	log.Printf("AmountIn: %v, AmountOutMin: %v, Path: %v, To: %v, Deadline: %v", amountIn, amountOutMin, path, to, deadline)
}

func handleSwapExactETHForTokens(params []interface{}) {
	amountOutMin := params[0].(*big.Int)
	path := params[1].([]common.Address)
	to := params[2].(common.Address)
	deadline := params[3].(*big.Int)

	log.Printf("Обработан swapExactETHForTokens:")
	log.Printf("AmountOutMin: %v, Path: %v, To: %v, Deadline: %v", amountOutMin, path, to, deadline)
}

func handleMulticall(params []interface{}) {
	deadline := params[0]
	calls := params[1].([][]byte)

	log.Printf("Обработка multicall с deadline: %v", deadline)

	for i, call := range calls {
		subMethod, err := UniswapRouterParsedABI.MethodById(call[:4])
		if err != nil {
			log.Printf("Неизвестный метод в multicall: %v", err)
			continue
		}

		subParams, err := subMethod.Inputs.Unpack(call[4:])
		if err != nil {
			log.Printf("Ошибка распаковки параметров для %s: %v", subMethod.Name, err)
			continue
		}

		log.Printf("Вложенный вызов %d: Метод: %s, Параметры: %v", i+1, subMethod.Name, subParams)
	}
}

func handleAddLiquidity(params []interface{}) {
	tokenA := params[0].(common.Address)   // Адрес токена A
	tokenB := params[1].(common.Address)   // Адрес токена B
	amountADesired := params[2].(*big.Int) // Желаемое количество токенов A
	amountBDesired := params[3].(*big.Int) // Желаемое количество токенов B
	amountAMin := params[4].(*big.Int)     // Минимально допустимое количество токенов A
	amountBMin := params[5].(*big.Int)     // Минимально допустимое количество токенов B
	to := params[6].(common.Address)       // Адрес получателя токенов ликвидности
	deadline := params[7].(*big.Int)       // Дедлайн

	log.Printf("Обработан вызов addLiquidity:")
	log.Printf("Token A: %v, Token B: %v", tokenA, tokenB)
	log.Printf("Amount A Desired: %v, Amount B Desired: %v", amountADesired, amountBDesired)
	log.Printf("Amount A Min: %v, Amount B Min: %v", amountAMin, amountBMin)
	log.Printf("Recipient: %v, Deadline: %v", to, deadline)
}

func handleRemoveLiquidity(params []interface{}) {
	tokenA := params[0].(common.Address) // Адрес токена A
	tokenB := params[1].(common.Address) // Адрес токена B
	liquidity := params[2].(*big.Int)    // Количество токенов ликвидности для удаления
	amountAMin := params[3].(*big.Int)   // Минимально допустимое количество токенов A
	amountBMin := params[4].(*big.Int)   // Минимально допустимое количество токенов B
	to := params[5].(common.Address)     // Адрес получателя токенов
	deadline := params[6].(*big.Int)     // Дедлайн

	log.Printf("Обработан вызов removeLiquidity:")
	log.Printf("Token A: %v, Token B: %v", tokenA, tokenB)
	log.Printf("Liquidity: %v", liquidity)
	log.Printf("Amount A Min: %v, Amount B Min: %v", amountAMin, amountBMin)
	log.Printf("Recipient: %v, Deadline: %v", to, deadline)
}

func handleSwapExactETHForTokensSupportingFeeOnTransferTokens(params []interface{}) {
	amountOutMin := params[0].(*big.Int) // Минимальное количество токенов
	path := params[1].([]common.Address) // Путь обмена (массив адресов)
	to := params[2].(common.Address)     // Адрес получателя
	deadline := params[3].(*big.Int)     // Дедлайн выполнения

	log.Printf("Обработан swapExactETHForTokensSupportingFeeOnTransferTokens:")
	log.Printf("Минимальное количество токенов (amountOutMin): %v", amountOutMin)
	log.Printf("Путь обмена (path): %v", path)
	log.Printf("Получатель (to): %v", to)
	log.Printf("Дедлайн (deadline): %v", deadline)
}
