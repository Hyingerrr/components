package mysql

import (
"fmt"
"strings"
"sync"
"time"

"github.com/Hyingerrr/components/logger"
"github.com/prometheus/client_golang/prometheus"

_ "github.com/go-sql-driver/mysql" // 必须引入！！！
"github.com/jinzhu/gorm"
)

var (
	onceClient     sync.Once
	clientInstance *Client
)

type Option func(*Client)

type Client struct {
	log         logger.Logger
	clients     sync.Map
	configs     ManagerConfig
	stateTicker time.Duration // 上报prometheus的间隔, 默认10s
}

func WithDBConfig(cfgs ManagerConfig) Option {
	return func(client *Client) {
		client.configs = cfgs
	}
}

func WithDBLogger(logger logger.Logger) Option {
	return func(client *Client) {
		client.log = logger
	}
}

func NewClient(options ...Option) *Client {
	onceClient.Do(func() {
		clientInstance = &Client{
			log:         logger.New(),
			configs:     make(map[string]*Config),
			stateTicker: 10 * time.Second,
		}

		// 注入组件
		for _, opt := range options {
			opt(clientInstance)
		}

		clientInstance.loading()
	})

	return clientInstance
}

func (c *Client) loading() {
	for schema, config := range c.configs {
		// 加载默认配置
		config.fillWithDefaults()

		// 开始连接
		client, err := c.NewMysqlConnect(schema, config)
		if err != nil {
			panic(err)
		}

		// prometheus
		if config.Metrics {
			// 每次执行sql都要上报
			c.RegisterMetricsCallbacks(client)
			// 间隔统计db
			go c.DBStats(client)
		}

		// todo tracing
		// go func()

		c.log.Infof("mysql[%s] init success!", schema)
	}
}

func (c *Client) NewMysqlConnect(schema string, config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(config.Driver, config.DSN)
	if err != nil {
		c.log.Panicf("gorm open panic, schema[%v], error: %+v", schema, err)
		return nil, err
	}

	if config.MaxOpenConns > 0 {
		db.DB().SetMaxOpenConns(config.MaxOpenConns)
	}

	if config.MaxIdleConns > 0 {
		db.DB().SetMaxIdleConns(config.MaxIdleConns)
	}

	if err := db.DB().Ping(); err != nil {
		c.log.Panicf("mysql init panic, schema[%v], error:%+v", schema, err)
	}

	// 存储多例db
	c.setDbClient(schema, db)

	if config.Debug {
		db.LogMode(true)
	}

	return db, nil
}

func (c *Client) setDbClient(schema string, db *gorm.DB) {
	c.clients.Store(strings.ToLower(schema), db)
}

// 获取实例
func (c *Client) GetDb(schema string) (*gorm.DB, error) {
	db, ok := c.clients.Load(strings.ToLower(schema))
	if ok {
		if client, ok := db.(*gorm.DB); ok {
			return client, nil
		}
	}

	c.log.Errorf("db instance not exist! schema[%v], now start new ...", schema)

	// schema不存在 新建实例
	if cfg, ok := c.configs[schema]; ok {
		return c.NewMysqlConnect(schema, cfg)
	}

	return nil, fmt.Errorf("get db error, schema[%v]", schema)
}

// ping
func (c *Client) Ping() []error {
	var (
		errs []error
	)

	c.clients.Range(func(key, db interface{}) bool {
		if client, ok := db.(*gorm.DB); ok {
			if err := client.DB().Ping(); err != nil {
				errs = append(errs, err)
			}
		}
		return true
	})

	return errs
}

// 统计mysql连接数
func (c *Client) DBStats(client *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			c.log.Errorf("DBStats metrics error:%+v", err)
		}
	}()

	timer := time.NewTicker(c.stateTicker)

	for {
		select {
		case <-timer.C:
			stats := client.DB().Stats()
			schema := client.Dialect().CurrentDatabase() // database

			// db最大连接数
			DBStats.With(prometheus.Labels{"schema": schema, "stats": "max_open_conn"}).Set(float64(stats.MaxOpenConnections))

			// 已建立的连接数（含空闲和在使用的连接）
			DBStats.With(prometheus.Labels{"schema": schema, "stats": "open_conn"}).Set(float64(stats.OpenConnections))

			// 当前正在使用的连接数
			DBStats.With(prometheus.Labels{"schema": schema, "stats": "in_use"}).Set(float64(stats.InUse))

			// 空闲连接数
			DBStats.With(prometheus.Labels{"schema": schema, "stats": "idle"}).Set(float64(stats.Idle))

			// 处于等待状态的连接
			DBStats.With(prometheus.Labels{"schema": schema, "stats": "wait_count"}).Set(float64(stats.WaitCount))
		}
	}

}

