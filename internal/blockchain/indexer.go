package blockchain

import (
	"Golang-DeFiExplorer/internal/models" 
	"database/sql"                        
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Indexer struct {
	database *sql.DB
}

func NewIndexer(database *sql.DB) *Indexer {
	return &Indexer{database: database}
}

func (i *Indexer) FetchLatestBlock() {
	url := os.Getenv("ETHEREUM_RPC_URL") 
	if url == "" {
		log.Fatal("ETHEREUM_RPC_URL не задан в переменных окружения")
	}

	body := strings.NewReader(`{
        "jsonrpc": "2.0",
        "method": "eth_blockNumber",
        "params": [],
        "id": 1
    }`)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Ответ на eth_blockNumber: %s", string(bodyBytes))

	var blockNumberResult struct {
		Result string `json:"result"`
	}

	if err := json.Unmarshal(bodyBytes, &blockNumberResult); err != nil {
		log.Fatalf("Ошибка при декодировании номера блока: %v", err)
	}

	log.Printf("Полученный номер блока (в шестнадцатеричном формате): %s", blockNumberResult.Result)

	height, err := strconv.ParseInt(blockNumberResult.Result, 0, 64)
	if err != nil {
		log.Printf("Ошибка преобразования высоты блока: %v", err)
		return
	}
	log.Printf("Преобразованная высота блока: %d", height)

	blockNumberHex := strconv.FormatInt(height, 16)
	log.Printf("Запрос данных о блоке с номером: 0x%s", blockNumberHex)

	body = strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": ["0x` + blockNumberHex + `", true],
		"id": 1
	}`)

	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ = io.ReadAll(resp.Body)
	log.Printf("Ответ на eth_getBlockByNumber: %s", string(bodyBytes))

	var blockResult struct {
		Result struct {
			Hash         string `json:"hash"`
			ParentHash   string `json:"parentHash"`
			Transactions []struct {
				Hash string `json:"hash"`
			} `json:"transactions"`
		} `json:"result"`
	}

	if err := json.Unmarshal(bodyBytes, &blockResult); err != nil {
		log.Fatalf("Ошибка при декодировании данных блока: %v", err)
	}

	log.Printf("Полученные данные блока: %+v", blockResult)

	block := models.Block{
		Hash:         blockResult.Result.Hash,
		PreviousHash: blockResult.Result.ParentHash,
		BlockNumber:  height,
	}

	log.Printf("Данные блока перед индексированием: %+v", block)

	i.IndexBlock(block)
}

func (i *Indexer) IndexBlock(block models.Block) {
	if block.Hash == "" {
		log.Println("Пропущена индексация блока из-за отсутствия необходимых данных")
		return
	}

	var existingID int64
	err := i.database.QueryRow("SELECT id FROM blocks WHERE hash = $1", block.Hash).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Ошибка при проверке существования блока: %v", err)
		return
	}

	if existingID > 0 {
		log.Printf("Блок с хешем %s уже существует с ID: %d", block.Hash, existingID)
		return
	}

       query := `
    INSERT INTO blocks (hash, previous_hash, block_number)
    VALUES ($1, $2, $3)
    RETURNING id;`

	var newID int64
	err = i.database.QueryRow(query, block.Hash, block.PreviousHash, block.BlockNumber).Scan(&newID)
	if err != nil {
		log.Printf("Ошибка при индексации блока: %v", err)
		return
	}

	block.Id = newID
	log.Printf("Блок успешно проиндексирован с ID: %d", block.Id)
	log.Printf("Обновленные данные блока после индексации: %+v", block)
}

func (i *Indexer) GetBlockByID(blockID int64) (*models.Block, error) {
	var block models.Block

	query := `SELECT id, hash, previous_hash FROM blocks WHERE id = $1`

	err := i.database.QueryRow(query, blockID).Scan(&block.Id, &block.Hash, &block.PreviousHash)
	if err != nil {
		log.Printf("Ошибка при извлечении блока по ID: %v", err)
		return nil, err
	}

	log.Printf("Извлеченный блок: %+v", block)
	return &block, nil
}
