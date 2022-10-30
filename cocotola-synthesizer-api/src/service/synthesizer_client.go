//go:generate mockery --output mock --name SynthesizerClient
package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/domain"
)

type SynthesizerClient interface {
	Synthesize(ctx context.Context, lang5 domain.Lang5, text string) (string, error)
}
