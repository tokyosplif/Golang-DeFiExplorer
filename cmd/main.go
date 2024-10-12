package main

import (
	"Golang-DeFiExplorer/cmd/server"
	"Golang-DeFiExplorer/internal/blockchain"
	"Golang-DeFiExplorer/internal/db"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dbConn := db.InitPostgres()
	defer dbConn.Close()

	blockchainURL := os.Getenv("ETHEREUM_RPC_URL")
	if blockchainURL == "" {
		log.Fatal("ETHEREUM_RPC_URL не задан в окружении")
	}

	indexer := blockchain.NewIndexer(dbConn)

	go func() {
		for {
			block, err := blockchain.GetLatestBlock(blockchainURL)
			if err != nil {
				log.Printf("Ошибка при получении блока: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}
			indexer.IndexBlock(block)
			time.Sleep(15 * time.Second)
		}
	}()

	r := server.NewRouter(indexer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
