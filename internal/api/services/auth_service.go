package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamiddarani/web-api-fiber/internal/api/dto"
	"github.com/hamiddarani/web-api-fiber/internal/api/models"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/db"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
	"github.com/hamiddarani/web-api-fiber/pkg/token"
	"github.com/hamiddarani/web-api-fiber/utils"
	"github.com/hamiddarani/web-api-fiber/utils/common"
	"gorm.io/gorm"
)

type AuthService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	database     *gorm.DB
	tokenService *token.TokenService
}

func NewAuthService(cfg *config.Config) *AuthService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg.Logger)
	return &AuthService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		tokenService: token.NewTokenService(cfg.JWT),
		otpService:   NewOtpService(cfg),
	}
}

func (s *AuthService) RegisterLoginByMobileNumber(req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := s.otpService.ValidateOtp(req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := s.existsByMobileNumber(req.MobileNumber)
	if err != nil {
		return nil, err
	}

	u := models.User{MobileNumber: req.MobileNumber}

	if exists {
		var user models.User
		err = s.database.
			Model(&models.User{}).
			Where("mobile_number = ?", u.MobileNumber).
			Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Role")
			}).
			Find(&user).Error
		if err != nil {
			return nil, err
		}
		tdto := dto.TokenDto{UserId: user.Id, MobileNumber: user.MobileNumber}

		if len(*user.UserRoles) > 0 {
			for _, ur := range *user.UserRoles {
				tdto.Roles = append(tdto.Roles, ur.Role.Name)
			}
		}

		token, err := s.tokenService.GenerateToken(&tdto)
		if err != nil {
			return nil, err
		}
		return token, nil

	}

	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return nil, err
	}

	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	var user models.User
	err = s.database.
		Model(&models.User{}).
		Where("mobile_number = ?", u.MobileNumber).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error
	if err != nil {
		return nil, err
	}
	tdto := dto.TokenDto{UserId: user.Id, MobileNumber: user.MobileNumber}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tdto.Roles = append(tdto.Roles, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (s *AuthService) RefreshToken(req *dto.RefreshToken) (*dto.AccessToken, error) {
	claimMap, err := s.tokenService.GetClaims(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	td := &dto.AccessToken{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()

	atc := jwt.MapClaims{}

	atc[utils.UserIdKey] = claimMap[utils.UserIdKey]
	atc[utils.MobileNumberKey] = claimMap[utils.MobileNumberKey]
	atc[utils.RolesKey] = claimMap[utils.RolesKey]
	atc[utils.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	td.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (s *AuthService) SendOtp(req *dto.GetOtpRequestDto) error {
	otp := common.GenerateOtp(s.cfg.Otp)
	err := s.otpService.SetOtp(otp, req.MobileNumber)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*) > 0").
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *AuthService) getDefaultRole() (roleId int, err error) {

	if err = s.database.Model(&models.Role{}).
		Select("id").
		Where("name = ?", utils.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
