package influx

import (
	"context"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/adapter"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/config"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func NewInfluxWriter(bucket string) *InfluxWriter {
	return &InfluxWriter{
		Bucket: bucket,
		lock:   sync.Mutex{},
	}
}

type InfluxWriter struct {
	Bucket    string `json:"bucket"`
	apiWriter api.WriteAPI
	lock      sync.Mutex
}

func (writer *InfluxWriter) WriteKlineToInfluxDB(ctx context.Context, wsKline binance.WsKline, writeTime time.Time) {
	writer.initApiWriter()
	p := influxdb2.NewPoint("kline",
		map[string]string{"symbol": wsKline.Symbol, "interval": wsKline.Interval},
		adapter.BinanceKlineToMap(wsKline),
		writeTime)
	// write point asynchronously
	writer.apiWriter.WritePoint(p)
}

func (writer *InfluxWriter) initApiWriter() {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if writer.apiWriter == nil {
		writer.apiWriter = connector.GetInfluxClient().WriteAPI(config.Org, writer.Bucket)
	}
}
