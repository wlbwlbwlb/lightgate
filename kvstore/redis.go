package kvstore

import (
	"time"

	"github.com/wl955/lightgate/config"

	"github.com/gomodule/redigo/redis"
)

var RedisPool = newPool(config.TOML.Redis.Server, config.TOML.Redis.Password, 0)

//server="localhost:6379"
//password=""
//db=0
func newPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err = c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err = c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		// Other pool configuration not shown in this example.
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
