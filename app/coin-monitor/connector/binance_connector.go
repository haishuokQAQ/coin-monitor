package connector

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
)

var client *binance.Client
var futuresClient *futures.Client
var deliveryClient *delivery.Client

func InitBinaceClient(apiKey, secretKey string) {
	proxyUrl, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	netTransport := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		Proxy: http.ProxyURL(proxyUrl),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
	}
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	client = &binance.Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   "https://api1.binance.com",
		UserAgent: "Binance/golang",
		HTTPClient: &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		},
		Debug:  true,
		Logger: log.New(devNull, "Binance-golang ", log.LstdFlags),
	}
	futuresClient = binance.NewFuturesClient(apiKey, secretKey)   // USDT-M Futures
	deliveryClient = binance.NewDeliveryClient(apiKey, secretKey) // Coin-M Futures
}

func GetBinanceClient() *binance.Client {
	return client
}

func GetBinanceFuturesClient() *futures.Client {
	return futuresClient
}

func GetBinanceDeliveryClient() *delivery.Client {
	return deliveryClient
}
