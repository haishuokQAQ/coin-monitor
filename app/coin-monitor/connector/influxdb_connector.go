package connector

import (
    `github.com/haishuokQAQ/coin-monitor/app/config`
    influxdb2 `github.com/influxdata/influxdb-client-go/v2`
)

var influxClient influxdb2.Client

func InitInfluxClient(){
    influxClient = influxdb2.NewClient(config.Endpoint, config.InfluxToken)
}

func GetInfluxClient() influxdb2.Client{
    return influxClient
}