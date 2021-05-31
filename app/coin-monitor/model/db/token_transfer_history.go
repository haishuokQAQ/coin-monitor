package db

type TokenTransferHistory struct {
	ExecuteTime  uint64 `json:"execute_time"`
	TokenAddress string `json:"token_address"`
	FromAddr     string `json:"from_addr"`
	TargetAddr   string `json:"target_addr"`
	Value        uint64 `json:"value"`
	TXAddress    string `json:"tx_address"`
}
