package converter

import (
	"context"
	"encoding/json"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

func ToProblemSearchCondition(ctx context.Context, param *entity.ProblemFindParameter, workbookID domain.WorkbookID) (domain.ProblemSearchCondition, error) {
	condition, err := domain.NewProblemSearchCondition(workbookID, param.PageNo, param.PageSize, param.Keyword)
	if err != nil {
		return nil, liberrors.Errorf("new Model. err: %w", err)
	}

	return condition, nil
}

func ToProblemFindResponse(ctx context.Context, result domain.ProblemSearchResult) (*entity.ProblemFindResponse, error) {
	problems := make([]*entity.Problem, len(result.GetResults()))
	for i, p := range result.GetResults() {
		properties := p.GetProperties(ctx)
		bytes, err := json.Marshal(properties)
		if err != nil {
			return nil, liberrors.Errorf("convert properties to json. err: %w, properties: %+v", err, properties)
		}

		model, err := entity.NewModel(p)
		if err != nil {
			return nil, liberrors.Errorf("new Model. err: %w", err)
		}

		problems[i] = &entity.Problem{
			Model:       model,
			Number:      p.GetNumber(),
			ProblemType: string(p.GetProblemType()),
			Properties:  bytes,
		}
	}

	e := &entity.ProblemFindResponse{
		TotalCount: result.GetTotalCount(),
		Results:    problems,
	}
	if err := libD.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToProblemFindAllResponse(ctx context.Context, result domain.ProblemSearchResult) (*entity.ProblemFindAllResponse, error) {
	problems := make([]*entity.SimpleProblem, len(result.GetResults()))
	for i, p := range result.GetResults() {
		bytes, err := json.Marshal(p.GetProperties(ctx))
		if err != nil {
			return nil, liberrors.Errorf("json.Marshal. err: %w", err)
		}

		model, err := entity.NewModel(p)
		if err != nil {
			return nil, liberrors.Errorf("entity.NewMode. err: %w", err)
		}

		problems[i] = &entity.SimpleProblem{
			ID:          model.ID,
			Version:     model.Version,
			Number:      p.GetNumber(),
			ProblemType: string(p.GetProblemType()),
			Properties:  bytes,
		}
	}

	e := &entity.ProblemFindAllResponse{
		TotalCount: result.GetTotalCount(),
		Results:    problems,
	}
	if err := libD.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToProblemResponse(ctx context.Context, problem domain.ProblemModel) (*entity.Problem, error) {
	logger := log.FromContext(ctx)
	// FIXME
	logger.Infof("------properties: %+v", problem.GetProperties(ctx))

	bytes, err := json.Marshal(problem.GetProperties(ctx))
	if err != nil {
		return nil, liberrors.Errorf("json.Marshal. err: %w", err)
	}

	model, err := entity.NewModel(problem)
	if err != nil {
		return nil, liberrors.Errorf("entity.NewModel. err: %w", err)
	}

	e := &entity.Problem{
		Model:       model,
		Number:      problem.GetNumber(),
		ProblemType: string(problem.GetProblemType()),
		Properties:  bytes,
	}

	if err := libD.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToProblemIDsCondition(ctx context.Context, param *entity.ProblemIDsParameter, workbookID domain.WorkbookID) (domain.ProblemIDsCondition, error) {
	ids := make([]domain.ProblemID, 0)
	for _, id := range param.IDs {
		ids = append(ids, domain.ProblemID(id))
	}
	domainParam, err := domain.NewProblemIDsCondition(workbookID, ids)
	if err != nil {
		return nil, liberrors.Errorf("service.NewProblemIDsCondition. err: %w", err)
	}

	return domainParam, nil

}

func ToProblemAddParameter(workbookID domain.WorkbookID, param *entity.ProblemAddParameter) (domain.ProblemAddParameter, error) {
	var properties map[string]string
	if err := json.Unmarshal(param.Properties, &properties); err != nil {
		return nil, liberrors.Errorf("Unmarshal. err: %w", err)
	}

	domainParam, err := domain.NewProblemAddParameter(workbookID /*param.Number, */, properties)
	if err != nil {
		return nil, liberrors.Errorf("service.NewProblemAddParameter. err: %w", err)
	}

	return domainParam, nil
}

func ToProblemUpdateParameter(param *entity.ProblemUpdateParameter) (domain.ProblemUpdateParameter, error) {
	var properties map[string]string
	if err := json.Unmarshal(param.Properties, &properties); err != nil {
		return nil, liberrors.Errorf("Unmarshal. err: %w", err)
	}

	domainParam, err := domain.NewProblemUpdateParameter( /*param.Number, */ properties)
	if err != nil {
		return nil, liberrors.Errorf("service.NewProblemUpdateParameter. err: %w", err)
	}

	return domainParam, nil
}
