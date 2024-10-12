package models

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Transaction struct {
	Id           int64
	Time         time.Time
	BlockId      int64
	Hash         string
	Size         int
	Weight       int
	IsCoinbase   bool
	Fee          int64
	FeeUSD       float64
	Details      interface{}
	Inputs       []Input
	Outputs      []Output
	Signatures   [][]byte
	RequiredSigs int
	From         string
	To           string
	Value        string
	UserId       int64
	CreatedAt    time.Time
}

type Input struct {
	TransactionId []byte
	OutIndex      int64
	PubKey        []byte
}

type Output struct {
	Value   string
	Address []byte
}

func New(from, to []byte, amount string) *Transaction {
	return &Transaction{
		Inputs:  []Input{{PubKey: from}},
		Outputs: []Output{{Value: amount, Address: to}},
		From:    string(from),
		To:      string(to),
		Value:   amount,
	}
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) error {
	hash := tx.CalculateHash()
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signatures = append(tx.Signatures, signature)
	return nil
}

func (tx *Transaction) CalculateHash() []byte {
	data, _ := json.Marshal(tx)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (tx *Transaction) Validate() error {
	if tx.From == "" || tx.To == "" || tx.Value == "" {
		return fmt.Errorf("неверные данные транзакции")
	}
	return nil
}

func isValidAddress(addr string) bool {
	return len(addr) > 0
}
