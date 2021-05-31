package db

type ContractInfo struct {
	Address        string `json:"address"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	CompileVersion string `json:"compile_version"`
	TotalSupply    uint64 `json:"total_supply"`
	CreatedAt      uint64 `json:"created_at"`
	ContractName   string `json:"contract_name"`
	Decimal        uint8  `json:"decimal"`
}
