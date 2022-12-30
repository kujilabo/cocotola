package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type WorkbookSearchCondition interface {
	GetPageNo() int
	GetPageSize() int
	GetSpaceIDs() []userD.SpaceID
}

type workbookSearchCondition struct {
	PageNo   int
	PageSize int
	SpaceIDs []userD.SpaceID
}

func NewWorkbookSearchCondition(pageNo, pageSize int, spaceIDs []userD.SpaceID) (WorkbookSearchCondition, error) {
	m := &workbookSearchCondition{
		PageNo:   pageNo,
		PageSize: pageSize,
		SpaceIDs: spaceIDs,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *workbookSearchCondition) GetPageNo() int {
	return p.PageNo
}

func (p *workbookSearchCondition) GetPageSize() int {
	return p.PageSize
}

func (p *workbookSearchCondition) GetSpaceIDs() []userD.SpaceID {
	return p.SpaceIDs
}

type WorkbookSearchResult interface {
	GetTotalCount() int
	GetResults() []WorkbookModel
}

type workbookSearchResult struct {
	TotalCount int
	Results    []WorkbookModel
}

func NewWorkbookSearchResult(totalCount int, results []WorkbookModel) (WorkbookSearchResult, error) {
	m := &workbookSearchResult{
		TotalCount: totalCount,
		Results:    results,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
func (m *workbookSearchResult) GetTotalCount() int {
	return m.TotalCount
}

func (m *workbookSearchResult) GetResults() []WorkbookModel {
	return m.Results
}

type WorkbookAddParameter interface {
	GetProblemType() ProblemTypeName
	GetName() string
	GetLang2() Lang2
	GetQuestionText() string
	GetProperties() map[string]string
}

type workbookAddParameter struct {
	ProblemType  ProblemTypeName
	Name         string
	Lang2        Lang2
	QuestionText string
	Properties   map[string]string
}

func NewWorkbookAddParameter(problemType ProblemTypeName, name string, lang2 Lang2, questionText string, properties map[string]string) (WorkbookAddParameter, error) {
	m := &workbookAddParameter{
		ProblemType:  problemType,
		Name:         name,
		Lang2:        lang2,
		QuestionText: questionText,
		Properties:   properties,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *workbookAddParameter) GetProblemType() ProblemTypeName {
	return p.ProblemType
}

func (p *workbookAddParameter) GetName() string {
	return p.Name
}

func (p *workbookAddParameter) GetLang2() Lang2 {
	return p.Lang2
}

func (p *workbookAddParameter) GetQuestionText() string {
	return p.QuestionText
}

func (p *workbookAddParameter) GetProperties() map[string]string {
	return p.Properties
}

type WorkbookUpdateParameter interface {
	GetName() string
	GetQuestionText() string
}

type workbookUpdateParameter struct {
	Name         string
	QuestionText string
}

func NewWorkbookUpdateParameter(name, questionText string) (WorkbookUpdateParameter, error) {
	m := &workbookUpdateParameter{
		Name:         name,
		QuestionText: questionText,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *workbookUpdateParameter) GetName() string {
	return p.Name
}

func (p *workbookUpdateParameter) GetQuestionText() string {
	return p.QuestionText
}
