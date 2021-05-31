package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/proxy/bscscan"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/dal/mysql"
)

func ScanForTokenDetail() {
	// 拉取所有的token信息
	pageNum := 1
	pageSize := 10
	ctx := context.Background()

	currentTotalCount := 0
	for {
		records, totalCount, err := mysql.SearchContractInfoForPage(ctx, pageNum, pageSize)
		if err != nil {
			panic(err)
		}
		for i, record := range records {
			currentTotalCount += 1
			if record.Symbol == "" {
				addr := common.HexToAddress(record.Address)
				ctx := context.Background()
				currentRecord := records[i]

				reuslt, err := bscscan.GetContractSourceCodeByAddress(ctx, addr.Hex())
				if err != nil {
					continue
				}
				for _, detail := range reuslt.Result {

					if !strings.HasPrefix(detail.ABI, "[") {
						continue
					}
					// 进行验证
					holder, err := NewErc20TokenInfoService(connector.GetBscClient(), record.Address, detail.ABI)
					if err != nil {
						fmt.Println("Fail to init caller.Address ", record.Address, "Error ", err)
						continue
					}

					caller := holder
					symbol, err := caller.GetSymbol()
					if err != nil {
						if record.Name != "" {
							fmt.Println("Fail to get symbol for ", record.Name, " address is ", record.Address, " error:", err)
						}
						continue
					}
					name, err := caller.GetName()
					if err != nil {
						continue
					}
					decimal, err := caller.GetDecimals()
					if err != nil {
						continue
					}
					totalSupply, err := caller.GetTotalSupply()
					if err != nil {
						totalSupply, err = caller.GetTotalSupplyBigInt()
						if err != nil {
							continue
						}
					}
					currentRecord.ContractName = currentRecord.Name
					currentRecord.Name = name
					currentRecord.Symbol = symbol
					currentRecord.Decimal = decimal
					currentRecord.TotalSupply = totalSupply
					err = mysql.UpdateContractMeta(ctx, currentRecord)
					if err != nil {
						fmt.Println(fmt.Sprintf("Fail to update contract info.Error %+v", err))
						continue
					}
					break
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
		if int64(currentTotalCount) >= totalCount {
			break
		}
		pageNum++
	}
}
