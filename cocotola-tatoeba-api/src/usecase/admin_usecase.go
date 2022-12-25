package usecase

import (
	"context"
	"errors"
	"io"

	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

const (
	commitSize = 1000
	logSize    = 100000
)

type AdminUsecase interface {
	ImportSentences(ctx context.Context, iterator service.TatoebaSentenceAddParameterIterator) error

	ImportLinks(ctx context.Context, iterator service.TatoebaLinkAddParameterIterator) error
}

type adminUsecase struct {
	transaction service.Transaction
}

func NewAdminUsecase(transaction service.Transaction) AdminUsecase {
	return &adminUsecase{
		transaction: transaction,
	}
}

func (u *adminUsecase) ImportSentences(ctx context.Context, iterator service.TatoebaSentenceAddParameterIterator) error {
	logger := log.FromContext(ctx)

	var readCount = 0
	var importCount = 0
	var skipCount = 0
	var loop = true
	for loop {
		if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
			repo := rf.NewTatoebaSentenceRepository(ctx)

			i := 0
			for {
				param, err := iterator.Next(ctx)
				if errors.Is(err, io.EOF) {
					loop = false
					break
				}
				readCount++
				if err != nil {
					return liberrors.Errorf("read next line. read count: %d, err: %w", readCount, err)
				}

				if param == nil {
					skipCount++
					continue
				}

				if err := repo.Add(ctx, param); err != nil {
					logger.Warnf("failed to Add. read count: %d, err: %v", readCount, err)
					skipCount++
					continue
				}

				i++
				importCount++
				if i >= commitSize {
					if importCount%logSize == 0 {
						logger.Infof("imported count: %d", importCount)
					}
					break
				}
			}

			return nil
		}); err != nil {
			return liberrors.Errorf("import sentence. err: %w", err)
		}
	}

	logger.Infof("imported count: %d", importCount)
	logger.Infof("skipped count: %d", skipCount)
	logger.Infof("read count: %d", readCount)

	return nil
}

func (u *adminUsecase) ImportLinks(ctx context.Context, iterator service.TatoebaLinkAddParameterIterator) error {
	logger := log.FromContext(ctx)

	var readCount = 0
	var importCount = 0
	var skipCount = 0
	var loop = true
	for loop {
		if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
			repo := rf.NewTatoebaLinkRepository(ctx)

			i := 0
			for {
				param, err := iterator.Next(ctx)
				if errors.Is(err, io.EOF) {
					loop = false
					break
				}
				readCount++
				if err != nil {
					return liberrors.Errorf("read next line. read count: %d, err: %w", readCount, err)
				}
				if param == nil {
					skipCount++
					continue
				}

				if err := repo.Add(ctx, param); err != nil {
					if !errors.Is(err, service.ErrTatoebaSentenceNotFound) {
						logger.Warnf("failed to Add. read count: %d, err: %v", readCount, err)
					}
					skipCount++
					continue
				}
				i++
				importCount++
				if i >= commitSize {
					if importCount%logSize == 0 {
						logger.Infof("imported count: %d", importCount)
					}
					break
				}
			}

			return nil
		}); err != nil {
			return liberrors.Errorf("import sentence. err: %w", err)
		}
	}

	logger.Infof("imported count: %d", importCount)
	logger.Infof("skipped count: %d", skipCount)
	logger.Infof("read count: %d", readCount)

	return nil
}
