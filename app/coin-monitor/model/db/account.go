package db

import (
	"time"
)

type Account struct {
	AccountAddress string     `json:"account_address"`
	LastActiveAt   *time.Time `json:"last_active_at"`
}
