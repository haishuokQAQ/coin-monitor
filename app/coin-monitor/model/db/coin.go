package db

func (Coin) TableName() string {
	return "coin"
}

type Coin struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
