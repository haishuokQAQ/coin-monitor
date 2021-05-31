package adapter

import (
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func BinanceKlineToMap(kline binance.WsKline) map[string]interface{} {
	openDouble, err := strconv.ParseFloat(kline.Open, 64)
	if err != nil {
		openDouble = 0
	}
	closeDouble, err := strconv.ParseFloat(kline.Close, 64)
	if err != nil {
		closeDouble = 0
	}
	highDouble, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		highDouble = 0
	}
	lowDouble, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		lowDouble = 0
	}
	volumeDouble, err := strconv.ParseFloat(kline.Volume, 64)
	if err != nil {
		volumeDouble = 0
	}
	quoteVolumeDouble, err := strconv.ParseFloat(kline.QuoteVolume, 64)
	if err != nil {
		quoteVolumeDouble = 0
	}
	activeBuyVolumeDouble, err := strconv.ParseFloat(kline.ActiveBuyVolume, 64)
	if err != nil {
		activeBuyVolumeDouble = 0
	}
	activeBuyQuoteVolumeDouble, err := strconv.ParseFloat(kline.ActiveBuyQuoteVolume, 64)
	if err != nil {
		activeBuyQuoteVolumeDouble = 0
	}

	return map[string]interface{}{
		"start_time":              kline.StartTime,
		"end_time":                kline.EndTime,
		"symbol":                  kline.Symbol,
		"interval":                kline.Interval,
		"first_trade_id":          kline.FirstTradeID,
		"last_trade_id":           kline.LastTradeID,
		"open":                    openDouble,
		"close":                   closeDouble,
		"high":                    highDouble,
		"low":                     lowDouble,
		"volume":                  volumeDouble,
		"trade_num":               kline.TradeNum,
		"quote_volume":            quoteVolumeDouble,
		"active_buy_volume":       activeBuyVolumeDouble,
		"active_buy_quote_volume": activeBuyQuoteVolumeDouble,
	}
}
