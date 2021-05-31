package db

type Symbol struct {
	Name           string `json:"name"`
	CoinName       string `json:"coin_name"`
	SourceCoinName string `json:"source_coin_name"`
}
