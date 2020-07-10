package redis

import (
	"testing"
	"time"

	"github.com/Hyingerrr/components/logger"
	"github.com/stretchr/testify/assert"
)

var (
	rdsConfigs = MgrRedisConfig{
		"redis1": &Config{
			Addr:         "127.0.0.1:6379",
			Password:     "",
			DB:           0,
			DialTimeout:  1000,
			ReadTimeout:  200,
			WriteTimeout: 200,
			PoolSize:     1000,
			MinIdleConns: 100,
			MaxRetries:   2,
			Metrics:      true,
		},
	}
)

func TestRedisMgr_NewMgr(t *testing.T) {
	var (
		it     = assert.New(t)
		client *Client
		key    = "text1"
		val    = "val1"
	)

	it.NotPanics(func() {
		client = NewMgr(WithConfigs(rdsConfigs), WithRedisLogger(logger.New()))
	})

	redis, err := client.GetClient("redis1")
	it.NoError(err)

	redis.Set(key, val, time.Millisecond*200)
}
