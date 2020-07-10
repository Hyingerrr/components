package mysql

import (
	"runtime"
	"time"
)

const (
	MaxDialTimeout   = 1000 // ms  最大连接超时间
	MaxReadTimeout   = 3000 // ms  读超时时间
	MaxWriteTimeout  = 5000 // ms  写超时时间
	MaxOpenConn      = 128
	MaxIdleConn      = 16
	MaxLifecycleConn = 300 // in second
)

type ManagerConfig map[string]*Config

type Config struct {
	Driver       string        `yaml:"driver"`
	DSN          string        `yaml:"dsn"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxLifeConns int           `yaml:"max_life_conns"`
	Debug        bool          `yaml:"debug"`   // 是否开启debug
	Metrics      bool          `yaml:"metrics"` // 是否上报prometheus
}

func (c *Config) fillWithDefaults() {
	var (
		maxCPU = runtime.NumCPU()
	)

	if c.DialTimeout <= 0 || c.DialTimeout > time.Duration(MaxDialTimeout*maxCPU) {
		c.DialTimeout = MaxDialTimeout
	}

	if c.ReadTimeout <= 0 || c.ReadTimeout > time.Duration(MaxReadTimeout*maxCPU) {
		c.ReadTimeout = MaxReadTimeout
	}

	if c.WriteTimeout <= 0 || c.WriteTimeout > time.Duration(MaxWriteTimeout*maxCPU) {
		c.WriteTimeout = MaxWriteTimeout
	}

	if c.MaxOpenConns <= 0 || c.MaxOpenConns > maxCPU*MaxOpenConn {
		c.MaxOpenConns = MaxOpenConn
	}

	if c.MaxIdleConns <= 0 || c.MaxIdleConns > maxCPU*MaxIdleConn {
		c.MaxIdleConns = MaxIdleConn
	}

	if c.MaxLifeConns <= 0 || c.MaxLifeConns > MaxLifecycleConn*maxCPU {
		c.MaxLifeConns = MaxLifecycleConn
	}
}
