package rexPgPool

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PgPool interface {
		GetPool() *pgxpool.Pool
		GetConn() *pgxpool.Conn
		Ping() error
		Close()
	}
	defaultPgPool struct {
		pool *pgxpool.Pool
	}
)

func NewPgPool(conf *pgxpool.Config) PgPool {
	// Create database connection
	connPool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}
	return &defaultPgPool{
		pool: connPool,
	}
}

func (d *defaultPgPool) GetPool() *pgxpool.Pool {
	return d.pool
}

func (d *defaultPgPool) GetConn() *pgxpool.Conn {
	connection, err := d.pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()
	return connection
}

func (d *defaultPgPool) Ping() error {
	// note: 设置5秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.GetConn().Ping(ctx)
}

func (d *defaultPgPool) Close() {
	d.pool.Close()
}
