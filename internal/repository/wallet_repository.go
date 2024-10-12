package repository

import (
	"Golang-DeFiExplorer/internal/models"
	"database/sql"
	"fmt"
)

func GetWalletByAddress(db *sql.DB, address string) (*models.Wallet, error) {
	var wallet models.Wallet

	query := `SELECT address, balance, created_at FROM wallets WHERE address = $1`
	err := db.QueryRow(query, address).Scan(&wallet.Address, &wallet.Balance, &wallet.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("кошелек не найден по адресу: %s", address)
		}
		return nil, fmt.Errorf("ошибка при получении кошелька: %w", err)
	}

	return &wallet, nil
}
