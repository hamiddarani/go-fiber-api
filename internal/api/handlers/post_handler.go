package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/services"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/utils"
)

type PostHandler struct {
	service *services.PostService
}

func NewPostHandler(cfg *config.Config) *PostHandler {
	return &PostHandler{
		service: services.NewPostService(cfg),
	}
}

// CreatePost godoc
// @Summary Create a post
// @Description Create a post
// @Tags Posts
// @Accept json
// @produces json
// @Param Request body dto.CreateUpdatePostRequest true "Create a post"
// @Success 201 {object} utils.BaseHttpResponse{result=dto.PostResponse} "Post response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/posts/ [post]
// @Security AuthBearer
func (h *PostHandler) Create(c *fiber.Ctx) error {
	req := &dto.CreateUpdatePostRequest{}

	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, err))
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, errors))
	}

	res, err := h.service.Create(c, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithAnyError(nil, false, utils.InternalError, err))
	}
	return c.Status(http.StatusCreated).JSON(utils.GenerateBaseResponse(res, true, utils.Success))
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post
// @Tags Posts
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Param Request body dto.CreateUpdatePostRequest true "Update a post"
// @Success 200 {object} utils.BaseHttpResponse{result=dto.PostResponse} "Post response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/posts/{id} [put]
// @Security AuthBearer
func (h *PostHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	req := &dto.CreateUpdatePostRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, err))
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, errors))
	}

	res, err := h.service.Update(c, id, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithAnyError(nil, false, utils.InternalError, err))
	}
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse(res, true, utils.Success))
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post
// @Tags Posts
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} utils.BaseHttpResponse "response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/posts/{id} [delete]
// @Security AuthBearer
func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(http.StatusNotFound).JSON(
			utils.GenerateBaseResponse(nil, false, utils.ValidationErrorCode))
	}

	err := h.service.Delete(c, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithAnyError(nil, false, utils.InternalError, err))
	}
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse(nil, true, utils.Success))
}

// GetPost godoc
// @Summary Get a post
// @Description Get a post
// @Tags Posts
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} utils.BaseHttpResponse{result=dto.PostResponse} "Post response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/posts/{id} [get]
// @Security AuthBearer
func (h *PostHandler) GetById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(http.StatusNotFound).JSON(
			utils.GenerateBaseResponse(nil, false, utils.ValidationErrorCode))
	}

	res, err := h.service.GetById(c, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithAnyError(nil, false, utils.InternalError, err))
	}
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse(res, true, utils.Success))
}

// GetPost godoc
// @Summary Get list of post
// @Description Get list of post
// @Tags Posts
// @Accept json
// @produces json
// @Success 200 {object} utils.BaseHttpResponse{result=dto.PostResponse} "Post response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/posts [get]
// @Security AuthBearer
func (h *PostHandler) List(c *fiber.Ctx) error {
	res, err := h.service.GetByFilter(c)
	if err != nil {
		return c.Status(utils.TranslateErrorToStatusCode(err)).JSON(
			utils.GenerateBaseResponseWithError(nil, false, utils.InternalError, err))

	}
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse(res, true, 0))
}
