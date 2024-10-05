package models

type Block struct {
	Id           int64
	Height       int64
	Hash         string
	PreviousHash string
	Timestamp    int64
	Transactions []Transaction
}
