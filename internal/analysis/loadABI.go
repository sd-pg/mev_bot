package analysis

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"os"
	"strings"
)

var UniswapRouterParsedABI abi.ABI // Глобальная переменная для хранения распарсенного ABI

// LoadABI загружает ABI из файла и парсит его
func LoadABI(filepath string) abi.ABI {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Ошибка чтения ABI файла: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Ошибка парсинга ABI: %v", err)
	}

	log.Println("ABI успешно загружен и распарсен")
	return parsedABI
}

// InitializeABI инициализирует глобальную переменную с ABI
func InitializeABI() {
	UniswapRouterParsedABI = LoadABI("abi/uniswap_router_abi.json")
}
