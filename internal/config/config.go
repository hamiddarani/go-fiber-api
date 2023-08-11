package config

import (
	"github.com/hamiddarani/web-api-fiber/pkg/cache"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
	"github.com/hamiddarani/web-api-fiber/pkg/token"
	"github.com/hamiddarani/web-api-fiber/utils/common"
)

type Config struct {
	Logger   *logging.Config `koanf:"logger"`
	Postgres *db.Config      `koanf:"postgres"`
	Redis    *cache.Config   `koanf:"redis"`
	JWT      *token.Config   `koanf:"jwt"`
	Otp      *common.Config  `koanf:"otp"`
}
