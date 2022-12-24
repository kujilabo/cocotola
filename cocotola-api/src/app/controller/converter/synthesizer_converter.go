package converter

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func ToAudioResponse(ctx context.Context, audio service.Audio) (*entity.AudioResponse, error) {
	audioModel := audio.GetAudioModel()
	e := &entity.AudioResponse{
		ID:      int(audioModel.GetID()),
		Lang2:   audioModel.GetLang2().String(),
		Text:    audioModel.GetText(),
		Content: audioModel.GetContent(),
	}

	if err := libD.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return e, nil
}
