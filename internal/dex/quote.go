package dex

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetQuote получает цену токена на DEX
func GetQuote(client *ethclient.Client, routerAddr string, amountIn *big.Int, path []common.Address, parsedABI abi.ABI) (*big.Int, error) {
	data, err := parsedABI.Pack("getAmountsOut", amountIn, path)
	if err != nil {
		log.Printf("Ошибка упаковки данных для getAmountsOut: %v", err)
		return nil, err
	}

	routerAddress := common.HexToAddress(routerAddr) // Сохраняем результат в переменную
	msg := ethereum.CallMsg{
		To:   &routerAddress, // Передаём указатель на переменную
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Printf("Ошибка вызова контракта: %v", err)
		return nil, err
	}

	var amounts []*big.Int
	err = parsedABI.UnpackIntoInterface(&amounts, "getAmountsOut", result)
	if err != nil {
		log.Printf("Ошибка распаковки результата getAmountsOut: %v", err)
		return nil, err
	}

	return amounts[len(amounts)-1], nil
}
