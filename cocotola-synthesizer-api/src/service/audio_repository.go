//go:generate mockery --output mock --name AudioRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/domain"
)

var ErrAudioNotFound = errors.New("Audio not found")

type AudioRepository interface {
	AddAudio(ctx context.Context, lang5 domain.Lang5, text, audioContent string) (domain.AudioID, error)

	FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (Audio, error)

	FindByLangAndText(ctx context.Context, lang5 domain.Lang5, text string) (Audio, error)

	FindAudioIDByText(ctx context.Context, lang5 domain.Lang5, text string) (domain.AudioID, error)
}
