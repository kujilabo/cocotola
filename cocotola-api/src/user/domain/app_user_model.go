//go:generate mockery --output mock --name AppUserModel
package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type AppUserID uint

type AppUserModel interface {
	Model
	GetAppUserID() AppUserID
	GetOrganizationID() OrganizationID
	GetLoginID() string
	GetUsername() string
	GetRoles() []string
	GetProperties() map[string]string
}

type appUserModel struct {
	Model
	OrganizationID OrganizationID `validate:"required,gte=1"`
	LoginID        string         `validate:"required"`
	Username       string         `validate:"required"`
	Roles          []string
	Properties     map[string]string
}

func NewAppUserModel(model Model, organizationID OrganizationID, loginID, username string, roles []string, properties map[string]string) (AppUserModel, error) {
	m := &appUserModel{
		Model:          model,
		OrganizationID: organizationID,
		LoginID:        loginID,
		Username:       username,
		Roles:          roles,
		Properties:     properties,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *appUserModel) GetAppUserID() AppUserID {
	return (AppUserID)(m.GetID())
}

func (m *appUserModel) GetOrganizationID() OrganizationID {
	return m.OrganizationID
}

func (m *appUserModel) GetLoginID() string {
	return m.LoginID
}

func (m *appUserModel) GetUsername() string {
	return m.Username
}

func (m *appUserModel) GetRoles() []string {
	return m.Roles
}

func (m *appUserModel) GetProperties() map[string]string {
	return m.Properties
}
