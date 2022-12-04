package domain

type OwnerModel interface {
	AppUserModel
}

type ownerModel struct {
	AppUserModel
}

func NewOwner(appUser AppUserModel) OwnerModel {
	return &ownerModel{
		AppUserModel: appUser,
	}
}
