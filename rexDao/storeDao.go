package rexDao

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/rootexit/rexLib/rexStore"
	"log"
	"time"
)

type (
	RedisDao interface {
		GetRD() *redis.Client
		Ping() error
		Close() error
		Set(key string, value interface{}) error
		SetCtx(ctx context.Context, key string, value interface{}) error
		SetEx(key string, value interface{}, seconds int) error
		SetExCtx(ctx context.Context, key string, value interface{}, seconds int) error
		Get(key string) (string, error)
		GetCtx(ctx context.Context, key string) (string, error)
		NewWatcher(channel string, fn func(msg *redis.Message)) error
		NewWatcherCtx(ctx context.Context, channel string, fn func(msg *redis.Message)) error
		Publish(channel string, msg string) error
		PublishCtx(ctx context.Context, channel string, msg string) error
		Ttl(key string) (int, error)
		TtlCtx(ctx context.Context, key string) (int, error)
		Del(keys ...string) (int, error)
		DelCtx(ctx context.Context, keys ...string) (int, error)
		Keys(pattern string) ([]string, error)
		KeysCtx(ctx context.Context, pattern string) ([]string, error)
	}
	defaultRedisDao struct {
		rd *redis.Client
	}
)

func NewRedisDaoWithRdConfig(rc *rexStore.RedisConfig) RedisDao {
	rdClient, err := rexStore.NewRedisClient(rc)
	if err != nil {
		log.Fatalf("connect redis store error: %v", err)
		return nil
	}
	return &defaultRedisDao{
		rd: rdClient,
	}
}

func NewRedisDao(rd *redis.Client) RedisDao {
	return &defaultRedisDao{
		rd: rd,
	}
}

func (d *defaultRedisDao) GetRD() *redis.Client {
	return d.rd
}

func (d *defaultRedisDao) Ping() error {
	// note: 设置5秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return d.rd.Ping(ctx).Err()
	}
	return nil
}

func (d *defaultRedisDao) Close() error {
	return d.rd.Close()
}

func (d *defaultRedisDao) Set(ckey string, value interface{}) error {
	return d.SetCtx(context.Background(), ckey, value)
}

func (d *defaultRedisDao) SetCtx(ctx context.Context, key string, value interface{}) error {
	if err := d.rd.Set(ctx, key, value, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (d *defaultRedisDao) SetEx(key string, value interface{}, seconds int) error {
	return d.SetExCtx(context.Background(), key, value, seconds)
}

func (d *defaultRedisDao) SetExCtx(ctx context.Context, key string, value interface{}, seconds int) error {
	if err := d.rd.SetEx(ctx, key, value, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}
	return nil
}

func (d *defaultRedisDao) GetCtx(ctx context.Context, key string) (string, error) {
	if val, err := d.rd.Get(ctx, key).Result(); errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func (d *defaultRedisDao) Get(key string) (string, error) {
	return d.GetCtx(context.Background(), key)
}

func (d *defaultRedisDao) NewWatcher(channel string, fn func(msg *redis.Message)) error {
	return d.NewWatcherCtx(context.Background(), channel, fn)
}

func (d *defaultRedisDao) NewWatcherCtx(ctx context.Context, channel string, fn func(msg *redis.Message)) error {
	sub := d.rd.Subscribe(ctx, channel)
	ch := sub.Channel()
	for msg := range ch {
		// 处理消息
		fn(msg)
	}
	return nil
}

func (d *defaultRedisDao) PublishCtx(ctx context.Context, channel string, msg string) error {
	if err := d.rd.Publish(ctx, channel, msg).Err(); err != nil {
		return err
	}
	return nil
}

func (d *defaultRedisDao) Publish(channel string, msg string) error {
	return d.PublishCtx(context.Background(), channel, msg)
}

func (d *defaultRedisDao) TtlCtx(ctx context.Context, key string) (int, error) {
	duration, err := d.rd.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if duration >= 0 {
		return int(duration / time.Second), nil
	}
	return 0, nil
}

func (d *defaultRedisDao) Ttl(key string) (int, error) {
	return d.TtlCtx(context.Background(), key)
}

func (d *defaultRedisDao) Del(keys ...string) (int, error) {
	return d.DelCtx(context.Background(), keys...)
}
func (d *defaultRedisDao) DelCtx(ctx context.Context, keys ...string) (int, error) {
	v, err := d.rd.Del(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (d *defaultRedisDao) Keys(pattern string) ([]string, error) {
	return d.KeysCtx(context.Background(), pattern)
}

// KeysCtx is the implementation of redis keys command.
func (d *defaultRedisDao) KeysCtx(ctx context.Context, pattern string) ([]string, error) {
	return d.rd.Keys(ctx, pattern).Result()
}
