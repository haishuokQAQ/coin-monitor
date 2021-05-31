package service

import (
	"context"
	"fmt"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/dal/mongodb"
)

func AnalyseHolder() {
	ctx := context.Background()
	// 启动时尝试读取热配置
	conf, err := mongodb.ReadHolderAnalyseConfig(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
}
