package db

import (
	"time"
)

type AddressBook struct {
	Id             uint64     `json:"id"`
	AccountAddress string     `json:"account_address"`
	TargetAddress  string     `json:"target_address"`
	ExchangeId     uint64     `json:"exchange_id"`
	Symbol         string     `json:"symbol"`
	ChainId        string     `json:"chain_id"`
	Version        int        `json:"version"`
	CreatedAt      *time.Time `json:"created_at"`
}
