package util

import (
    `fmt`
    `github.com/adshao/go-binance/v2`
    `testing`
)

func TestEncodeToMap(t *testing.T) {
    resultMap := map[string]interface{}{}
    err := EncodeToMap(&binance.Kline{
        OpenTime: 123,
        Open: "1",
        High: "2",
    }, resultMap)
    if err != nil {
        panic(err)
    }
    fmt.Println(resultMap)
}
