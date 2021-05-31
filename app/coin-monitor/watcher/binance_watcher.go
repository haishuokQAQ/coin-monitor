package watcher

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/dal/influx"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"

	"github.com/haishuokQAQ/coin-monitor/app/config"

	"github.com/adshao/go-binance/v2"
)

type Config struct {
	Symbol string `json:"symbol"`
}

func NewCoinWatcher(conf *Config) *BinanceWatcher {
	return &BinanceWatcher{
		Symbol:        conf.Symbol,
		InfluxWriter:  influx.NewInfluxWriter(config.Bucket),
		watcherConfig: conf,
		stopChans:     map[string]chan struct{}{},
		stopChansLock: sync.Mutex{},
	}
}

type BinanceWatcher struct {
	Symbol               string               `json:"symbol"`
	InfluxWriter         *influx.InfluxWriter `json:"influx_writer"`
	watcherConfig        *Config              `json:"watcher_config"`
	stopChans            map[string]chan struct{}
	stopChansLock        sync.Mutex
	retryTimes           int
	nextRetryShouldAfter *time.Time
	errorChan            chan error
	doneChan             chan struct{}
}

func (watcher *BinanceWatcher) Cancel() {
	for _, stopChan := range watcher.stopChans {
		stopChan <- struct{}{}
	}
}

func (watcher *BinanceWatcher) StartWatch() (doneChan chan struct{}, errChan chan error, err error) {
	defer func() {
		if err != nil {
			watcher.Cancel()
		}
	}()
	watcher.doneChan = make(chan struct{})
	watcher.errorChan = make(chan error)
	// 开启K线监听
	if err := watcher.initKLineWatcher(); err != nil {
		return nil, nil, err
	}
	// 开启交易监听

	// 开启挂单监听
	if err := watcher.initDeptWatcher(); err != nil {
		return nil, nil, err
	}
	return watcher.doneChan, watcher.errorChan, nil
}

func (watcher *BinanceWatcher) initDeptWatcher() error {
	ws := NewBinanceOrderBookWatcher(watcher.Symbol)
	if err := ws.StartWatch(); err != nil {
		return err
	}
	watcher.stopChans["depth"] = ws.CancelChan
	return nil
}

func (watcher *BinanceWatcher) restartSingleKlineWatcher(interval string, recordFinal bool) error {
	key := fmt.Sprintf("%s%v", interval, recordFinal)
	watcher.stopChansLock.Lock()
	defer watcher.stopChansLock.Unlock()
	if _, ok := watcher.stopChans[key]; ok {
		watcher.stopChans[key] <- struct{}{}
		delete(watcher.stopChans, key)
	}
	err := watcher.initSingleKlineWatcher(interval, recordFinal)
	if err != nil {
		return err
	}
	return nil
}

func (watcher *BinanceWatcher) initSingleTradeWatcher() error {
	return nil
}

func (watcher *BinanceWatcher) initSingleKlineWatcher(interval string, recordFinal bool) error {
	_, stopChan, err := binance.WsKlineServe(watcher.Symbol, interval, watcher.solveEvent(recordFinal), func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		return err
	}
	watcher.stopChans[fmt.Sprintf("%s%v", interval, recordFinal)] = stopChan
	return nil
}

func (watcher *BinanceWatcher) initKLineWatcher() (err error) {
	defer func() {
		if err != nil {
			for _, stopChan := range watcher.stopChans {
				stopChan <- struct{}{}
			}
		}
	}()
	// 监听实时k线
	err = watcher.initSingleKlineWatcher("1m", false)
	if err != nil {
		return err
	}
	err = watcher.initSingleKlineWatcher("1m", true)
	if err != nil {
		return err
	}
	err = watcher.initSingleKlineWatcher("5m", true)
	if err != nil {
		return err
	}
	err = watcher.initSingleKlineWatcher("15m", true)
	if err != nil {
		return err
	}
	err = watcher.initSingleKlineWatcher("1h", true)
	if err != nil {
		return err
	}
	err = watcher.initSingleKlineWatcher("4h", true)
	if err != nil {
		return err
	}
	err = watcher.initSingleKlineWatcher("1d", true)
	if err != nil {
		return err
	}
	return nil
}

func (watcher *BinanceWatcher) solveEvent(recordFinalKline bool) binance.WsKlineHandler {
	return func(event *binance.WsKlineEvent) {
		if !event.Kline.IsFinal && recordFinalKline {
			return
		}
		var recordTime time.Time
		if recordFinalKline {
			recordTime = time.Unix(event.Kline.EndTime/1000, 0)
		} else {
			recordTime = time.Now()
			event.Kline.Interval = "1s"
		}
		watcher.InfluxWriter.WriteKlineToInfluxDB(context.Background(), event.Kline, recordTime)
	}
}

type BinanceOrderBookWatcher struct {
	Symbol     string
	BookLock   sync.Mutex
	Bid        BookItemList
	Ask        BookItemList
	CancelChan chan struct{}
}

type BookItem struct {
	Price        float64 `json:"price"`
	Quantity     float64 `json:"quantity"`
	LastUpdateId int64   `json:"last_update_id"`
}

type BookItemList []*BookItem

func (p BookItemList) Len() int {
	return len(p)
}

func (p BookItemList) Less(i, j int) bool { return p[i].Price < p[j].Price }

func (p BookItemList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func NewBinanceOrderBookWatcher(symbol string) *BinanceOrderBookWatcher {
	watcher := &BinanceOrderBookWatcher{
		Symbol:   symbol,
		BookLock: sync.Mutex{},
		Bid:      BookItemList{},
		Ask:      BookItemList{},
	}
	return watcher
}

func (ow *BinanceOrderBookWatcher) Cancel() {
	ow.CancelChan <- struct{}{}
}

func (ow *BinanceOrderBookWatcher) StartWatch() error {
	_, sc, err := binance.WsDepthServe100Ms(ow.Symbol, func(event *binance.WsDepthEvent) {
		ow.BookLock.Lock()
		defer ow.BookLock.Unlock()
		ow.Bid = UpdateNewBidListFromEvent(ow.Bid, event.Bids, event.FirstUpdateID, event.UpdateID)
		ow.Ask = UpdateNewAskListFromEvent(ow.Ask, event.Asks, event.FirstUpdateID, event.UpdateID)
	}, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		return err
	}
	ow.CancelChan = sc
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		client := connector.GetBinanceClient()
		for {
			time.Sleep(1 * time.Second)
			select {
			case <-ow.CancelChan:
				return
			default:
			}
			res, err := client.NewDepthService().Symbol(ow.Symbol).Limit(5000).Do(context.Background())
			if err != nil {
				fmt.Println("Depth watcher ", ow.Symbol, " get depth service error ", err)
				continue
			}
			ow.BookLock.Lock()

			ow.Bid = UpdateNewBidListFromSnapshot(ow.Bid, res.Bids, res.LastUpdateID)
			ow.Ask = UpdateNewAskListFromSnapshot(ow.Ask, res.Asks, res.LastUpdateID)
			ow.BookLock.Unlock()
		}
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		mongoCollection := connector.GetMongoConnection().Database("coin_monitor").Collection("depth")
		for {
			time.Sleep(1 * time.Second)
			select {
			case <-ow.CancelChan:
				return
			default:
			}
			currentAsk := ow.Ask
			currentBid := ow.Bid
			// 写入
			_, err := mongoCollection.InsertOne(context.Background(), map[string]interface{}{
				"symbol":    ow.Symbol,
				"platform":  "binance",
				"ask":       currentAsk,
				"bid":       currentBid,
				"timestamp": time.Now().Unix(),
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	return nil
}

func UpdateNewBidListFromSnapshot(existList BookItemList, respList []binance.Bid, lastUpdatedId int64) BookItemList {
	bidMap := map[float64]*BookItem{}

	for _, bid := range respList {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			continue
		}
		item := &BookItem{
			Price:        price,
			Quantity:     quantity,
			LastUpdateId: lastUpdatedId,
		}
		bidMap[item.Price] = item
	}
	for i, item := range existList {
		// 如果更新id小于等于lastUpdatedId则抛弃
		if item.LastUpdateId <= lastUpdatedId {
			continue
		}
		if _, ok := bidMap[item.Price]; !ok {
			bidMap[item.Price] = existList[i]
		}
	}
	newBidList := BookItemList{}
	for k := range bidMap {
		newBidList = append(newBidList, bidMap[k])
	}
	sort.Sort(newBidList)
	return newBidList
}

func UpdateNewAskListFromSnapshot(existList BookItemList, respList []binance.Ask, lastUpdatedId int64) BookItemList {
	askMap := map[float64]*BookItem{}

	for _, bid := range respList {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			continue
		}
		item := &BookItem{
			Price:        price,
			Quantity:     quantity,
			LastUpdateId: lastUpdatedId,
		}
		askMap[item.Price] = item
	}
	for i, item := range existList {
		// 如果更新id小于等于lastUpdatedId则抛弃
		if item.LastUpdateId <= lastUpdatedId {
			continue
		}
		if _, ok := askMap[item.Price]; !ok {
			askMap[item.Price] = existList[i]
		}
	}
	newBidList := BookItemList{}
	for k := range askMap {
		newBidList = append(newBidList, askMap[k])
	}
	sort.Sort(newBidList)
	return newBidList
}

func UpdateNewBidListFromEvent(existList BookItemList, respList []binance.Bid, firstUpdateId, updateId int64) BookItemList {
	bidMap := map[float64]*BookItem{}
	for index, item := range existList {
		bidMap[item.Price] = existList[index]
	}
	for _, bid := range respList {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			continue
		}
		item := &BookItem{
			Price:        price,
			Quantity:     quantity,
			LastUpdateId: updateId,
		}
		if _, ok := bidMap[item.Price]; !ok {
			bidMap[item.Price] = item
		} else {
			exist := bidMap[item.Price]
			if !(exist.LastUpdateId > firstUpdateId && exist.LastUpdateId < updateId) {
				continue
			}
			bidMap[item.Price] = item
		}
	}
	newBidList := BookItemList{}
	for k := range bidMap {
		newBidList = append(newBidList, bidMap[k])
	}
	sort.Sort(newBidList)
	return newBidList
}
func UpdateNewAskListFromEvent(existList BookItemList, respList []binance.Ask, firstUpdateId, updateId int64) BookItemList {
	askMap := map[float64]*BookItem{}
	for index, item := range existList {
		askMap[item.Price] = existList[index]
	}
	for _, bid := range respList {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			continue
		}
		item := &BookItem{
			Price:        price,
			Quantity:     quantity,
			LastUpdateId: updateId,
		}
		if _, ok := askMap[item.Price]; !ok {
			askMap[item.Price] = item
		} else {
			exist := askMap[item.Price]
			if !(exist.LastUpdateId > firstUpdateId && exist.LastUpdateId < updateId) {
				continue
			}
			askMap[item.Price] = item
		}
	}
	newAskList := BookItemList{}
	for k := range askMap {
		newAskList = append(newAskList, askMap[k])
	}
	sort.Sort(newAskList)
	return newAskList
}
