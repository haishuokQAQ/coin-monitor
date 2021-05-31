package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/dal/mysql"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/db"
)

func SaveReportLog(ctx context.Context, dataStr string, remoteIp string) {
	// 先base64解码
	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		fmt.Println(dataStr, err)
		return
	}
	actualByte, err := RsaDecrypt(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	reportLog := &db.ReportLog{
		SourceIp:        remoteIp,
		ReportTimestamp: uint64(time.Now().Unix()),
		ReportContent:   string(actualByte),
	}
	err = mysql.InsertReportLog(ctx, reportLog)
	if err != nil {
		fmt.Println(err)
		return
	}
}
