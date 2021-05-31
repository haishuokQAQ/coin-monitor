package db

func (PlatformCoin) TableName() string {
	return "platform_coin"
}

type PlatformCoin struct {
	PlatformId uint64 `json:"platform_id"`
	CoinName   string `json:"coin_name"`
}
