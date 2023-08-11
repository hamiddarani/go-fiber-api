package config

import (
	"github.com/hamiddarani/web-api-fiber/pkg/cache"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
	"github.com/hamiddarani/web-api-fiber/pkg/token"
	"github.com/hamiddarani/web-api-fiber/utils/common"
)

func Default() *Config {
	return &Config{
		Logger: &logging.Config{
			FilePath: "./logs/",
			Encoding: "json",
			Logger:   "zap",
			Level:    "debug",
		},
		Postgres: &db.Config{
			Host:            "localhost",
			Port:            "5432",
			User:            "postgres",
			Password:        "admin",
			SSLMode:         "disable",
			DbName:          "test",
			MaxIdleConns:    15,
			ConnMaxLifetime: 5,
			MaxOpenConns:    100,
		},
		Redis: &cache.Config{
			Host:         "localhost",
			Port:         "6379",
			Password:     "password",
			Db:           "0",
			DialTimeout:  10,
			ReadTimeout:  10,
			WriteTimeout: 10,
			PoolSize:     15,
			PoolTimeout:  10,
		},
		JWT: &token.Config{
			AccessTokenExpireDuration:  1440,
			RefreshTokenExpireDuration: 60,
			Secret:                     "MySecret",
			RefreshSecret:              "MyRefreshSecret",
		},
		Otp: &common.Config{
			ExpireTime: 120,
			Limiter:    100,
			Digits:     6,
		},
	}
}
