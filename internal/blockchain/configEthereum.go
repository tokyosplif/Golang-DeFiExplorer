package blockchain

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EthereumConfig struct {
	RPCURL string 
}

func LoadConfig() (*EthereumConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Ошибка при загрузке .env файла: %v", err)
		return nil, err
	}

	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		log.Println("Предупреждение: переменная окружения ETHEREUM_RPC_URL не установлена")
		return nil, nil
	}

	return &EthereumConfig{
		RPCURL: rpcURL,
	}, nil
}
