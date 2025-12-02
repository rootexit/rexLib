package rexCasBin

import (
	"encoding/base64"
	"errors"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rootexit/rexLib/rexStore"
	"gorm.io/gorm"
)

const (
	DefaultCasBinPrefix       = ""
	DefaultCasBinTable        = "casbin_rule"
	DefaultCasBinConfPath     = "etc/rbac_model.conf"
	DefaultCasBinWatchChannel = "casbin"
)

type (
	CasBinConf struct {
		ConfPath       string   `json:",default=etc/rbac_model.conf"`
		WatcherChannel string   `json:",default=casbin"`
		WhiteList      []string `json:",default=[]"`
	}
	CasBinTool interface {
		GetAdapter() *gormadapter.Adapter
		GetEnforcer() *casbin.Enforcer
		SetFilter(in []gormadapter.Filter)
		SetRedisConfig(in *rexStore.RedisConfig)
		NewWatcher(watchChannel string, watchCallback func(enforcer *casbin.Enforcer, msg string)) error
		Bootstrap() error
		Reload() error
	}
	defaultCasBinTool struct {
		db            *gorm.DB
		adapter       *gormadapter.Adapter
		enforcer      *casbin.Enforcer
		prefix        string
		tableName     string
		confPath      string
		isFiltered    bool
		filter        []gormadapter.Filter
		isWatching    bool
		rc            *rexStore.RedisConfig
		watchChannel  string
		watchCallback func(enforcer *casbin.Enforcer, msg string)
	}
)

func NewCasBinTool(DB *gorm.DB, prefix, tableName, confPath string, rc *rexStore.RedisConfig, watchChannel string, watchCallback func(enforcer *casbin.Enforcer, msg string)) CasBinTool {
	return NewFullCasBinTool(DB, rc, prefix, tableName, confPath, false, []gormadapter.Filter{}, true, watchChannel, watchCallback)
}

func EasyNewCasBinWatchTool(DB *gorm.DB, rc *rexStore.RedisConfig, watchChannel string, watchCallback func(enforcer *casbin.Enforcer, msg string)) CasBinTool {
	return NewFullCasBinTool(DB, rc, DefaultCasBinPrefix, DefaultCasBinTable, DefaultCasBinConfPath, false, []gormadapter.Filter{}, true, watchChannel, watchCallback)
}

func EasyNewCasBinTool(DB *gorm.DB) CasBinTool {
	return NewFullCasBinTool(DB, nil, DefaultCasBinPrefix, DefaultCasBinTable, DefaultCasBinConfPath, false, []gormadapter.Filter{}, false, DefaultCasBinWatchChannel, func(enforcer *casbin.Enforcer, msg string) {})
}

func NewFullCasBinTool(DB *gorm.DB, rc *rexStore.RedisConfig, prefix, tableName, confPath string, isFiltered bool, filters []gormadapter.Filter, isWatching bool, watchChannel string, watchCallback func(enforcer *casbin.Enforcer, msg string)) CasBinTool {
	return &defaultCasBinTool{
		db:            DB,
		prefix:        prefix,
		tableName:     tableName,
		confPath:      confPath,
		isFiltered:    isFiltered,
		filter:        filters,
		isWatching:    isWatching,
		rc:            rc,
		watchChannel:  watchChannel,
		watchCallback: watchCallback,
	}
}

func (t *defaultCasBinTool) GetAdapter() *gormadapter.Adapter {
	return t.adapter
}

func (t *defaultCasBinTool) GetEnforcer() *casbin.Enforcer {
	return t.enforcer
}

func (t *defaultCasBinTool) SetFilter(in []gormadapter.Filter) {
	t.filter = in
	t.isFiltered = true
	return
}

func (t *defaultCasBinTool) SetRedisConfig(in *rexStore.RedisConfig) {
	t.rc = in
	return
}

func (t *defaultCasBinTool) SetWatcher(watchChannel string, watchCallback func(*casbin.Enforcer, string)) {
	t.watchChannel = watchChannel
	t.watchCallback = watchCallback
	return
}

func (t *defaultCasBinTool) SetWatcherCallback(watchCallback func(*casbin.Enforcer, string)) {
	t.watchCallback = watchCallback
	return
}

func (t *defaultCasBinTool) NewWatcher(watchChannel string, watchCallback func(enforcer *casbin.Enforcer, msg string)) error {
	t.watchChannel = watchChannel
	t.watchCallback = watchCallback
	decodedBytes, err := base64.StdEncoding.DecodeString(t.rc.Pass)
	if err != nil {
		return err
	}
	w, err := rediswatcher.NewWatcher(t.rc.Host, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: string(decodedBytes),
		},
		Channel: t.watchChannel,
		// Only exists in test, generally be true
		IgnoreSelf: true,
	})
	if err != nil {
		return err
	}
	if err = t.enforcer.SetWatcher(w); err != nil {
		return err
	}
	if err = w.SetUpdateCallback(func(msg string) {
		t.watchCallback(t.enforcer, msg)
	}); err != nil {
		return err
	}
	return nil
}

func (t *defaultCasBinTool) Bootstrap() error {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(t.db, t.prefix, t.tableName)
	if err != nil {
		return err
	}
	t.adapter = adapter
	e, err := casbin.NewEnforcer(t.confPath, adapter)
	if err != nil {
		return err
	}
	t.enforcer = e
	if t.isFiltered {
		err = t.enforcer.LoadFilteredPolicy(t.filter)
		if err != nil {
			return err
		}
	} else {
		err = t.enforcer.LoadPolicy()
		if err != nil {
			return err
		}
	}
	if t.isWatching {
		if t.watchCallback != nil {
			return errors.New("watch callback is nil")
		}
		if t.watchChannel == "" {
			return errors.New("watch callback is nil")
		}
		if err = t.NewWatcher(t.watchChannel, t.watchCallback); err != nil {
			return err
		}
	}
	return err
}

func (t *defaultCasBinTool) Reload() error {
	if t.isFiltered {
		if err := t.enforcer.LoadFilteredPolicy(t.filter); err != nil {
			return err
		}
	} else {
		if err := t.enforcer.LoadPolicy(); err != nil {
			return err
		}
	}
	return nil
}
