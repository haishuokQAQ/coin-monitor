package service

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/dal/mysql"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/db"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/proxy/bscscan"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/util"
	"github.com/haishuokQAQ/coin-monitor/app/config"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
)

func ScanForBscTokens(startCount, endCount uint64) {

	ethHttpConnector := connector.GetBscClient()
	ctx := context.Background()
	targetCount := uint64(0)
	if endCount == 0 {
		blockCount, err := ethHttpConnector.BlockNumber(ctx)
		if err != nil {
			panic(err)
		}
		targetCount = blockCount
	} else {
		targetCount = endCount
	}

	for i := uint64(startCount); i < targetCount; i++ {
		fmt.Println("In block ", i)
		currentBlock, err := ethHttpConnector.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			fmt.Println(fmt.Sprintf("Read block %+v error.Err :%+v", i, err))
			continue
		}
		fmt.Println("Transaction count ", currentBlock.Transactions().Len())
		for _, transaction := range currentBlock.Transactions() {
			if transaction.To() == nil {
				// 检查是否为合约
				receipt, err := ethHttpConnector.TransactionReceipt(context.Background(), transaction.Hash())
				if err != nil {
					fmt.Println(fmt.Sprintf("Read receipt of tx %s error.Err :%+v", transaction.Hash(), err))
					continue
				}
				if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
					err := SaveSmartContract(ctx, receipt, currentBlock.Time())
					if err != nil {
						continue
					}
				}
			}
		}
		if endCount == 0 {
			for {
				newBlockCount, err := ethHttpConnector.BlockNumber(ctx)
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
	}
}

func SaveSmartContract(ctx context.Context, receipt *types.Receipt, blockTime uint64) error {
	fmt.Println(receipt.ContractAddress.Hex())
	// 根据地址获取合约信息
	sourceCodeResp, err := bscscan.GetContractSourceCodeByAddress(ctx, receipt.ContractAddress.Hex())
	if err != nil {
		fmt.Println(fmt.Sprintf("Get contract of tx %s error.Err :%+v", receipt.ContractAddress.Hex(), err))
		return err
	}
	if len(sourceCodeResp.Result) > 0 {
		// 写入元数据
		contractInfo := &db.ContractInfo{
			Address:        receipt.ContractAddress.Hex(),
			Name:           sourceCodeResp.Result[0].ContractName,
			Symbol:         "",
			CompileVersion: sourceCodeResp.Result[0].CompilerVersion,
			TotalSupply:    0,
			CreatedAt:      blockTime,
		}
		meta, err := mysql.GetContractMetaById(ctx, contractInfo.Address)
		if err != nil {
			fmt.Println(fmt.Sprintf("get contract meta of  %+v err.%+v", contractInfo.Address, err))
			return err
		}
		if meta == nil {
			err = mysql.InsertContractMeta(ctx, contractInfo)
			if err != nil {
				fmt.Println(fmt.Sprintf("insert contract %+v err.%+v", *contractInfo, err))
				return err
			}
		} else {
			if meta.CompileVersion == "" && contractInfo.CompileVersion != "" {
				err = mysql.UpdateContractMeta(ctx, contractInfo)
				if err != nil {
					fmt.Println(fmt.Sprintf("update contract %+v err.%+v", *contractInfo, err))
					return err
				}
			}
		}
		// 写入代码文件
		dirPath := fmt.Sprintf("%s/contracts/%s", config.ContractBasicPath, contractInfo.Address)
		if !util.CheckFileExist(dirPath) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				fmt.Println(fmt.Sprintf("mkdir %+v err.%+v", dirPath, err))
				return err
			}
		}
		for _, detail := range sourceCodeResp.Result {
			fileName := contractInfo.Name
			if fileName == "" {
				fileName = contractInfo.Address
			}
			filePath := fmt.Sprintf("%s/contracts/%s/%s.sol", config.ContractBasicPath, contractInfo.Address, fileName)
			if util.CheckFileExist(filePath) {
				continue
			}
			err := ioutil.WriteFile(filePath, []byte(detail.SourceCode), fs.ModeAppend)
			if err != nil {
				fmt.Println(fmt.Sprintf("write file %+v err.%+v", filePath, err))
				continue
			}
		}
	}
	return nil
}
