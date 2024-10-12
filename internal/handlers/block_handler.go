package handlers

import (
	"Golang-DeFiExplorer/internal/blockchain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetBlock(indexer *blockchain.Indexer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockIDStr := chi.URLParam(r, "id")
		if blockIDStr == "" {
			http.Error(w, "Block ID is required", http.StatusBadRequest)
			return
		}

		blockID, err := strconv.ParseInt(blockIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid block ID", http.StatusBadRequest)
			return
		}

		block, err := indexer.GetBlockByID(blockID)
		if err != nil {
			http.Error(w, "Block not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(block); err != nil {
			http.Error(w, "Failed to encode block", http.StatusInternalServerError)
		}
	}
}
