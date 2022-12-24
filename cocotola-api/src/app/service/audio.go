//go:generate mockery --output mock --name Audio
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type Audio interface {
	GetAudioModel() domain.AudioModel
}

type audio struct {
	AudioModel domain.AudioModel
}

func NewAudio(audioModel domain.AudioModel) (Audio, error) {
	m := &audio{
		AudioModel: audioModel,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (s *audio) GetAudioModel() domain.AudioModel {
	return s.AudioModel
}
