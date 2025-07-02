package rexStore

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	// ClusterType means redis cluster.
	ClusterType = "cluster"
	// NodeType means redis node.
	NodeType = "node"
	// Nil is an alias of redis.Nil.
	Nil = redis.Nil

	blockingQueryTimeout = 5 * time.Second
	readWriteTimeout     = 2 * time.Second
	defaultSlowThreshold = time.Millisecond * 100
	defaultPingTimeout   = time.Second
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
	// ErrEmptyDB is an error that indicates no redis key is set.
	ErrEmptyDB = errors.New("empty redis db")
)

type (
	RedisConfig struct {
		Host     string
		Type     string `json:",default=node,options=node|cluster"`
		User     string `json:",optional"`
		Pass     string `json:",optional"`
		DB       int    `json:",default=0,optional"`
		tls      bool   `json:",optional"`
		NonBlock bool   `json:",default=true"`
		// PingTimeout is the timeout for ping redis.
		PingTimeout time.Duration `json:",default=1s"`
	}
)

func NewRedisClient(rc *RedisConfig) (*redis.Client, error) {
	if err := rc.Validate(); err != nil {
		return nil, err
	}

	// note: 对密码进行base64解码
	decodedBytes, err := base64.StdEncoding.DecodeString(rc.Pass)
	if err != nil {
		log.Fatalln("base64 decode error:", err)
		return nil, err
	}

	var tlsConfig *tls.Config
	if rc.tls {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	opts := redis.Options{
		Addr:         rc.Host,              // Redis 地址
		Password:     string(decodedBytes), // 没有密码则留空
		DB:           rc.DB,                // 默认 DB
		ReadTimeout:  readWriteTimeout,
		WriteTimeout: readWriteTimeout,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
		TLSConfig:    tlsConfig,
	}
	//if rc.tls {
	//	opts.TLSConfig = &tls.Config{}
	//}

	return redis.NewClient(&opts), nil
}

func (rc *RedisConfig) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}
