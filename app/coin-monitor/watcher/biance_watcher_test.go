package watcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func TestTradeServe(t *testing.T) {
	connector.InitBinaceClient(config.ApiKey, config.SecretKey)
	dc, sc, err := binance.WsTradeServe("DOGEUSDT", func(event *binance.WsTradeEvent) {
		fmt.Println(*event)
	}, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(15 * time.Second)
		sc <- struct{}{}
	}()
	<-dc
}

func TestDepthServe(t *testing.T) {
	connector.InitBinaceClient(config.ApiKey, config.SecretKey)
	dc, sc, err := binance.WsDepthServe("DOGEUSDT", func(event *binance.WsDepthEvent) {
		fmt.Println(*event)
	}, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(200 * time.Millisecond)
		sc <- struct{}{}
	}()
	<-dc
}

func TestBinanceOrderBookWatcher_StartWatch(t *testing.T) {
	connector.InitBinaceClient(config.ApiKey, config.SecretKey)
	w := NewBinanceOrderBookWatcher("DOGEUSDT")
	err := w.StartWatch()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100000; i++ {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("ASK")
		for _, item := range w.Ask {
			fmt.Println(item)
		}
		fmt.Println("BID")
		for _, item := range w.Bid {
			fmt.Println(item)
		}

	}

}
