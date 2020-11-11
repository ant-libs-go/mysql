/* ######################################################################
# Author: (jzy20admin@qq.com)
# Created Time: 2020-11-10 14:50:38
# File Name: client.go
# Description:
####################################################################### */

package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func NewMysqlPool(cfg *Cfg) *xorm.EngineGroup {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&autocommit=true",
		cfg.DialUser, cfg.DialPawd, cfg.DialHost, cfg.DialPort, cfg.DialName))
	if err != nil {
		panic(fmt.Sprintf("mysql connect errr: %s", err))
	}
	orm, err := xorm.NewEngineGroup(engine, []*xorm.Engine{})
	orm.ShowSQL(cfg.Debug)
	orm.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	if cfg.PoolMaxOpenConn > 0 {
		orm.SetMaxOpenConns(cfg.PoolMaxOpenConn)
	}
	if cfg.PoolMaxIdleConn > 0 {
		orm.SetMaxIdleConns(cfg.PoolMaxIdleConn)
	}
	if cfg.PoolConnMaxLifetime > 0 {
		orm.SetConnMaxLifetime(cfg.PoolConnMaxLifetime * time.Millisecond)
	}
	return orm
}
