package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Health Check
// @Tags health
// @Accept  json
// @Success 200 {object} utils.BaseHttpResponse "Success"
// @Failure 400 {object} utils.BaseHttpResponse "Failed"
// @Produce  json
// @Router /v1/health/ [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse("working...", true, utils.Success))
}
