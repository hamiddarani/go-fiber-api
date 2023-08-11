package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/handlers"
	"github.com/hamiddarani/web-api-fiber/internal/config"
)

func FilesRoutes(r fiber.Router, cfg *config.Config) {
	h := handlers.NewFileHandler(cfg)

	r.Post("/", h.Create)
	r.Delete("/:id", h.Delete)
	r.Get("/:id", h.GetById)
}
