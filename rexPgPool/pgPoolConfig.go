package rexPgPool

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPoolConfig struct {
	Host     string `json:",default=pg.xx.com"`
	User     string `json:",default=admin"`
	Password string `json:",default=admin123"`
	DbName   string `json:",default=public"`
	Port     uint   `json:",default=5432"`
	SslMode  string `json:",default=disable"`
	Loc      string `json:",default=Asia/Shanghai"`
	Debug    bool   `json:",default=true"`

	MaxConns          int32         `json:",default=4"`
	MinConns          int32         `json:",default=0"`
	MaxConnLifetime   time.Duration `json:",default=10m"`
	MaxConnIdleTime   time.Duration `json:",default=30m"`
	HealthCheckPeriod time.Duration `json:",default=1m"`
	ConnectTimeout    time.Duration `json:",default=5s"`
}

func DefaultPoolConfig(conf *PgPoolConfig) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Minute * 5
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// note: 对密码进行base64解码
	decodedBytes, err := base64.StdEncoding.DecodeString(conf.Password)
	if err != nil {
		log.Fatalln("base64 decode error:", err)
		return nil
	}
	realPassword := string(decodedBytes)

	// 自定义数据库 URL
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(conf.User, realPassword),
		Host:   fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:   conf.DbName,
	}

	q := u.Query()
	q.Set("sslmode", conf.SslMode)
	q.Set("options", "-c timezone="+conf.Loc)
	u.RawQuery = q.Encode()
	dsn := u.String()
	if conf.Debug {
		log.Println("Postgres connection url:", dsn)
	}

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("创建配置失败，错误：", err)
	}

	if conf.MaxConns != 0 {
		dbConfig.MaxConns = conf.MaxConns
	} else {
		dbConfig.MaxConns = defaultMaxConns
	}
	if conf.MaxConns != 0 {
		dbConfig.MinConns = conf.MaxConns
	} else {
		dbConfig.MaxConns = defaultMinConns
	}
	if conf.MaxConnLifetime != 0 {
		dbConfig.MaxConnLifetime = conf.MaxConnLifetime
	} else {
		dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	}
	if conf.MaxConnIdleTime != 0 {
		dbConfig.MaxConnIdleTime = conf.MaxConnIdleTime
	} else {
		dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	}
	if conf.HealthCheckPeriod != 0 {
		dbConfig.HealthCheckPeriod = conf.HealthCheckPeriod
	} else {
		dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	}

	if conf.ConnectTimeout != 0 {
		dbConfig.ConnConfig.ConnectTimeout = conf.ConnectTimeout
	} else {
		dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout
	}

	if conf.Debug {
		dbConfig.ConnConfig.Tracer = &QueryTracer{}
	}

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		if conf.Debug {
			log.Println("Before acquiring the connection pool to the database!!")
		}
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		if conf.Debug {
			log.Println("After releasing the connection pool to the database!!")
		}
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		if conf.Debug {
			log.Println("Closed the connection pool to the database!!")
		}
	}

	return dbConfig
}
