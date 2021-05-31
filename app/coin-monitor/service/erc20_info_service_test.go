package service

import (
	"fmt"
	"testing"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/proxy/bscscan"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
)

var caller *ERC20TokenCodeHolder
var holder *ERC20TokenCodeHolder

func initErc20TokenCaller() {
	err := connector.InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	bscscan.InitEtherScanClient()
	holder, err = NewErc20TokenInfoService(connector.GetBscClient(), "0xf09b7B6bA6dAb7CccC3AE477a174b164c39f4C66", "")
	if err != nil {
		panic(err)
	}
	caller = holder
}

func TestERC20TokenCaller_GetDecimals(t *testing.T) {
	initErc20TokenCaller()
	fmt.Println(caller.GetDecimals())
}

func TestERC20TokenCaller_GetName(t *testing.T) {
	initErc20TokenCaller()
	fmt.Println(caller.GetName())
}

func TestERC20TokenCaller_GetSymbol(t *testing.T) {
	initErc20TokenCaller()
	fmt.Println(caller.GetSymbol())
}

func TestERC20TokenCaller_GetTotalSupply(t *testing.T) {
	initErc20TokenCaller()
	fmt.Println(caller.GetTotalSupply())
	fmt.Println(caller.GetTotalSupplyBigInt())
}
