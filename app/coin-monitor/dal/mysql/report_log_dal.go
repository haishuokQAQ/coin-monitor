package mysql

import (
	"context"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/db"
)

func InsertReportLog(ctx context.Context, reportLog *db.ReportLog) error {
	err := connector.GetMysqlConnector().Create(reportLog).Error
	if err != nil {
		return err
	}
	return nil
}
