package converter

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func ToWorkbookSearchResponse(result domain.WorkbookSearchResult) (*entity.WorkbookSearchResponse, error) {
	workbooks := make([]*entity.WorkbookResponseHTTPEntity, len(result.GetResults()))
	for i, w := range result.GetResults() {
		model, err := entity.NewModel(w)
		if err != nil {
			return nil, liberrors.Errorf("entity.NewModel. err: %w", err)
		}

		workbooks[i] = &entity.WorkbookResponseHTTPEntity{
			Model:        model,
			Name:         w.GetName(),
			Lang2:        w.GetLang2().String(),
			ProblemType:  string(w.GetProblemType()),
			QuestionText: w.GetQuestionText(),
		}
	}

	e := &entity.WorkbookSearchResponse{
		TotalCount: result.GetTotalCount(),
		Results:    workbooks,
	}

	if err := libD.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToWorkbookHTTPEntity(workbook domain.WorkbookModel) (entity.WorkbookResponseHTTPEntity, error) {
	e := entity.WorkbookResponseHTTPEntity{
		Model: entity.Model{
			ID:        workbook.GetID(),
			Version:   workbook.GetVersion(),
			CreatedBy: workbook.GetCreatedBy(),
			UpdatedBy: workbook.GetUpdatedBy(),
		},
		Name:         workbook.GetName(),
		Lang2:        workbook.GetLang2().String(),
		ProblemType:  string(workbook.GetProblemType()),
		QuestionText: workbook.GetQuestionText(),
	}

	if err := libD.Validator.Struct(e); err != nil {
		return entity.WorkbookResponseHTTPEntity{}, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToWorkbookAddParameter(param *entity.WorkbookAddParameter) (domain.WorkbookAddParameter, error) {
	domainParam, err := domain.NewWorkbookAddParameter(domain.ProblemTypeName(param.ProblemType), param.Name, domain.Lang2JA, param.QuestionText, map[string]string{
		"audioEnabled": "false",
	})

	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return domainParam, nil
}

func ToWorkbookUpdateParameter(param *entity.WorkbookUpdateParameter) (domain.WorkbookUpdateParameter, error) {
	domainParam, err := domain.NewWorkbookUpdateParameter(param.Name, param.QuestionText)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return domainParam, nil
}
