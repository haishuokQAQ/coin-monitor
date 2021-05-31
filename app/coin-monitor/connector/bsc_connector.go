package connector

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

var bscClient *ethclient.Client

func InitBscHttpClient(net string) error {
	ethClientInstance, err := ethclient.Dial(net)
	if err != nil {
		return err
	}
	bscClient = ethClientInstance
	return nil
}

func GetBscClient() *ethclient.Client {
	return bscClient
}
