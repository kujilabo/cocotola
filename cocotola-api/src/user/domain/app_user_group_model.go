package domain

import libD "github.com/kujilabo/cocotola/lib/domain"

type AppUserGroupID uint

type AppUserGroupModel interface {
	Model
	GetAppUerGroupID() AppUserGroupID
	GetOrganizationID() OrganizationID
	GetKey() string
	GetName() string
	GetDescription() string
}

type appUserGroupModel struct {
	Model
	OrganizationID OrganizationID
	Key            string `validate:"required"`
	Name           string `validate:"required"`
	Description    string
}

// NewAppUserGroup returns a new AppUserGroup
func NewAppUserGroup(model Model, organizationID OrganizationID, key, name, description string) (AppUserGroupModel, error) {
	m := &appUserGroupModel{
		Model:          model,
		OrganizationID: organizationID,
		Key:            key,
		Name:           name,
		Description:    description,
	}

	return m, libD.Validator.Struct(m)
}

func (m *appUserGroupModel) GetAppUerGroupID() AppUserGroupID {
	return AppUserGroupID(m.GetID())
}

func (m *appUserGroupModel) GetOrganizationID() OrganizationID {
	return m.OrganizationID
}

func (m *appUserGroupModel) GetKey() string {
	return m.Key
}

func (m *appUserGroupModel) GetName() string {
	return m.Name
}

func (m *appUserGroupModel) GetDescription() string {
	return m.Description
}
