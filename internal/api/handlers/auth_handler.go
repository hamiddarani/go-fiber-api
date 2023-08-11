package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/services"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: services.NewAuthService(cfg),
	}
}

// SendOtp godoc
// @Summary Send otp to user
// @Description Send otp to user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Request body dto.GetOtpRequestDto true "GetOtpRequest"
// @Success 201 {object} utils.BaseHttpResponse "Success"
// @Failure 400 {object} utils.BaseHttpResponse "Failed"
// @Failure 409 {object} utils.BaseHttpResponse "Failed"
// @Router /v1/auth/send-otp [post]
func (h *AuthHandler) SendOtp(c *fiber.Ctx) error {
	req := new(dto.GetOtpRequestDto)

	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, errors))
	}

	err := h.service.SendOtp(req)

	if err != nil {
		return c.Status(utils.TranslateErrorToStatusCode(err)).
			JSON(utils.GenerateBaseResponseWithError(nil, false, -1, err))
	}

	return c.Status(http.StatusCreated).
		JSON(utils.GenerateBaseResponse(nil, true, utils.Success))
}

// RegisterLoginByMobileNumber godoc
// @Summary RegisterLoginByMobileNumber
// @Description RegisterLoginByMobileNumber
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Request body dto.RegisterLoginByMobileRequest true "RegisterLoginByMobileRequest"
// @Success 201 {object} utils.BaseHttpResponse "Success"
// @Failure 400 {object} utils.BaseHttpResponse "Failed"
// @Failure 409 {object} utils.BaseHttpResponse "Failed"
// @Router /v1/auth/login [post]
func (h *AuthHandler) RegisterLoginByMobileNumber(c *fiber.Ctx) error {
	req := new(dto.RegisterLoginByMobileRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, errors))
	}

	token, err := h.service.RegisterLoginByMobileNumber(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.InternalError, err))
	}

	return c.Status(http.StatusCreated).
		JSON(utils.GenerateBaseResponse(token, true, utils.Success))
}

// RefreshToken godoc
// @Summary RefreshToken
// @Description RefreshToken
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Request body dto.RefreshToken true "RefreshToken"
// @Success 201 {object} utils.BaseHttpResponse "Success"
// @Failure 400 {object} utils.BaseHttpResponse "Failed"
// @Failure 409 {object} utils.BaseHttpResponse "Failed"
// @Router /v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(dto.RefreshToken)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, errors))
	}

	token, err := h.service.RefreshToken(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, err))
	}

	return c.Status(http.StatusCreated).
		JSON(utils.GenerateBaseResponse(token, true, utils.Success))
}
