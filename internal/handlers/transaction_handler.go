package handlers

import (
	"Golang-DeFiExplorer/internal/db"
	"Golang-DeFiExplorer/internal/models"
	"Golang-DeFiExplorer/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
		return
	}

	if err := tx.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbInstance := db.GetDBInstance()

	balanceStr, err := repository.GetBalanceForAddress(dbInstance, tx.From)
	if err != nil {
		http.Error(w, "Ошибка при получении баланса", http.StatusInternalServerError)
		return
	}

	balance, err := strconv.ParseInt(strconv.FormatInt(balanceStr, 10), 10, 64)
	if err != nil {
		http.Error(w, "Ошибка при преобразовании баланса", http.StatusInternalServerError)
		return
	}

	txValue, err := strconv.ParseInt(tx.Value, 10, 64)
	if err != nil {
		http.Error(w, "Ошибка при преобразовании значения транзакции", http.StatusBadRequest)
		return
	}

	if balance < txValue {
		http.Error(w, "Недостаточно средств", http.StatusBadRequest)
		return
	}

	privateKeys, err := repository.GetUserPrivateKeys(dbInstance, tx.UserId)
	if err != nil {
		http.Error(w, "Ошибка при получении приватных ключей", http.StatusInternalServerError)
		return
	}

	for _, privKey := range privateKeys {
		if err := tx.Sign(privKey); err != nil {
			http.Error(w, "Ошибка при подписывании транзакции", http.StatusInternalServerError)
			return
		}
	}

	if err := repository.SaveTransaction(dbInstance, &tx); err != nil {
		http.Error(w, "Ошибка при сохранении транзакции", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}
