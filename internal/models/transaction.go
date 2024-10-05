package models

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
)

type Transaction struct {
	Id        []byte
	Inputs    []Input
	Outputs   []Output
	Signature []byte
}

type Input struct {
	TransactionId []byte
	OutIndex      int64
	PubKey        []byte
}

type Output struct {
	Value   int64
	Address []byte
}

func New(from, to []byte, amount int64) *Transaction {
	return &Transaction{
		Inputs:  []Input{{PubKey: from}},
		Outputs: []Output{{Value: amount, Address: to}},
	}
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) error {
	hash := tx.Hash()
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return err
	}
	tx.Signature = append(r.Bytes(), s.Bytes()...)
	return nil
}

func (tx *Transaction) Hash() []byte {
	data, _ := json.Marshal(tx)
	hash := sha256.Sum256(data)
	return hash[:]
}
