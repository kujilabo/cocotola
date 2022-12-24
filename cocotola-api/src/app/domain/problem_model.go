//go:generate mockery --output mock --name ProblemModel
package domain

import (
	"context"
	"errors"

	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type ProblemID uint

type ProblemModel interface {
	userD.Model
	GetNumber() int
	GetProblemType() ProblemTypeName
	GetProperties(ctx context.Context) map[string]interface{}
}

type problemModel struct {
	userD.Model
	Number      int                    `validate:"required"`
	ProblemType ProblemTypeName        `validate:"required"`
	Properties  map[string]interface{} `validate:"required"`
}

func NewProblemModel(model userD.Model, number int, problemType ProblemTypeName, properties map[string]interface{}) (ProblemModel, error) {
	m := &problemModel{
		Model:       model,
		Number:      number,
		ProblemType: problemType,
		Properties:  properties,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *problemModel) GetNumber() int {
	return m.Number
}

func (m *problemModel) GetProblemType() ProblemTypeName {
	return m.ProblemType
}

func (m *problemModel) GetProperties(ctx context.Context) map[string]interface{} {
	panic(errors.New("not implemented"))
}
