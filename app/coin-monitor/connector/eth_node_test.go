package connector

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth"
)

func TestEthTxPool(t *testing.T) {
	txPool := makeTxPool()
	time.Sleep(2 * time.Second)
	txs, err := txPool.Pending()
	if err != nil {
		panic(err)
	}
	for address, transactions := range txs {
		fmt.Println("===")
		fmt.Println(address.Hex())
		for _, transaction := range transactions {
			fmt.Println("--")
			fmt.Println(transaction.Hash())
		}
	}
}

func TestStartMiner(t *testing.T) {
	ethIns := makeEthereum()
	var srv interface{}
	for i, api := range ethIns.APIs() {
		if api.Namespace == "miner" {
			srv = ethIns.APIs()[i].Service
			break
		}
	}
	minerApi := srv.(*eth.PrivateMinerAPI)
	thread := 1
	if err := minerApi.Start(&thread); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Hour)
}
