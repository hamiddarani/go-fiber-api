package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/models"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
)

type FileService struct {
	base *BaseService[models.File, dto.CreateFileRequest, dto.UpdateFileRequest, dto.FileResponse, dto.FileFilter]
}

func NewFileService(cfg *config.Config) *FileService {
	return &FileService{
		base: &BaseService[models.File, dto.CreateFileRequest, dto.UpdateFileRequest, dto.FileResponse, dto.FileFilter]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg.Logger),
		},
	}
}

func (s *FileService) Create(ctx *fiber.Ctx, req *dto.CreateFileRequest) (*dto.FileResponse, error) {
	return s.base.Create(ctx, req)
}

func (s *FileService) Update(ctx *fiber.Ctx, id int, req *dto.UpdateFileRequest) (*dto.FileResponse, error) {
	return s.base.Update(ctx, id, req)
}

func (s *FileService) Delete(ctx *fiber.Ctx, id int) error {
	return s.base.Delete(ctx, id)
}

func (s *FileService) GetById(ctx *fiber.Ctx, id int) (*dto.FileResponse, error) {
	return s.base.GetById(ctx, id)
}
