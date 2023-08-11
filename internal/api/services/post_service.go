package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/models"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
)

type PostService struct {
	base *BaseService[models.Post, dto.CreateUpdatePostRequest, dto.CreateUpdatePostRequest, dto.PostResponse, dto.PostFilter]
}

func NewPostService(cfg *config.Config) *PostService {
	return &PostService{
		base: &BaseService[models.Post, dto.CreateUpdatePostRequest, dto.CreateUpdatePostRequest, dto.PostResponse, dto.PostFilter]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg.Logger),
			Preloads: []preload{preload{string: "Image"}},
		},
	}
}

func (s *PostService) Create(ctx *fiber.Ctx, req *dto.CreateUpdatePostRequest) (*dto.PostResponse, error) {
	return s.base.Create(ctx, req)
}

func (s *PostService) Update(ctx *fiber.Ctx, id int, req *dto.CreateUpdatePostRequest) (*dto.PostResponse, error) {
	return s.base.Update(ctx, id, req)
}

func (s *PostService) Delete(ctx *fiber.Ctx, id int) error {
	return s.base.Delete(ctx, id)
}

func (s *PostService) GetById(ctx *fiber.Ctx, id int) (*dto.PostResponse, error) {
	return s.base.GetById(ctx, id)
}

func (s *PostService) GetByFilter(ctx *fiber.Ctx) (*dto.PagedList[dto.PostResponse], error) {
	return s.base.GetByFilter(ctx)
}
