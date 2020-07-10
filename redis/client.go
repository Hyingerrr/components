package redis

import (
	"github.com/Hyingerrr/components/logger"
	goredis "github.com/go-redis/redis/v7"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

var (
	RedisMgr   *Client
	onceClient sync.Once
)

type Options func(*Client)

type Client struct {
	mux         sync.RWMutex
	clientMaps  map[string]*goredis.Client
	configMaps  MgrRedisConfig
	log         logger.Logger
	stateTicker time.Duration
}

func WithConfigs(configs MgrRedisConfig) Options {
	return func(mgr *Client) {
		mgr.configMaps = configs
	}
}

func WithRedisLogger(log logger.Logger) Options {
	return func(mgr *Client) {
		mgr.log = log
	}
}

func NewMgr(opts ...Options) *Client {
	onceClient.Do(func() {
		RedisMgr = &Client{
			clientMaps:  make(map[string]*goredis.Client),
			configMaps:  make(map[string]*Config),
			log:         logger.New(),
			stateTicker: 10 * time.Second,
		}

		for _, opt := range opts {
			opt(RedisMgr)
		}

		RedisMgr.loading()
	})

	return RedisMgr
}

func (mgr *Client) loading() {
	for schema, config := range mgr.configMaps {
		config.fillDefaultConfig(schema)

		cli := mgr.newRedisConnect(schema, config)
		if _, err := cli.Ping().Result(); err != nil {
			mgr.log.Errorf("redis schema[%v] ping error", schema)
			panic(err)
		}

		// prometheus
		if config.Metrics {
			// 每条指令都统计
			cli.AddHook(NewHook(schema))

			// 间隔统计连接数等
			mgr.MonitorStats(cli)
		}

		mgr.log.Infof("redis[%s] init success!", schema)
	}
}

func (mgr *Client) newRedisConnect(schema string, config *Config) *goredis.Client {
	cli := goredis.NewClient(&goredis.Options{
		Network:      "tcp",
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  time.Duration(config.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Millisecond,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	mgr.mux.Lock()
	mgr.clientMaps[schema] = cli
	mgr.mux.Unlock()

	return cli
}

func (mgr *Client) GetClient(schema string) *goredis.Client {
	if cli := mgr.hasClient(schema); cli != nil {
		return cli
	}

	mgr.log.Infof("redis schema[%s] not exist, new connect now...", schema)

	// new client
	if cfg, ok := mgr.configMaps[schema]; ok {
		return mgr.newRedisConnect(schema, cfg)
	}

	return nil
}

func (mgr *Client) hasClient(schema string) *goredis.Client {
	mgr.mux.RLock()
	defer mgr.mux.RUnlock()
	return mgr.clientMaps[schema]
}

func (mgr *Client) MonitorStats(client *goredis.Client) {
	var (
		timer = time.NewTicker(mgr.stateTicker)
	)

	for {
		select {
		case <-timer.C:
			stats := client.PoolStats()
			GoRedisStats.With(prometheus.Labels{"name": client.String(), "pool": "idle_pool"}).Set(float64(stats.IdleConns))

			GoRedisStats.With(prometheus.Labels{"name": client.String(), "pool": "cur_pool"}).Set(float64(stats.TotalConns))
		}
	}
}
