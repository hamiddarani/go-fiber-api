package middlewares

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/token"
	"github.com/hamiddarani/web-api-fiber/utils"
)

func Authentication(cfg *config.Config) fiber.Handler {
	var tokenService = token.NewTokenService(cfg.JWT)
	return func(c *fiber.Ctx) error {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.Request().Header.Peek(utils.AuthorizationHeaderKey)
		token := strings.TrimPrefix(string(auth), "Bearer ")

		if len(auth) == 0 {
			response := "please provide your authentication information"
			return c.Status(http.StatusUnauthorized).SendString(response)
		} else {
			claimMap, err = tokenService.GetClaims(token)
			if err != nil {
				response := "invalid token header, please login again"
				return c.Status(http.StatusUnauthorized).SendString(response)
			}
		}
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(utils.GenerateBaseResponseWithError(
				nil, false, -2, err,
			))
		}

		c.Locals(utils.UserIdKey, claimMap[utils.UserIdKey])
		c.Locals(utils.MobileNumberKey, claimMap[utils.MobileNumberKey])
		c.Locals(utils.RolesKey, claimMap[utils.RolesKey])
		c.Locals(utils.ExpireTimeKey, claimMap[utils.ExpireTimeKey])

		return c.Next()
	}
}

func Authorization(validRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Request().Header.PeekKeys()) == 0 {
			return c.Status(http.StatusForbidden).JSON(utils.GenerateBaseResponse(nil, false, -403))
		}
		rolesVal := c.Locals(utils.RolesKey)
		if rolesVal == nil {
			return c.Status(http.StatusForbidden).JSON(utils.GenerateBaseResponse(nil, false, -403))

		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}

		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				return c.Next()
			}
		}
		return c.Status(http.StatusForbidden).JSON(utils.GenerateBaseResponse(nil, false, -403))
	}
}
