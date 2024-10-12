package blockchain

import (
	"Golang-DeFiExplorer/internal/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Monitor struct {
	blockchainURL string
	indexer       *Indexer
}

func NewMonitor(blockchainURL string, indexer *Indexer) *Monitor {
	return &Monitor{
		blockchainURL: blockchainURL,
		indexer:       indexer,
	}
}

func (m *Monitor) Start() {
	for {
		m.checkForNewBlocks()
		time.Sleep(10 * time.Second)
	}
}

func (m *Monitor) checkForNewBlocks() {
	log.Println("Checking for new blocks...")

	lastBlock, err := GetLatestBlock(m.blockchainURL) // Здесь вызываем функцию
	if err != nil {
		log.Printf("Ошибка при получении последнего блока: %v", err)
		return
	}

	m.indexer.IndexBlock(lastBlock)
}

func GetLatestBlock(blockchainURL string) (models.Block, error) {
	var blockResult struct {
		Result models.Block `json:"result"`
	}

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{"latest", true},
		"id":      1,
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(blockchainURL, "application/json", jsonBody(reqBody))
	if err != nil {
		return models.Block{}, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&blockResult); err != nil {
		return models.Block{}, err
	}

	return blockResult.Result, nil
}

func jsonBody(data interface{}) io.Reader {
	body, _ := json.Marshal(data)
	return bytes.NewReader(body)
}
