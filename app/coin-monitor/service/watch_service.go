package service

import (
	"errors"
	"sync"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/watcher"
)

var binanceWatcherMap map[string]*watcher.BinanceWatcher
var binanceWatcherLock sync.Mutex

func StartWatcher(symbol, platform string) error {
	binanceWatcherLock.Lock()
	defer binanceWatcherLock.Unlock()
	if _, ok := binanceWatcherMap[symbol]; ok {
		return errors.New("已存在一个watcher")
	}
	return nil
}
