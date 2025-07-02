package rexDatabase

import (
	"encoding/base64"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"strings"
	"time"
)

type PgDbConfig struct {
	Host        string
	User        string
	Password    string
	DbName      string
	Port        uint
	SslMode     string
	Loc         string
	Debug       bool
	TablePrefix string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime int64

	//Charset              string
	//LogMode              bool
	//AllowNativePasswords bool
}

func NewPgDbClient(c *PgDbConfig) (*gorm.DB, error) {
	// note: 对密码进行base64解码
	decodedBytes, err := base64.StdEncoding.DecodeString(c.Password)
	if err != nil {
		log.Fatalln("base64 decode error:", err)
		return nil, err
	}
	realPassword := string(decodedBytes)

	format := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s connect_timeout=30"
	dsn := fmt.Sprintf(format,
		c.Host,
		c.User,
		realPassword,
		c.DbName,
		c.Port,
		c.SslMode,
		c.Loc)

	// 处理表前缀
	var newLogger logger.Interface
	newLogger = NewGormZapLogger()
	if c.Debug {
		newLogger = newLogger.LogMode(logger.Info)
	} else {
		newLogger = newLogger.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s_", c.TablePrefix), // 表名前缀，`User`表为`t_users`
			SingularTable: true,                              // 使用单数表名，启用该选项后，`User` 表将是`auth`
			NameReplacer:  strings.NewReplacer("/", "_"),     // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
		},
	})

	if err != nil {
		return nil, err
	}

	sqlDB, sqlErr := db.DB()
	if c.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(c.MaxLifetime))
	} else {
		sqlDB.SetConnMaxLifetime(time.Second * 1000)
	}
	sqlDB.SetMaxOpenConns(c.MaxOpen)
	sqlDB.SetMaxIdleConns(c.MaxIdle)

	if sqlErr != nil {
		return nil, sqlErr
	}

	// Ping
	if pingErr := sqlDB.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
