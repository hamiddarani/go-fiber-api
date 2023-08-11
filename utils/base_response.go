package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hamiddarani/web-api-fiber/internal/api/validations"
)

type BaseHttpResponse struct {
	Result           any                `json:"result"`
	Success          bool               `json:"success"`
	ResultCode       ResultCode         `json:"resultCode"`
	ValidationErrors *[]ValidationError `json:"validationErrors"`
	Error            any                `json:"error"`
}

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
}

var validate = validations.NewValidator()

func ValidateStruct[Td any](dto *Td) error {

	err := validate.Struct(dto)

	if err != nil {
		return err
	}
	return nil
}

func getValidationErrors(err error) *[]ValidationError {
	var ve validator.ValidationErrors
	var validationErrors []ValidationError

	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("%v", err)
			var el ValidationError
			el.Property = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Error()
			validationErrors = append(validationErrors, el)
		}
		return &validationErrors
	}
	return nil
}

func GenerateBaseResponse(result any, success bool, resultCode ResultCode) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:    success,
		ResultCode: resultCode,
	}
}

func GenerateBaseResponseWithError(result any, success bool, resultCode ResultCode, err error) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:    success,
		ResultCode: resultCode,
		Error:      err.Error(),
	}

}

func GenerateBaseResponseWithAnyError(result any, success bool, resultCode ResultCode, err any) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:    success,
		ResultCode: resultCode,
		Error:      err,
	}
}

func GenerateBaseResponseWithValidationError(result any, success bool, resultCode ResultCode, err error) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:          success,
		ResultCode:       resultCode,
		ValidationErrors: getValidationErrors(err),
	}
}