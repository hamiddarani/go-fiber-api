package db

import "time"

type Config struct {
	Host            string        `koanf:"host"`
	Port            string        `koanf:"port"`
	User            string        `koanf:"user"`
	Password        string        `koanf:"password"`
	DbName          string        `koanf:"db_name"`
	SSLMode         string        `koanf:"SSLMode"`
	MaxIdleConns    int           `koanf:"max_idle_connections"`
	MaxOpenConns    int           `koanf:"max_open_connections"`
	ConnMaxLifetime time.Duration `koanf:"connection_max_life_time"`
}
