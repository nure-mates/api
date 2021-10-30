package redis

import (
	"sync"

	"github.com/go-redis/redis"

	"github.com/nure-mates/api/src/config"
)

var (
	client *Client
	once   = &sync.Once{}
)

// New - create redis client.
func New(cfg *config.Redis) (client *Client, err error) {
	once.Do(func() {
		cli := redis.NewClient(
			&redis.Options{
				Password: cfg.Password,
				PoolSize: cfg.PoolSize,
				Addr:     cfg.Address,
			})

		err = cli.Ping().Err()
		if err != nil {
			return
		}

		client = &Client{cli: cli}
	})

	return
}

// GetRedis returns redis client.
func GetRedis() RedisCli {
	return client
}
