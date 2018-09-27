package contexts

import (
	"fmt"

	"github.com/Xuanwo/migrant/common/cache"
	"github.com/Xuanwo/migrant/common/db"
	"github.com/Xuanwo/migrant/config"
)

// Global shared variables, mostly connections and configs.
var (
	MySQL *db.MySQL

	Redis *cache.Redis
)

// Setup setups the context for migrant.
func Setup(c *config.Config) (err error) {
	if c.MySQL != nil {
		MySQL, err = db.NewMySQL(&db.MySQLOptions{
			Address: fmt.Sprintf(
				"%s:%d",
				c.MySQL.Host,
				c.MySQL.Port,
			),
			Database:           c.MySQL.Database,
			User:               c.MySQL.User,
			Password:           c.MySQL.Password,
			ConnectionTimeout:  c.MySQL.Timeout,
			MaxConnections:     c.MySQL.MaxConnections,
			MaxIdleConnections: c.MySQL.MaxIdleConnections,
		})
		if err != nil {
			return
		}
	}

	// Set redis first, so we can always pick redis sentinel if both set.
	if c.Redis != nil {
		Redis, err = cache.NewRedis(&cache.RedisOptions{
			Address: fmt.Sprintf(
				"%s:%d", c.Redis.Host, c.Redis.Port,
			),
			DB:      c.Redis.DB,
			Timeout: c.Redis.Timeout,
		})
		if err != nil {
			return
		}
	}
	if c.RedisSentinel != nil {
		Redis, err = cache.NewRedisSentinel(&cache.RedisSentinelOptions{
			Addresses:  c.RedisSentinel.Addresses,
			MasterName: c.RedisSentinel.MasterName,
			DB:         c.RedisSentinel.DB,
			Timeout:    c.RedisSentinel.Timeout,
		})
		if err != nil {
			return
		}
	}

	return nil
}
