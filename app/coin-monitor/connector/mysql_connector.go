package connector

import (
	"fmt"

	"github.com/haishuokQAQ/coin-monitor/app/config"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"

	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

func InitMysql() {
	err := InitMysqlConnector(config.MysqlHost, config.MysqlPort, config.MysqlUserName, config.MysqlPasswd, config.MysqlDBName)
	if err != nil {
		panic(err)
	}
}

func InitMysqlConnector(ip string, port int, username, password, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?", username, password, ip, port, dbName)
	db, err := gorm.Open(mysql.Dialector{
		Config: &mysql.Config{
			DSN: dsn,
		},
	}, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		Plugins:                                  nil,
	})
	if err != nil {
		return err
	}
	db = db.Set("gorm:save_associations", false).Set("gorm:association_save_reference", false)
	gormDB = db
	return nil
}

func GetMysqlConnector() *gorm.DB {
	return gormDB
}
