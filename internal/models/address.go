package models

import "time"

type Address struct {
	Id        int64
	Address   string
	WalletId  int64
	CreatedAt time.Time
}
