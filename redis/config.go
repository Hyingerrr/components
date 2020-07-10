package redis

type Config struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"database"`
	DialTimeout  int    `yaml:"dial_timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	MaxRetries   int    `yaml:"max_retries"`
	Metrics      bool   `yaml:"metrics"`
}

type MgrRedisConfig map[string]*Config

func (c *Config) fillDefaultConfig(name string) {
	if c.DialTimeout == 0 {
		c.DialTimeout = 5000
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 1000
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 1000
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 3
	}
}
