package mongodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/mongo"
)

func Test_InsertHolderInfo(t *testing.T) {
	connector.InitMongoConnection()
	err := InsertHolderInfo(context.Background(), "current", &mongo.ERC20TokenHolderInf{
		Address:     "0X0000123",
		TotalSupply: 0,
		Decimal:     0,
		Timestamp:   0,
		Holders:     0,
		Transfers:   0,
		HolderInf: map[string]uint64{
			"abc": 123,
		},
	})
	if err != nil {
		panic(err)
	}
}

func Test_GetCurrentHolderInfoForAddress(t *testing.T) {
	connector.InitMongoConnection()
	result, err := GetCurrentHolderInfoForAddress(context.Background(), "0X0000123")
	if err != nil {
		panic(err)
	}
	fmt.Println(*result)
}
