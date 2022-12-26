package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/domain"
	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type audioEntity struct {
	ID           uint   `validate:"required"`
	Lang5        string `validate:"required"`
	Text         string `validate:"required"`
	AudioContent string `validate:"required"`
}

func (e *audioEntity) TableName() string {
	return "audio"
}

func (e *audioEntity) toAudioModel() (domain.AudioModel, error) {
	lang5, err := domain.NewLang5(e.Lang5)
	if err != nil {
		return nil, err
	}

	return domain.NewAudioModel(e.ID, lang5, e.Text, e.AudioContent)
}

type audioRepository struct {
	db *gorm.DB
}

func newAudioRepository(db *gorm.DB) service.AudioRepository {
	return &audioRepository{
		db: db,
	}
}

func (r *audioRepository) AddAudio(ctx context.Context, lang5 domain.Lang5, text, audioContent string) (domain.AudioID, error) {
	_, span := tracer.Start(ctx, "audioRepository.AddAudio")
	defer span.End()

	entity := audioEntity{
		Lang5:        lang5.String(),
		Text:         text,
		AudioContent: audioContent,
	}
	if result := r.db.Create(&entity); result.Error != nil {
		return 0, result.Error
	}
	return domain.AudioID(entity.ID), nil
}

func (r *audioRepository) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error) {
	_, span := tracer.Start(ctx, "audioRepository.FindAudioByAudioID")
	defer span.End()

	entity := audioEntity{}
	if result := r.db.Where("id = ?", uint(audioID)).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAudioNotFound
		}
		return nil, result.Error
	}
	audioModel, err := entity.toAudioModel()
	if err != nil {
		return nil, err
	}
	audio, err := service.NewAudio(audioModel)
	if err != nil {
		return nil, err
	}
	return audio, nil
}

func (r *audioRepository) FindByLangAndText(ctx context.Context, lang5 domain.Lang5, text string) (service.Audio, error) {
	_, span := tracer.Start(ctx, "audioRepository.FindByLangAndText")
	defer span.End()

	entity := audioEntity{}
	if result := r.db.Where("lang5 = ? and text = ?", lang5.String(), text).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAudioNotFound
		}
		return nil, result.Error
	}
	audio, err := entity.toAudioModel()
	if err != nil {
		return nil, err
	}
	audioService, err := service.NewAudio(audio)
	if err != nil {
		return nil, err
	}
	return audioService, nil
}

func (r *audioRepository) FindAudioIDByText(ctx context.Context, lang5 domain.Lang5, text string) (domain.AudioID, error) {
	_, span := tracer.Start(ctx, "audioRepository.FindAudioIDByText")
	defer span.End()

	entity := audioEntity{}
	if result := r.db.Where("lang5 = ? and text = ?", lang5.String(), text).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, service.ErrAudioNotFound
		}
		return 0, result.Error
	}
	model, err := entity.toAudioModel()
	if err != nil {
		return 0, liberrors.Errorf("failed to toAudio. entity: %v, err: %w", entity, err)
	}
	return domain.AudioID(model.GetID()), nil
}
