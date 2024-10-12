package models

import "time"

type Wallet struct {
	Id        int64
	Address   string
	UserId    int64
	Balance   int64
	CreatedAt time.Time
}
