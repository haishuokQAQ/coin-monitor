package mysql

import (
	"context"
	"fmt"
	"testing"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func initMysqlConnector() {
	err := connector.InitMysqlConnector(config.MysqlHost, config.MysqlPort, config.MysqlUserName, config.MysqlPasswd, config.MysqlDBName)
	if err != nil {
		panic(err)
	}
}

func TestListPlatform(t *testing.T) {
	initMysqlConnector()
	platform, err := ListPlatform(context.Background())
	if err != nil {
		panic(err)
	}
	for _, m := range platform {
		fmt.Println(*m)
	}
}

func TestGetCoinByPlatform(t *testing.T) {
	initMysqlConnector()
	platform, err := ListPlatform(context.Background())
	if err != nil {
		panic(err)
	}
	for _, m := range platform {
		fmt.Println(m.PlatformName)
		coins, err := GetCoinByPlatform(context.Background(), m.Id)
		if err != nil {
			panic(err)
		}
		for _, coin := range coins {
			fmt.Println(*coin)
		}
	}
}
