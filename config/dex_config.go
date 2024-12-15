package config

// DEXConfig представляет конфигурацию для работы с конкретным DEX
type DEXConfig struct {
	Name       string
	RouterAddr string
	ABI        string // Путь к ABI
}

// DexConfigs содержит список поддерживаемых DEX
var DexConfigs = []DEXConfig{
	{
		Name:       "UniswapV2",
		RouterAddr: "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
		ABI:        "abi/uniswap_router_abi.json", // Путь к ABI-файлу
	},
	{
		Name:       "Sushiswap",
		RouterAddr: "0xd9e1ce17f2641f24ae83637ab66a2cca9c378b9f",
		ABI:        "abi/uniswap_router_abi.json", // Можно использовать тот же ABI
	},
}
