package server

import (
	"Golang-DeFiExplorer/internal/blockchain"
	"Golang-DeFiExplorer/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(indexer *blockchain.Indexer) http.Handler {
	r := chi.NewRouter()

	r.Post("/user", handlers.CreateUser)
	r.Post("/transaction", handlers.CreateTransaction)
	r.Get("/block/{id}", handlers.GetBlock(indexer))
	r.Get("/wallet/{id}", handlers.GetWallet)

	return r
}
