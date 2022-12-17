package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

var jst *time.Location

func init() {
	jst = time.Now().Local().Location()
}

type userQuotaEntity struct {
	ID             string
	OrganizationID uint
	AppUserID      uint
	Date           time.Time
	Name           string
	Unit           string
	Count          int
}

func (e *userQuotaEntity) TableName() string {
	return "user_quota"
}

type userQuotaRepository struct {
	db *gorm.DB
}

func NewUserQuotaRepository(db *gorm.DB) (service.UserQuotaRepository, error) {
	return &userQuotaRepository{
		db: db,
	}, nil
}

func (r *userQuotaRepository) IsExceeded(ctx context.Context, organizationID userD.OrganizationID, appUserID userD.AppUserID, name string, unit service.QuotaUnit, limit int) (bool, error) {
	now := time.Now()
	var date time.Time
	if unit == "month" {
		date = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, jst)
	} else if unit == "day" {
		date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	} else {
		date = time.Date(1900, 1, 1, 0, 0, 0, 0, jst)
	}
	entity := userQuotaEntity{}
	if result := r.db.Where(&userQuotaEntity{
		OrganizationID: uint(organizationID),
		AppUserID:      uint(appUserID),
		Date:           date,
		Unit:           string(unit),
		Name:           name,
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	if entity.Count > limit {
		return true, nil
	}
	return false, nil
}

func (r *userQuotaRepository) Increment(ctx context.Context, organizationID userD.OrganizationID, appUserID userD.AppUserID, name string, unit service.QuotaUnit, limit int, count int) (bool, error) {
	now := time.Now()
	var date time.Time
	if unit == "month" {
		date = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	} else if unit == "day" {
		date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		date = time.Date(1900, 1, 1, 0, 0, 0, 0, now.Location())
	}
	entity := userQuotaEntity{}
	if result := r.db.Where(&userQuotaEntity{
		OrganizationID: uint(organizationID),
		AppUserID:      uint(appUserID),
		Date:           date,
		Unit:           string(unit),
		Name:           name,
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			entity.ID = libD.NewULID()
			entity.OrganizationID = uint(organizationID)
			entity.AppUserID = uint(appUserID)
			entity.Date = date
			entity.Name = name
			entity.Unit = string(unit)
			entity.Count = count
			if result := r.db.Create(entity); result.Error != nil {
				return false, result.Error
			}
			if entity.Count > limit {
				return true, nil
			}
			return false, nil
		}
		return false, result.Error
	}

	if result := r.db.Model(&userQuotaEntity{}).Where(&userQuotaEntity{
		OrganizationID: uint(organizationID),
		AppUserID:      uint(appUserID),
		Date:           date,
		Unit:           string(unit),
		Name:           name,
	}).UpdateColumn("count", gorm.Expr("count + ?", count)); result.Error != nil {
		return false, result.Error
	}
	if entity.Count > limit {
		return true, nil
	}
	return false, nil
}
