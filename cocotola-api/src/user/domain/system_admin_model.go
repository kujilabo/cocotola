package domain

const SystemAdminID = 1

type SystemAdminModel interface {
	GetID() uint
	IsSystemAdminModel() bool
}

type systemAdminModel struct {
}

func NewSystemAdminModel() SystemAdminModel {
	return &systemAdminModel{}
}

func (s *systemAdminModel) GetID() uint {
	return SystemAdminID
}

func (s *systemAdminModel) IsSystemAdminModel() bool {
	return true
}
