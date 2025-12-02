package rexRbac

import (
	"encoding/base64"
	"errors"
	"log"

	rediswatcher "github.com/billcobbler/casbin-redis-watcher/v2"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rootexit/rexLib/rexStore"
	"gorm.io/gorm"
)

const (
	defaultDatabaseName = "casbin"
	defaultTableName    = "casbin_rule"
)

type CasbinEngine struct {
	Prefix           string
	TableName        string
	RbacPath         string `json:"rbacPath"`
	RbacChannel      string `json:"rbacChannel"`
	IsCustomCallback bool
	UpdateCallback   func(*casbin.Enforcer, string)
	IsFiltered       bool
	Filter           []gormadapter.Filter
	Redis            *rexStore.RedisConfig
	Db               *gorm.DB
	Enforcer         *casbin.Enforcer
	Watcher          persist.Watcher
}

// Deprecated: Use CasBinTool instead.
func EasyNewCasbinEngine(RbacPath, RbacChannel string, redis *rexStore.RedisConfig, db *gorm.DB, v ...any) *CasbinEngine {
	if RbacPath == "" {
		panic("CasbinEngine RbacPath is nil")
	}
	if RbacChannel == "" {
		panic("CasbinEngine RbacChannel is nil")
	}
	if redis == nil {
		panic("CasbinEngine redis is nil")
	}
	if db == nil {
		panic("CasbinEngine db is nil")
	}

	prefix := ""
	tableName := ""
	if len(v) == 0 {
		prefix = ""
		tableName = defaultTableName
	} else if len(v) == 1 {
		prefix = v[0].(string)
		tableName = defaultTableName
	} else if len(v) == 2 {
		prefix = v[0].(string)
		tableName = v[1].(string)
	} else {
		panic(errors.New("wrong parameters"))
	}

	enforcer, err := InitRbac(RbacPath, db, prefix, tableName)
	if err != nil {
		panic(err)
	}
	return &CasbinEngine{
		RbacPath:    RbacPath,
		RbacChannel: RbacChannel,
		Redis:       redis,
		Db:          db,
		Enforcer:    enforcer,
	}
}

// Deprecated: Use CasBinTool instead.
func NewCasbinEngine(params *CasbinEngine) *CasbinEngine {
	if params == nil {
		panic("CasbinEngine params is nil")
	}
	enforcer, err := InitRbac(params.RbacPath, params.Db, params.Prefix, params.TableName)
	if err != nil {
		panic(err)
	}
	return &CasbinEngine{
		RbacPath:         params.RbacPath,
		RbacChannel:      params.RbacChannel,
		IsCustomCallback: params.IsCustomCallback,
		UpdateCallback:   params.UpdateCallback,
		IsFiltered:       params.IsFiltered,
		Filter:           params.Filter,
		Redis:            params.Redis,
		Db:               params.Db,
		Enforcer:         enforcer,
	}
}

// Deprecated: Use CasBinTool instead.
func InitRbac(RbacPath string, Db *gorm.DB, prefix, tableName string) (*casbin.Enforcer, error) {
	if prefix == "" {
		prefix = ""
	}
	if tableName == "" {
		tableName = defaultTableName
	}
	adapter, err := gormadapter.NewAdapterByDBUseTableName(Db, prefix, tableName)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(RbacPath, adapter)
	if err != nil {
		return nil, err
	}
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Deprecated: Use CasBinTool instead.
func (engine *CasbinEngine) EasyNewWatcher() *CasbinEngine {
	// note: 如果是定制的过滤器模式下，不判断是否是过滤模式，也不判断是否是定制回调模式
	engine.IsCustomCallback = false
	engine.IsFiltered = false
	return engine.NewWatcher()
}

// Deprecated: Use CasBinTool instead.
func (engine *CasbinEngine) CustomFilterNewWatcher(filters []gormadapter.Filter) *CasbinEngine {
	// note: 如果是定制的过滤器模式下，不判断是否是过滤模式，也不判断是否是定制回调模式
	engine.IsCustomCallback = true
	engine.IsFiltered = true
	engine.Filter = filters
	engine.UpdateCallback = func(enforcer *casbin.Enforcer, msg string) {
		_ = enforcer.LoadFilteredPolicy(filters)
	}
	return engine.NewWatcher()
}

// Deprecated: Use CasBinTool instead.
func (engine *CasbinEngine) NewWatcher() *CasbinEngine {
	decodedBytes, err := base64.StdEncoding.DecodeString(engine.Redis.Pass)
	if err != nil {
		log.Fatalln("base64 decode error:", err)
		return nil
	}
	watcher, err := rediswatcher.NewWatcher(engine.Redis.Host, rediswatcher.Password(string(decodedBytes)), rediswatcher.Channel(engine.RbacChannel))
	if err != nil {
		panic(err)
	}
	err = engine.Enforcer.SetWatcher(watcher)
	if err != nil {
		panic(err)
	}
	err = engine.Enforcer.SavePolicy()
	if err != nil {
		panic(err)
	}
	if engine.IsCustomCallback {
		err = watcher.SetUpdateCallback(func(msg string) {
			engine.UpdateCallback(engine.Enforcer, msg)
		})
		if err != nil {
			panic(err)
		}
	}
	if engine.IsFiltered {
		err = engine.Enforcer.LoadFilteredPolicy(engine.Filter)
		if err != nil {
			panic(err)
		}
	} else {
		err = engine.Enforcer.LoadPolicy()
		if err != nil {
			panic(err)
		}
	}
	engine.Watcher = watcher
	return engine
}
