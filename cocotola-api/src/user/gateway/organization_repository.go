package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
)

type organizationRepository struct {
	db *gorm.DB
}

type organizationEntity struct {
	SinmpleModelEntity
	Name string
}

func (e *organizationEntity) TableName() string {
	return "organization"
}

func (e *organizationEntity) toModel() (service.Organization, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	organizationModel, err := domain.NewOrganizationModel(model, e.Name)
	if err != nil {
		return nil, err
	}

	return service.NewOrganization(organizationModel)
}

func NewOrganizationRepository(ctx context.Context, db *gorm.DB) (service.OrganizationRepository, error) {
	if db == nil {
		return nil, liberrors.Errorf("db is inl. err: %w", libD.ErrInvalidArgument)
	}

	return &organizationRepository{
		db: db,
	}, nil
}

func (r *organizationRepository) GetOrganization(ctx context.Context, operator domain.AppUserModel) (service.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.GetOrganization")
	defer span.End()

	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			ID: uint(operator.GetOrganizationID()),
		},
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) FindOrganizationByName(ctx context.Context, operator domain.SystemAdminModel, name string) (service.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.FindOrganizationByName")
	defer span.End()

	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		Name: name,
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) FindOrganizationByID(ctx context.Context, operator domain.SystemAdminModel, id domain.OrganizationID) (service.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.FindOrganizationByID")
	defer span.End()

	organization := organizationEntity{}

	if result := r.db.Where(organizationEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			ID: uint(id),
		},
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}
		return nil, result.Error
	}

	return organization.toModel()
}

func (r *organizationRepository) AddOrganization(ctx context.Context, operator domain.SystemAdminModel, param service.OrganizationAddParameter) (domain.OrganizationID, error) {
	_, span := tracer.Start(ctx, "organizationRepository.AddOrganization")
	defer span.End()

	organization := organizationEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		Name: param.GetName(),
	}

	if result := r.db.Create(&organization); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrOrganizationAlreadyExists)
	}

	return domain.OrganizationID(organization.ID), nil
}
