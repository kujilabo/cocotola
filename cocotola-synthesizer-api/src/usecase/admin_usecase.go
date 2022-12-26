//go:generate mockery --output mock --name AdminUsecase
package usecase

import (
	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/service"
)

type AdminUsecase interface {
}

type adminUsecase struct {
	transaction service.Transaction
}

func NewAdminUsecase(transaction service.Transaction) AdminUsecase {
	return &adminUsecase{
		transaction: transaction,
	}
}
