package mongo

type ERC20TokenHolderInf struct {
	Address     string            `json:"address"`
	TotalSupply uint64            `json:"total_supply"`
	Decimal     int               `json:"decimal"`
	Timestamp   uint64            `json:"timestamp"`
	Holders     uint64            `json:"holders"`
	Transfers   uint64            `json:"transfers"`
	HolderInf   map[string]uint64 `json:"holder_inf"`
}

func (ERC20TokenHolderInf) BaseCollectionName() string {
	return "token_holder_inf"
}
