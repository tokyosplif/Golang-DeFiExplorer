package models

import "time"

type Block struct {
	Id           int64
	Hash         string
	PreviousHash string
	ParentHash   string
	BlockNumber  int64
	CreatedAt    time.Time
	Transactions []Transaction
}
