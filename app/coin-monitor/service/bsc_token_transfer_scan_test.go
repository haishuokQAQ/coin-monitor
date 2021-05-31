package service

import (
	"testing"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func TestScanForTokenTransfer(t *testing.T) {
	err := connector.InitMysqlConnector(config.MysqlHost, config.MysqlPort, config.MysqlUserName, config.MysqlPasswd, config.MysqlDBName)
	if err != nil {
		panic(err)
	}
	err = connector.InitBscHttpClient("https://bsc-dataseed1.ninicoin.io")
	if err != nil {
		panic(err)
	}
	ScanForTokenTransfer(7500000, 0)
}
