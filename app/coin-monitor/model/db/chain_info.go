package db

type ChainInfo struct {
	ChainId        string `json:"chain_id"`
	ChainName      string `json:"chain_name"`
	EndPoint       string `json:"end_point"`
	ApplicationKey string `json:"application_key"`
	ScanUrl        string `json:"scan_url"`
}
