package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/db"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/dal/mysql"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/constant"
)

func ScanForTokenTransfer(startCount, endCount uint64) {
	bscHttpConnector := connector.GetBscClient()
	ctx := context.Background()
	targetCount := uint64(0)
	if endCount == 0 {
		blockCount, err := bscHttpConnector.BlockNumber(ctx)
		if err != nil {
			panic(err)
		}
		targetCount = blockCount
	} else {
		targetCount = endCount
	}

	for i := uint64(startCount); i < targetCount; i++ {
		fmt.Println("In block ", i)
		currentBlock, err := bscHttpConnector.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			fmt.Println(fmt.Sprintf("Read block %+v error.Err :%+v", i, err))
			continue
		}
		fmt.Println("Transaction count ", currentBlock.Transactions().Len())
		for _, transaction := range currentBlock.Transactions() {
			if transaction.To() != nil {
				receipt, err := bscHttpConnector.TransactionReceipt(ctx, transaction.Hash())
				if err != nil {
					fmt.Println(fmt.Sprintf(""))
					continue
				}
				_ = SaveErc20TransferDetail(ctx, receipt, currentBlock.Time(), transaction)
			}
		}
		if endCount == 0 {
			for {
				newBlockCount, err := bscHttpConnector.BlockNumber(ctx)
				if err != nil {
					fmt.Println(fmt.Sprintf("Fail to update block count.Error %+v", err))
					continue
				}
				if newBlockCount > targetCount {
					targetCount = newBlockCount
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func SaveErc20TransferDetail(ctx context.Context, receipt *types.Receipt, blockTime uint64, transaction *types.Transaction) error {
	for _, log := range receipt.Logs {
		if log.Topics[0].Hex() == constant.TopicTransfer && len(log.Topics) >= 3 {
			history := &db.TokenTransferHistory{
				ExecuteTime:  blockTime,
				TokenAddress: log.Address.Hex(),
				FromAddr:     util.HexTrim(log.Topics[1].Hex()),
				TargetAddr:   util.HexTrim(log.Topics[2].Hex()),
				TXAddress:    transaction.Hash().Hex(),
				Value:        common.HexToHash(common.Bytes2Hex(log.Data)).Big().Uint64(),
			}
			currentHistory, err := mysql.GetTokenTransferHistoryByTokenAndAddress(ctx, transaction.Hash().Hex(), log.Address.Hex())
			if err != nil {
				fmt.Println(fmt.Sprintf("get transfer history by tx id %+v and token addr %+v err.%+v", transaction.Hash().Hex(), log.Address.Hex(), err))
				continue
			}
			if currentHistory != nil {
				if *currentHistory == *history {
					continue
				}
			}
			err = mysql.InsertTokenTransferHistory(ctx, history)
			if err != nil {
				fmt.Println(fmt.Sprintf("insert transfer history %+v err.%+v", *history, err))
				continue
			}
		}
	}
	return nil
}
