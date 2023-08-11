package cache

import "time"

type Config struct {
	Host               string        `koanf:"host"`
	Port               string        `koanf:"port"`
	Password           string        `koanf:"password"`
	Db                 string        `koanf:"db"`
	DialTimeout        time.Duration `koanf:"dial_timeout"`
	ReadTimeout        time.Duration `koanf:"read_timeout"`
	WriteTimeout       time.Duration `koanf:"write_timeout"`
	PoolSize           int           `koanf:"pool_size"`
	PoolTimeout        time.Duration `koanf:"pool_timeout"`
	IdleCheckFrequency time.Duration `koanf:"idle_check_frequency"`
}
