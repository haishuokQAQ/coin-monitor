package service

import (
	"testing"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/proxy/bscscan"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func TestScanForBscTokens(t *testing.T) {
	err := connector.InitMysqlConnector(config.MysqlHost, config.MysqlPort, config.MysqlUserName, config.MysqlPasswd, config.MysqlDBName)
	if err != nil {
		panic(err)
	}
	bscscan.InitEtherScanClient()
	ScanForBscTokens(7500000, 0)
}
