package handlers

import (
	"Golang-DeFiExplorer/internal/db"
	"Golang-DeFiExplorer/internal/repository"
	"encoding/json"
	"net/http"
)

func GetWallet(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Адрес кошелька обязателен", http.StatusBadRequest)
		return
	}

	dbInstance := db.GetDBInstance()

	wallet, err := repository.GetWalletByAddress(dbInstance, address)
	if err != nil {
		http.Error(w, "Ошибка при получении кошелька: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}
