//go:generate mockery --output mock --name UserUsecase
package usecase

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/domain"
	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type UserUsecase interface {
	Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (service.Audio, error)

	FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error)
}

type userUsecase struct {
	transaction       service.Transaction
	synthesizerClient service.SynthesizerClient
}

func NewUserUsecase(transaction service.Transaction, synthesizerClient service.SynthesizerClient) UserUsecase {
	return &userUsecase{
		transaction:       transaction,
		synthesizerClient: synthesizerClient,
	}
}

func (u *userUsecase) Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (service.Audio, error) {
	var audio service.Audio
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		// try to find the audio content from the DB
		repo := rf.NewAudioRepository(ctx)
		if tmpAudioContentFromDB, err := repo.FindByLangAndText(ctx, lang2.ToLang5(), text); err == nil {
			audio = tmpAudioContentFromDB
			return nil
		} else if !errors.Is(err, service.ErrAudioNotFound) {
			return liberrors.Errorf("FindByLangAndText. err: %w", err)
		}

		// synthesize text via the Web API
		audioContent, err := u.synthesizerClient.Synthesize(ctx, lang2.ToLang5(), text)
		if err != nil {
			return liberrors.Errorf("to u.synthesizerClient.Synthesize. err: %w", err)
		}

		audioID, err := repo.AddAudio(ctx, lang2.ToLang5(), text, audioContent)
		if err != nil {
			return liberrors.Errorf("toAddAudio. err: %w", err)
		}

		audioModel, err := domain.NewAudioModel(uint(audioID), lang2.ToLang5(), text, audioContent)
		if err != nil {
			return liberrors.Errorf("NewAudioModel. err: %w", err)
		}

		tmpAudio, err := service.NewAudio(audioModel)
		if err != nil {
			return err
		}
		audio = tmpAudio
		return nil
	}); err != nil {
		return nil, err
	}

	return audio, nil
}

func (u *userUsecase) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error) {
	var audio service.Audio
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		repo := rf.NewAudioRepository(ctx)
		tmpAudio, err := repo.FindAudioByAudioID(ctx, audioID)
		if err != nil {
			return err
		}

		audio = tmpAudio
		return nil
	}); err != nil {
		return nil, err
	}

	return audio, nil
}
