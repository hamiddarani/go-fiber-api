package services

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/models"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
	"github.com/hamiddarani/web-api-fiber/utils"
	"github.com/hamiddarani/web-api-fiber/utils/common"
	"github.com/hamiddarani/web-api-fiber/utils/service_errors"
	"gorm.io/gorm"
)

type preload struct {
	string
}

type BaseService[T any, Tc any, Tu any, Tr any, Tf any] struct {
	Database *gorm.DB
	Logger   logging.Logger
	Preloads []preload
}

func NewBaseService[T any, Tc any, Tu any, Tr any, Tf any](cfg *config.Config) *BaseService[T, Tc, Tu, Tr, Tf] {
	return &BaseService[T, Tc, Tu, Tr, Tf]{
		Database: db.GetDb(),
		Logger:   logging.NewLogger(cfg.Logger),
	}
}

func (s *BaseService[T, Tc, Tu, Tr, Tf]) Create(ctx *fiber.Ctx, req *Tc) (*Tr, error) {

	model, _ := common.TypeConverter[T](req)
	tx := s.Database.Begin()
	err := tx.
		Create(model).
		Error
	if err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	bm, _ := common.TypeConverter[models.BaseModel](model)
	return s.GetById(ctx, bm.Id)
}

func (s *BaseService[T, Tc, Tu, Tr, Tf]) Update(ctx *fiber.Ctx, id int, req *Tu) (*Tr, error) {

	updateMap, _ := common.TypeConverter[map[string]interface{}](req)
	snakeMap := map[string]interface{}{}
	for k, v := range *updateMap {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Locals(utils.UserIdKey).(float64)), Valid: true}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	model := new(T)
	tx := s.Database.Begin()
	if err := tx.Model(model).
		Where("id = ? and deleted_by is null", id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	return s.GetById(ctx, id)

}

func (s *BaseService[T, Tc, Tu, Tr, Tf]) Delete(ctx *fiber.Ctx, id int) error {
	tx := s.Database.Begin()

	model := new(T)

	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Locals(utils.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	if ctx.Locals(utils.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	if cnt := tx.
		Model(model).
		Where("id = ? and deleted_by is null", id).
		Updates(deleteMap).
		RowsAffected; cnt == 0 {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Update, service_errors.RecordNotFound, nil)
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (s *BaseService[T, Tc, Tu, Tr, Tf]) GetById(ctx *fiber.Ctx, id int) (*Tr, error) {
	model := new(T)
	db := Preload(s.Database, s.Preloads)
	err := db.
		Where("id = ? and deleted_by is null", id).
		First(model).
		Error
	if err != nil {
		return nil, err
	}
	return common.TypeConverter[Tr](model)
}

func (s *BaseService[T, Tc, Tu, Tr, Tf]) GetByFilter(ctx *fiber.Ctx) (*dto.PagedList[Tr], error) {
	res, err := Paginate[T, Tr, Tf](ctx, s.Preloads, s.Database)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewPagedList[T any](items *[]T, count int64, pageNumber int, pageSize int) *dto.PagedList[T] {
	pl := &dto.PagedList[T]{
		PageNumber: pageNumber,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviousPage = pl.PageNumber > 1

	return pl
}

// Paginate
func Paginate[T any, Tr any, Tf any](ctx *fiber.Ctx, preloads []preload, db *gorm.DB) (*dto.PagedList[Tr], error) {
	model := new(T)
	var items *[]T
	var rItems *[]Tr
	db = Preload(db, preloads)

	filter := new(Tf)
	ctx.QueryParser(filter)
	query := getQuery[Tf](filter)
	sort := getSort(ctx)

	offset := GetOffset(ctx)
	limit := GetPageSize(ctx)
	page := GetPageNumber(ctx)

	var totalRows int64 = 0

	db.
		Model(model).
		Where(query).
		Count(&totalRows)

	err := db.
		Where(query).
		Offset(offset).
		Limit(limit).
		Order(sort).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}
	rItems, err = common.TypeConverter[[]Tr](items)
	if err != nil {
		return nil, err
	}
	return NewPagedList(rItems, totalRows, page, limit), err

}

func getQuery[Tf any](filter *Tf) string {

	query := make([]string, 0)
	query = append(query, "deleted_by is null")

	values := reflect.ValueOf(*filter)
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {
		if !values.Field(i).IsZero() {
			switch types.Field(i).Tag.Get("type") {
			case "contains":
				query = append(query, fmt.Sprintf("%s ILike '%%%v%%'", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "notContains":
				query = append(query, fmt.Sprintf("%s not ILike '%%%v%%'", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "startsWith":
				query = append(query, fmt.Sprintf("%s ILike '%v%%'", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "endsWith":
				query = append(query, fmt.Sprintf("%s ILike '%%%v'", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "equal":
				query = append(query, fmt.Sprintf("%s = '%v'", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "notEqual":
				query = append(query, fmt.Sprintf("%s != '%s'", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "lessThan":
				query = append(query, fmt.Sprintf("%s < %s", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "lessThanOrEqual":
				query = append(query, fmt.Sprintf("%s <= %s", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "greaterThan":
				query = append(query, fmt.Sprintf("%s > %s", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			case "greaterThanOrEqual":
				query = append(query, fmt.Sprintf("%s >= %s", types.Field(i).Tag.Get("db_field"), values.Field(i).Interface()))
			}
		}
	}
	return strings.Join(query, " AND ")
}

// getSort
func getSort(ctx *fiber.Ctx) string {
	sort := ctx.Query("sort", "asc")
	if sort == "asc" || sort == "desc" {
		return fmt.Sprintf("%s %s", "id", sort)
	}
	return fmt.Sprintf("%s %s", "id", "asc")
}

// Preload
func Preload(db *gorm.DB, preloads []preload) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.string)
	}
	return db
}

func GetOffset(ctx *fiber.Ctx) int {
	return (GetPageNumber(ctx) - 1) * GetPageSize(ctx)
}

func GetPageSize(ctx *fiber.Ctx) int {
	limit := ctx.QueryInt("limit", 10)
	if limit > 100 {
		limit = 100
	}
	return limit
}

func GetPageNumber(ctx *fiber.Ctx) int {
	page := ctx.QueryInt("page", 1)
	return page
}
