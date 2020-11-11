/* ######################################################################
# Author: (jzy20admin@qq.com)
# Created Time: 2020-11-10 15:25:26
# File Name: client_mgr.go
# Description:
####################################################################### */

package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	once  sync.Once
	lock  sync.RWMutex
	pools map[string]*xorm.EngineGroup
)

func init() {
	pools = map[string]*xorm.EngineGroup{}
}

type mysqlConfig struct {
	Cfgs map[string]*Cfg `toml:"mysql"`
}

type Cfg struct {
	// dial
	DialUser string `toml:"user"`
	DialPawd string `toml:"pawd"`
	DialHost string `toml:"host"`
	DialPort string `toml:"port"`
	DialName string `toml:"name"`

	Debug bool `toml:"debug"` // 是否显示sql语句

	// pool
	PoolMaxOpenConn     int           `toml:"pool_max_open_conn"`     // 最大连接数大小
	PoolMaxIdleConn     int           `toml:"pool_max_idle_conn"`     // 最大空闲的连接的个数
	PoolConnMaxLifetime time.Duration `toml:"pool_conn_max_lifetime"` // 连接的生命时间,超过此时间，连接将关闭后重新建立新的，0代表忽略相关判断,单位:second
}

// 验证mysql实例的配置正确性与连通性。
// 参数names是实例的名称列表，如果为空则检测所有配置的实例
func Valid(names ...string) (err error) {
	if len(names) == 0 {
		var cfgs map[string]*Cfg
		if cfgs, err = loadCfgs(); err != nil {
			return
		}
		for k, _ := range cfgs {
			names = append(names, k)
		}
	}
	for _, name := range names {
		var cli *xorm.EngineGroup
		cli, err = SafeClient(name)
		if err == nil {
			err = cli.Ping()
		}
		if err != nil {
			err = fmt.Errorf("mysql#%s is invalid, %s", name, err)
			return
		}
	}
	return
}

func DefaultClient() (r *xorm.EngineGroup) {
	return Client("default")
}

func DefaultPool() (r *xorm.EngineGroup) {
	return Pool("default")
}

func Client(name string) (r *xorm.EngineGroup) {
	return Pool(name)
}

func SafeClient(name string) (r *xorm.EngineGroup, err error) {
	return SafePool(name)
}

func Pool(name string) (r *xorm.EngineGroup) {
	var err error
	if r, err = getPool(name); err != nil {
		panic(err)
	}
	return
}

func SafePool(name string) (r *xorm.EngineGroup, err error) {
	return getPool(name)
}

func getPool(name string) (r *xorm.EngineGroup, err error) {
	lock.RLock()
	r = pools[name]
	lock.RUnlock()
	if r == nil {
		r, err = addPool(name)
	}
	return
}

func addPool(name string) (r *xorm.EngineGroup, err error) {
	var cfg *Cfg
	if cfg, err = loadCfg(name); err != nil {
		return
	}
	r = NewMysqlPool(cfg)

	lock.Lock()
	pools[name] = r
	lock.Unlock()
	return
}

func loadCfg(name string) (r *Cfg, err error) {
	var cfgs map[string]*Cfg
	if cfgs, err = loadCfgs(); err != nil {
		return
	}
	if r = cfgs[name]; r == nil {
		err = fmt.Errorf("mysql#%s not configed", name)
		return
	}
	return
}

func loadCfgs() (r map[string]*Cfg, err error) {
	r = map[string]*Cfg{}

	cfg := &mysqlConfig{}
	once.Do(func() {
		_, err = config.Load(cfg, options.WithOnChangeFn(func(cfg interface{}) {
			lock.Lock()
			defer lock.Unlock()
			pools = map[string]*xorm.EngineGroup{}
		}))
	})

	config.Get(cfg)
	if err == nil && (cfg == nil || cfg.Cfgs == nil || len(cfg.Cfgs) == 0) {
		err = fmt.Errorf("not configed")
	}
	if err != nil {
		err = fmt.Errorf("mysql load cfgs error, %s", err)
		return
	}
	r = cfg.Cfgs
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
