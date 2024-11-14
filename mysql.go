/* ######################################################################
# Author: (jzy20admin@qq.com)
# Created Time: 2020-11-10 14:50:38
# File Name: mysql.go
# Description:
####################################################################### */

package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysqlPool(cfg *Cfg) *gorm.DB {
	gcfg := &gorm.Config{}
	if cfg.Debug == true {
		gcfg.Logger = logger.Default.LogMode(logger.Info)
	}
	orm, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.DialUser, cfg.DialPawd, cfg.DialHost, cfg.DialPort, cfg.DialName)), gcfg)
	if err != nil {
		panic(fmt.Sprintf("mysql connect errr: %s", err))
	}

	db, err := orm.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get DB instance: %s", err))
	}

	if cfg.PoolMaxOpenConn > 0 {
		db.SetMaxOpenConns(cfg.PoolMaxOpenConn)
	}
	if cfg.PoolMaxIdleConn > 0 {
		db.SetMaxIdleConns(cfg.PoolMaxIdleConn)
	}
	if cfg.PoolConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.PoolConnMaxLifetime * time.Millisecond)
	}
	return orm
}
