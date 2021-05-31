package chain

type Erc20Token struct {
	Address     string            `json:"address"`
	Name        string            `json:"name"`
	Symbol      string            `json:"symbol"`
	TotalSupply uint64            `json:"total_supply"`
	Demicals    uint64            `json:"demicals"`
	Balance     map[string]uint64 `json:"balance"`
}
