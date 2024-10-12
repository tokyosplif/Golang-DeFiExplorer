package repository

import (
	"Golang-DeFiExplorer/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
)

func GetBalanceForAddress(db *sql.DB, address string) (int64, error) {
	var balance int64
	query := `SELECT balance FROM wallets WHERE address = $1`
	err := db.QueryRow(query, address).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("адрес не найден: %s", address)
		}
		return 0, fmt.Errorf("ошибка при получении баланса: %w", err)
	}
	return balance, nil
}

func SaveTransaction(db *sql.DB, tx *models.Transaction) error {
	query := `INSERT INTO transactions (from_address, to_address, value, signatures, created_at) VALUES ($1, $2, $3, $4, $5)`
	signatures, err := json.Marshal(tx.Signatures)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации подписей: %w", err)
	}

	_, err = db.Exec(query, tx.From, tx.To, tx.Value, signatures, tx.CreatedAt)
	return err
}
