package db

import (
	"time"
)

type TxHistory struct {
	ChainId        string     `json:"chain_id"`
	TransactionId  string     `json:"transaction_id"`
	IP             string     `json:"ip"`
	UserAgent      string     `json:"user_agent"`
	Characteristic string     `json:"characteristic"`
	CreatedAt      *time.Time `json:"created_at"`
}
