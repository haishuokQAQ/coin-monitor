package main

import (
	"context"
	"os"
	"strconv"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/proxy/bscscan"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/service"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func init() {
	connector.InitMysql()
	err := connector.InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
}

func main() {
	config.ContractBasicPath = "/home/data"
	startCount, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	threadCount := 4
	if len(os.Args) > 2 {
		threadCountStr := os.Args[2]
		count, err := strconv.ParseInt(threadCountStr, 10, 64)
		if err == nil && count > 0 && count < 32 {
			threadCount = int(count)
		}
	}
	err = connector.InitMysqlConnector(config.MysqlHost, config.MysqlPort, config.MysqlUserName, config.MysqlPasswd, config.MysqlDBName)
	if err != nil {
		panic(err)
	}
	bscscan.InitEtherScanClient()
	// 获取高度
	err = connector.InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	height, err := connector.GetBscClient().BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	iterateCount := height - startCount
	batchSize := iterateCount / uint64(threadCount)
	currentStart := startCount
	for i := 0; i < threadCount-1; i++ {
		if i == threadCount-1 {
			continue
		}
		go service.ScanForBscTokens(currentStart, currentStart+batchSize)
		currentStart += batchSize
	}
	service.ScanForBscTokens(currentStart, 0)
}
