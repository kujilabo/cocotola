package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
)

type tatoebaLinkRepository struct {
	db *gorm.DB
	rf service.RepositoryFactory
}

type tatoebaLinkEntity struct {
	From int
	To   int
}

func (e *tatoebaLinkEntity) TableName() string {
	return "tatoeba_link"
}

func newTatoebaLinkRepository(db *gorm.DB, rf service.RepositoryFactory) service.TatoebaLinkRepository {
	return &tatoebaLinkRepository{
		db: db,
		rf: rf,
	}
}

func (r *tatoebaLinkRepository) Add(ctx context.Context, param service.TatoebaLinkAddParameter) error {
	sentenceRepo := r.rf.NewTatoebaSentenceRepository(ctx)
	fromContained, err := sentenceRepo.ContainsSentenceBySentenceNumber(ctx, param.GetFrom())
	if err != nil {
		return err
	}

	toContained, err := sentenceRepo.ContainsSentenceBySentenceNumber(ctx, param.GetTo())
	if err != nil {
		return err
	}

	if !fromContained || !toContained {
		return service.ErrTatoebaSentenceNotFound
	}

	entity := tatoebaLinkEntity{
		From: param.GetFrom(),
		To:   param.GetTo(),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		if err := libG.ConvertDuplicatedError(result.Error, service.ErrTatoebaLinkAlreadyExists); errors.Is(err, service.ErrTatoebaLinkAlreadyExists) {
			return liberrors.Errorf("failed to Add tatoebaLink. err: %w", err)
		}

		if err := libG.ConvertRelationError(result.Error, service.ErrTatoebaLinkSourceNotFound); errors.Is(err, service.ErrTatoebaLinkSourceNotFound) {
			// fmt.Printf("relation %v, %v\n", fromContained, toContained)
			// nothing
			return nil
		}

		return liberrors.Errorf("failed to Add tatoebaLink. err: %w", result.Error)
	}

	return nil
}
