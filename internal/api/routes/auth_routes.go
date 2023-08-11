package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/handlers"
	"github.com/hamiddarani/web-api-fiber/internal/config"
)

func AuthRoutes(r fiber.Router, cfg *config.Config) {
	h := handlers.NewAuthHandler(cfg)

	r.Post("/send-otp", h.SendOtp)

	r.Post("/login", h.RegisterLoginByMobileNumber)
	r.Post("/refresh-token", h.RefreshToken)
}
