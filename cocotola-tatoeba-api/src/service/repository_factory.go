package service

import (
	"context"
)

type RepositoryFactory interface {
	NewTatoebaLinkRepository(ctx context.Context) TatoebaLinkRepository

	NewTatoebaSentenceRepository(ctx context.Context) TatoebaSentenceRepository
}
