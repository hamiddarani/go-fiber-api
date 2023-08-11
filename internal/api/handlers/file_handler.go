package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/services"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/utils"
)

type FileHandler struct {
	service *services.FileService
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		service: services.NewFileService(cfg),
	}
}

// CreateFile godoc
// @Summary Create a file
// @Description Create a file
// @Tags Files
// @Accept x-www-form-urlencoded
// @produces json
// @Param file formData dto.UploadFileRequest true "Create a file"
// @Param file formData file true "Create a file"
// @Success 201 {object} utils.BaseHttpResponse{result=dto.FileResponse} "File response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/files/ [post]
// @Security AuthBearer
func (h *FileHandler) Create(c *fiber.Ctx) error {
	var err error

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// upload := &dto.UploadFileRequest{}
	// if err := c.BodyParser(upload); err != nil {
	// 	return c.Status(http.StatusBadRequest).SendString(err.Error())
	// }

	// if errors := utils.ValidateStruct(upload); errors != nil {
	// 	return c.Status(http.StatusBadRequest).
	// 		JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, errors))
	// }

	req := dto.CreateFileRequest{}
	// req.Description = upload.Description
	req.MimeType = file.Header.Get("Content-Type")
	req.Directory = "uploads"
	req.Name, err = saveUploadFile(file, req.Directory)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, err))

	}

	res, err := h.service.Create(c, &req)
	if err != nil {
		return c.Status(utils.TranslateErrorToStatusCode(err)).JSON(utils.GenerateBaseResponseWithError(nil, false, -1, err))
	}
	return c.Status(http.StatusCreated).JSON(utils.GenerateBaseResponse(res, true, utils.Success))

}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file
// @Tags Files
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} utils.BaseHttpResponse "response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/files/{id} [delete]
// @Security AuthBearer
func (h *FileHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponseWithValidationError(nil, false, utils.ValidationErrorCode, nil))
	}
	file, err := h.service.GetById(c, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponse(nil, false, -1))

	}
	err = os.Remove(fmt.Sprintf("%s/%s", file.Directory, file.Name))
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(utils.GenerateBaseResponse(nil, false, -1000))
	}
	err = h.service.Delete(c, id)
	if err != nil {
		return c.Status(utils.TranslateErrorToStatusCode(err)).
			JSON(utils.GenerateBaseResponseWithError(nil, false, -1, err))
	}
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse(nil, true, 0))
}

// GetFile godoc
// @Summary Get a file
// @Description Get a file
// @Tags Files
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} utils.BaseHttpResponse{result=dto.FileResponse} "File response"
// @Failure 400 {object} utils.BaseHttpResponse "Bad request"
// @Router /v1/files/{id} [get]
// @Security AuthBearer
func (h *FileHandler) GetById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	res, _ := h.service.GetById(c, id)
	return c.Status(http.StatusOK).JSON(utils.GenerateBaseResponse(res, true, utils.Success))
}

func saveUploadFile(file *multipart.FileHeader, directory string) (string, error) {
	randFileName := uuid.New()
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileName := file.Filename
	fileNameArr := strings.Split(fileName, ".")
	fileExt := fileNameArr[len(fileNameArr)-1]
	fileName = fmt.Sprintf("%s.%s", randFileName, fileExt)
	dst := fmt.Sprintf("%s/%s", directory, fileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
