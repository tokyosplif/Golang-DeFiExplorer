package models

type UserPrivateKey struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	PrivateKey string `json:"private_key"`
}
