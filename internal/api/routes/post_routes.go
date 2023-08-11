package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/handlers"
	"github.com/hamiddarani/web-api-fiber/internal/config"
)

func PostRoutes(r fiber.Router, cfg *config.Config) {
	h := handlers.NewPostHandler(cfg)

	r.Post("/", h.Create)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
	r.Get("/:id", h.GetById)
	r.Get("/", h.List)
}
