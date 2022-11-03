package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type Organization interface {
	domain.OrganizationModel
}

type organization struct {
	domain.OrganizationModel
}

func NewOrganization(organizationModel domain.OrganizationModel) (Organization, error) {
	m := &organization{
		organizationModel,
	}
	return m, libD.Validator.Struct(m)
}
