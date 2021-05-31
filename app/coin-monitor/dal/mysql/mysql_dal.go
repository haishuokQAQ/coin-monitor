package mysql

import (
	"context"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/db"

	"gorm.io/gorm"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
)

func ListPlatform(ctx context.Context) ([]*db.Platform, error) {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	result := []*db.Platform{}
	err := dbConn.Model(&db.Platform{}).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ListCoin(ctx context.Context) ([]*db.Coin, error) {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	result := []*db.Coin{}
	err := dbConn.Model(&db.Coin{}).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetCoinByPlatform(ctx context.Context, platformId uint64) ([]*db.Coin, error) {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	coinPlatformLinks := []*db.PlatformCoin{}
	err := dbConn.Model(&db.PlatformCoin{}).Where("platform_id = ?", platformId).Find(&coinPlatformLinks).Error
	if err != nil {
		return nil, err
	}
	coinNames := []string{}
	for _, link := range coinPlatformLinks {
		coinNames = append(coinNames, link.CoinName)
	}
	result := []*db.Coin{}
	err = dbConn.Model(&db.Coin{}).Where("name in (?)", coinNames).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetContractMetaById(ctx context.Context, address string) (*db.ContractInfo, error) {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	result := &db.ContractInfo{}
	err := dbConn.Model(result).Where("address = ?", address).First(result).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func SearchContractInfoForPage(ctx context.Context, pageNum, pageSize int) ([]*db.ContractInfo, int64, error) {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	result := []*db.ContractInfo{}
	count := int64(0)
	err := dbConn.Model(&db.ContractInfo{}).Where("symbol = ?", "").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = dbConn.Model(&db.ContractInfo{}).Where("symbol = ?", "").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

func InsertContractMeta(ctx context.Context, contractInfo *db.ContractInfo) error {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	err := dbConn.Create(contractInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateContractMeta(ctx context.Context, contractInfo *db.ContractInfo) error {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	err := dbConn.Where("address = ?", contractInfo.Address).Save(contractInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertTokenTransferHistory(ctx context.Context, transferHistory *db.TokenTransferHistory) error {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	err := dbConn.Create(transferHistory).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTokenTransferHistoryByTokenAndAddress(ctx context.Context, address string, tokenAddress string) (*db.TokenTransferHistory, error) {
	dbConn := connector.GetMysqlConnector().WithContext(ctx)
	result := &db.TokenTransferHistory{}
	err := dbConn.Model(result).Where("tx_address = ?", address).Where("token_address = ?", tokenAddress).First(result).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
