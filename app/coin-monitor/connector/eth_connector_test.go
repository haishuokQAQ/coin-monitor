package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/constant"

	"github.com/ethereum/go-ethereum/common/compiler"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestCompileContract(t *testing.T) {
	contracts, err := compiler.CompileSolidity("", "/Users/konghaishuo/Documents/go/src/github.com/haishuokQAQ/coin-monitor/test_contract/0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3.sol")
	if err != nil {
		panic(err)
	}
	for _, contract := range contracts {
		fmt.Println(reflect.TypeOf(contract.Info.AbiDefinition))
	}
}

func TestGetPendingBlock(t *testing.T) {
	ctx := context.Background()
	rpcClient, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/3385d6a07d13405faaa658d91f44af25")
	if err != nil {
		panic(err)
	}
	err = InitEthHttpClient("https://mainnet.infura.io/v3/3385d6a07d13405faaa658d91f44af25")
	if err != nil {
		panic(err)
	}
	interChan := make(chan interface{}, 100)
	sub1, err := rpcClient.Subscribe(ctx, "eth", interChan, "newPendingTransactions")
	if err != nil {
		panic(err)
	}
	go func() {
		for txId := range interChan {
			if _, ok := txId.(string); !ok {
				continue
			}
			txIdStr := txId.(string)
			fmt.Println(txIdStr)
			tx, pending, err := ethClient.TransactionByHash(ctx, common.HexToHash(txIdStr))
			if err != nil {
				fmt.Println("Get tx error!", err)
				continue
			}
			if !pending {
				fmt.Println("Not pending!")
				continue
			}
			msg, err := tx.AsMessage(types.NewEIP155Signer(big.NewInt(1)))
			if err != nil {
				fmt.Println("As msg error!", err)
				continue
			}
			fmt.Println(tx.To().Hash(), msg.From().Hex())
		}
	}()
	for err2 := range sub1.Err() {
		panic(err2)
	}
}

func TestGetTokenTransfer(t *testing.T) {
	ctx := context.Background()
	err := InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	bscHttpConnector := GetBscClient()
	tx, _, err := bscHttpConnector.TransactionByHash(ctx, common.HexToHash("0xc62e61aacb5983d5c92e13bda363a5ced09a83d2203b0ca1c558513868d3508c"))
	if err != nil {
		panic(err)
	}
	receipt, err := bscHttpConnector.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		panic(err)
	}
	for _, log := range receipt.Logs {
		if log.Topics[0].Hex() == constant.TopicTransfer {
			fmt.Println("From:", log.Topics[1].Hex(), "|To:", common.BigToHash(log.Topics[2].Big()).Hex(), "|Token Address:", log.Address.Hex(), "|Value:", common.HexToHash(common.Bytes2Hex(log.Data)).Big().Uint64())
			//fmt.Println("From:", log.Topics[1].Hex(), "|To:", log.Topics[2].Hex(), "|Token Address:", log.Address.Hex(), "|Value:", common.HexToHash(common.Bytes2Hex(log.Data)).Big().Uint64())
		}
	}
}

func TestPrintInput(t *testing.T) {
	fmt.Println(common.HexToHash("000000000000000000000000000000000000000000000000002386f26fc10000").Big().Uint64())
}

func TestJson(t *testing.T) {
	//jsonStr := "{\"ext\":{\"a\":1,\"b\":\"c\"}}"
	//jsonStr := "{\"ext\":1}"
	jsonStr := "{\"ext\":\"2\"}"
	result := &struct {
		Ext interface{} `json:"ext"`
	}{}
	fmt.Println(json.Unmarshal([]byte(jsonStr), result))
	fmt.Println(result)
}

func TestGetTxDetail(t *testing.T) {
	ctx := context.Background()
	err := InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	bscHttpConnector := GetBscClient()
	tx, _, err := bscHttpConnector.TransactionByHash(ctx, common.HexToHash("0x422c301b27ab585c0276afdd6bb22d561209cbf7109920690184e1684ade19ef"))
	if err != nil {
		panic(err)
	}
	hexStr := common.Bytes2Hex(tx.Data())
	fmt.Println(hexStr)
	a := "0000000000000000000000000000000000000000000000000000000000000003"
	b := "0000000000000000000000000000000000000000000000000000000064bdaba9"
	fmt.Println(len([]rune(a)), len([]rune(b)))
	receipt, err := bscHttpConnector.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		panic(err)
	}
	StoreABI, err := LoadStoreABI()
	if err != nil {
		panic(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		panic(err)
	}
	for _, log := range receipt.Logs {
		if log.Address.Hex() == strings.ToLower("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3") {

		}
		event, err := contractAbi.Unpack("ItemSet", log.Data)
		if err != nil {
			panic(err)
		}
		fmt.Println(event)
	}
}

func LoadStoreABI() (string, error) {
	fileName := "/Users/konghaishuo/Documents/go/src/github.com/haishuokQAQ/coin-monitor/test_contract/0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3.sol"
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TestGetStrLen(t *testing.T) {
	fmt.Println(len([]rune("0x000000000000000000000000773355277126cbdcf8eb80702f6bc1a3cb843bbb")))
	fmt.Println(len([]rune("0x773355277126cbdcf8eb80702f6bc1a3cb843bbb")))
	fmt.Println(len([]rune("0x55d398326f99059ff775485246999027b3197955")))
	fmt.Println(fmt.Sprintf("0x%+v", strings.TrimPrefix(strings.TrimPrefix("0x000000000000000000000000773355277126cbdcf8eb80702f6bc1a3cb843bbb", "0x"), "0")))
	fmt.Println(fmt.Sprintf("0x%+v", string([]rune("0x000000000000000000000000773355277126cbdcf8eb80702f6bc1a3cb843bbb")[26:])))
}

func TestExecuteContract(t *testing.T) {
	err := InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	client := GetBscClient()
	client.CallContract(context.Background(), ethereum.CallMsg{}, nil)
}

func TestCallContract(t *testing.T) {
	err := InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	client := GetBscClient()
	parsed, err := abi.JSON(strings.NewReader(TestABI))
	if err != nil {
		panic(err)
	}

	address := common.HexToAddress("0xAA2f4ce8cF3687a8b4F98176BE48a746583d5627")

	// Value 参数编码
	valueInput, err := parsed.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	fmt.Println(reflect.TypeOf(valueInput))
	fmt.Println(valueInput)
	fmt.Println(common.Bytes2Hex(valueInput))

	method, err := parsed.MethodById(valueInput)
	if err != nil {
		panic(err)
	}
	fmt.Println(method.String())
	for s, m := range parsed.Methods {
		fmt.Println(s, m.String())
	}
	bc := bind.NewBoundContract(address, parsed, client, client, client)
	var result string
	param := &[]interface{}{
		&result,
	}
	err = bc.Call(&bind.CallOpts{
		Pending:     false,
		From:        common.HexToAddress("0x1b65bB376FF0a882cB3026D0c782021e84f75848"),
		BlockNumber: nil,
		Context:     context.Background(),
	}, param, "name")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

const TestABI = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"minTokensBeforeSwap","type":"uint256"}],"name":"MinTokensBeforeSwapUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"tokensSwapped","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"ethReceived","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tokensIntoLiqudity","type":"uint256"}],"name":"SwapAndLiquify","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bool","name":"enabled","type":"bool"}],"name":"SwapAndLiquifyEnabledUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[],"name":"_liquidityFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_maxTxAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_taxFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tAmount","type":"uint256"}],"name":"deliver","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"excludeFromFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"excludeFromReward","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"geUnlockTime","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"includeInFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"includeInReward","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"isExcludedFromFee","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"isExcludedFromReward","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"time","type":"uint256"}],"name":"lock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tAmount","type":"uint256"},{"internalType":"bool","name":"deductTransferFee","type":"bool"}],"name":"reflectionFromToken","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"liquidityFee","type":"uint256"}],"name":"setLiquidityFeePercent","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"maxTxPercent","type":"uint256"}],"name":"setMaxTxPercent","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bool","name":"_enabled","type":"bool"}],"name":"setSwapAndLiquifyEnabled","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"taxFee","type":"uint256"}],"name":"setTaxFeePercent","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"swapAndLiquifyEnabled","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"rAmount","type":"uint256"}],"name":"tokenFromReflection","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalFees","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"uniswapV2Pair","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"uniswapV2Router","outputs":[{"internalType":"contract IUniswapV2Router02","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unlock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`
