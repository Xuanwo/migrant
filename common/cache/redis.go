package cache

import (
	"time"

	"github.com/go-redis/redis"
)

// Redis is a redis client.
type Redis struct {
	*redis.Client
}

// RedisOptions is options for redis client.
type RedisOptions struct {
	Address string
	DB      int
	Timeout int
}

// NewRedis creates new redis client.
func NewRedis(opt *RedisOptions) (r *Redis, err error) {
	r = &Redis{}
	r.Client = redis.NewClient(&redis.Options{
		Addr:         opt.Address,
		DB:           opt.DB,
		DialTimeout:  time.Duration(opt.Timeout) * time.Second,
		ReadTimeout:  time.Duration(opt.Timeout) * time.Second,
		WriteTimeout: time.Duration(opt.Timeout) * time.Second,
	})

	_, err = r.Ping().Result()
	if err != nil {
		return
	}

	return
}

// RedisSentinelOptions is options for redis client.
type RedisSentinelOptions struct {
	MasterName string
	Addresses  []string
	DB         int
	Timeout    int
}

// NewRedisSentinel creates new redis sentinel client.
func NewRedisSentinel(opt *RedisSentinelOptions) (r *Redis, err error) {
	r = &Redis{}
	r.Client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    opt.MasterName,
		SentinelAddrs: opt.Addresses,
		DB:            opt.DB,
		DialTimeout:   time.Duration(opt.Timeout) * time.Second,
		ReadTimeout:   time.Duration(opt.Timeout) * time.Second,
		WriteTimeout:  time.Duration(opt.Timeout) * time.Second,
	})

	_, err = r.Ping().Result()
	if err != nil {
		return
	}

	return
}
