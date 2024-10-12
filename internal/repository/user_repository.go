package repository

import (
	"Golang-DeFiExplorer/internal/models"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
)

func SaveUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (username, email, password, created_at) 
              VALUES ($1, $2, $3, NOW()) RETURNING id`

	err := db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

func GetUserPrivateKeys(db *sql.DB, userID int64) ([]*ecdsa.PrivateKey, error) {
	var privateKeys [][]byte // Храним ключи в виде среза байтов
	query := `SELECT private_key FROM user_private_keys WHERE user_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key []byte
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании ключа: %w", err)
		}
		privateKeys = append(privateKeys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после завершения сканирования: %w", err)
	}

	var keys []*ecdsa.PrivateKey
	for _, keyBytes := range privateKeys {
		block, _ := pem.Decode(keyBytes)
		if block == nil || block.Type != "EC PRIVATE KEY" {
			return nil, fmt.Errorf("неверный формат приватного ключа")
		}

		privKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге ключа: %w", err)
		}
		keys = append(keys, privKey)
	}

	return keys, nil
}
