//go:generate mockery --output mock --name CustomTranslationRepository
package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
)

// var ErrCustomTranslationNotFound = errors.New("azure translation not found")
// var ErrCustomTranslationAlreadyExists = errors.New("azure translation already exists")

type CustomTranslationRepository interface {
	Add(ctx context.Context, param domain.TranslationAddParameter) error

	Update(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos, param domain.TranslationUpdateParameter) error

	Remove(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) error

	FindByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error)

	FindByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)

	FindByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error)

	Contain(ctx context.Context, lang2 domain.Lang2, text string) (bool, error)
}
