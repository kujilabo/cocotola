package converter

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/controller/entity"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func ToTatoebaSentenceSearchCondition(ctx context.Context, param *entity.TatoebaSentenceFindParameter) (service.TatoebaSentenceSearchCondition, error) {
	condition, err := service.NewTatoebaSentenceSearchCondition(param.PageNo, param.PageSize, param.Keyword, param.Random)
	if err != nil {
		return nil, liberrors.Errorf("service.NewTatoebaSentenceSearchCondition. err: %w", err)
	}

	return condition, nil
}

func ToTatoebaSentenceResponse(ctx context.Context, result *service.TatoebaSentencePairSearchResult) (*entity.TatoebaSentenceFindResponse, error) {
	entities := make([]entity.TatoebaSentencePair, len(result.Results))
	for i, m := range result.Results {
		src := entity.TatoebaSentence{
			SentenceNumber: m.GetSrc().GetSentenceNumber(),
			Lang2:          m.GetSrc().GetLang2().String(),
			Text:           m.GetSrc().GetText(),
			Author:         m.GetSrc().GetAuthor(),
			UpdatedAt:      m.GetSrc().GetUpdatedAt(),
		}
		dst := entity.TatoebaSentence{
			SentenceNumber: m.GetDst().GetSentenceNumber(),
			Lang2:          m.GetDst().GetLang2().String(),
			Text:           m.GetDst().GetText(),
			Author:         m.GetDst().GetAuthor(),
			UpdatedAt:      m.GetDst().GetUpdatedAt(),
		}
		entities[i] = entity.TatoebaSentencePair{
			Src: src,
			Dst: dst,
		}
	}

	return &entity.TatoebaSentenceFindResponse{
		TotalCount: result.TotalCount,
		Results:    entities,
	}, nil
}
