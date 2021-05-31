package main

import (
	"fmt"

	watcher2 "github.com/haishuokQAQ/coin-monitor/app/coin-monitor/watcher"

	"golang.org/x/sync/errgroup"

	"github.com/fvbock/endless"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func initConnectors() {
	connector.InitBinaceClient(config.ApiKey, config.SecretKey)
	connector.InitInfluxClient()
	connector.InitMysql()
	connector.InitMongoConnection()
}

func initWatcher() {
	watcher := watcher2.NewCoinWatcher(&watcher2.Config{Symbol: "DOGEUSDT"})
	err, _, _ := watcher.StartWatch()
	if err != nil {
		panic(err)
	}
}

var g errgroup.Group

func main() {
	initConnectors()
	endlessServer := endless.NewServer(fmt.Sprintf(":%d", 26512), nil)
	initWatcher()
	g.Go(func() error {
		return endlessServer.ListenAndServe()
	})
	err := g.Wait()
	if err != nil {
		panic(err)
	}
}
