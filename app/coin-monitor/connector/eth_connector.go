package connector

import "github.com/ethereum/go-ethereum/ethclient"

var ethClient *ethclient.Client

func InitEthHttpClient(net string) error {
	ethClientInstance, err := ethclient.Dial(net)
	if err != nil {
		return err
	}
	ethClient = ethClientInstance
	return nil
}

func GetEthClient() *ethclient.Client {
	return ethClient
}
