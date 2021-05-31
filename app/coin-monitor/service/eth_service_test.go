package service

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/proxy/bscscan"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
)

func TestSyncEthBlock(t *testing.T) {
	startCount := 0
	ctx := context.Background()
	bscscan.InitEtherScanClient()
	err := connector.InitEthHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	ethHttpConnector := connector.GetEthClient()
	blockCount, err := ethHttpConnector.BlockNumber(ctx)
	if err != nil {
		panic(err)
	}
	startCount = int(blockCount - 1000)
	for i := uint64(startCount); i < blockCount; i++ {
		fmt.Println("In block ", i)
		currentBlock, err := ethHttpConnector.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			panic(err)
		}
		fmt.Println("Transaction count ", currentBlock.Transactions().Len())
		for _, transaction := range currentBlock.Transactions() {
			if transaction.To() == nil {
				// 检查是否为合约
				receipt, err := ethHttpConnector.TransactionReceipt(context.Background(), transaction.Hash())
				if err != nil {
					panic(err)
				}
				if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
					fmt.Println(receipt.ContractAddress.Hex())
					// 根据地址获取合约信息
					sourceCodeResp, err := bscscan.GetContractSourceCodeByAddress(ctx, receipt.ContractAddress.Hex())
					if err != nil {
						panic(err)
					}
					fmt.Println(len(sourceCodeResp.Result))
				}
			}
		}
		blockCount, err = ethHttpConnector.BlockNumber(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func TestGetTx(t *testing.T) {
	err := connector.InitEthHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	ethHttpConnector := connector.GetEthClient()
	tx, _, err := ethHttpConnector.TransactionByHash(context.Background(), common.HexToHash("0x0f5a88ab1540fca2b3300f72994977a8d74232b588077101e2ef82fd325422b9")) //0x37b3215a229ab60bd574c14a1f408066fe37ec5b1e7424c3ec7efc0d32094e7d
	if err != nil {
		panic(err)
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(big.NewInt(56)))
	if err != nil {
		panic(err)
	}
	fmt.Println(msg.To())
	fmt.Println(msg.From().Hex())
	receipt, err := ethHttpConnector.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(receipt.ContractAddress.Hex())
	fmt.Println(receipt.Type)
}
