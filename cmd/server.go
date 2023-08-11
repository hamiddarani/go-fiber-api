package cmd

import (
	"log"
	"os"

	"github.com/hamiddarani/web-api-fiber/internal/api"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/cache"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	migrations "github.com/hamiddarani/web-api-fiber/pkg/db/migration"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
	"github.com/joho/godotenv"

	"github.com/spf13/cobra"
)

type Server struct{}

func (s Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		s.run(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "server",
		Short: "run PhoneBook server",
		Run:   run,
	}
}

func (s *Server) run(cfg *config.Config, trap chan os.Signal) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := logging.NewLogger(cfg.Logger)

	if err := db.InitDb(cfg.Postgres); err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	if err := cache.InitRedis(cfg.Redis); err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	migrations.Up_1()

	server := api.NewServer(cfg, logger)
	go server.InitServer()

	logger.Info(logging.OS, logging.SystemKill, "signal trap", map[logging.ExtraKey]interface{}{"signal": (<-trap).String()})
}
