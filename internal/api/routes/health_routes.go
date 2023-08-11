package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/handlers"
)

func Health(r fiber.Router) {
	h := handlers.NewHealthHandler()

	r.Get("/", h.HealthCheck)
}
