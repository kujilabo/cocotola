//go:generate mockery --output mock --name ProblemAddParameter
//go:generate mockery --output mock --name ProblemSelectParameter1
//go:generate mockery --output mock --name ProblemSelectParameter2
//go:generate mockery --output mock --name ProblemUpdateParameter
//go:generate mockery --output mock --name ProblemSearchCondition
package domain

import (
	"errors"
	"strconv"

	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type ProblemAddParameter interface {
	GetWorkbookID() WorkbookID
	GetProperties() map[string]string
	GetStringProperty(name string) (string, error)
	GetIntProperty(name string) (int, error)
}

type problemAddParameter struct {
	WorkbookID WorkbookID `validate:"required"`
	Properties map[string]string
}

func NewProblemAddParameter(workbookID WorkbookID, properties map[string]string) (ProblemAddParameter, error) {
	m := &problemAddParameter{
		WorkbookID: workbookID,
		Properties: properties,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *problemAddParameter) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}

func (p *problemAddParameter) GetProperties() map[string]string {
	return p.Properties
}

func (p *problemAddParameter) GetStringProperty(name string) (string, error) {
	s, ok := p.Properties[name]
	if !ok {
		return "", errors.New("key not found")
	}
	return s, nil
}
func (p *problemAddParameter) GetIntProperty(name string) (int, error) {
	i, err := strconv.Atoi(p.Properties[name])
	if err != nil {
		return 0, liberrors.Errorf("%q property is not a integer. value: %s, err: %w", name, p.Properties[name], err)
	}
	return i, nil
}

type ProblemSelectParameter1 interface {
	GetWorkbookID() WorkbookID
	GetProblemID() ProblemID
}

type problemSelectParameter1 struct {
	WorkbookID WorkbookID
	ProblemID  ProblemID
}

func NewProblemSelectParameter1(WorkbookID WorkbookID, problemID ProblemID) (ProblemSelectParameter1, error) {
	m := &problemSelectParameter1{
		WorkbookID: WorkbookID,
		ProblemID:  problemID,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *problemSelectParameter1) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}
func (p *problemSelectParameter1) GetProblemID() ProblemID {
	return p.ProblemID
}

type ProblemSelectParameter2 interface {
	GetWorkbookID() WorkbookID
	GetProblemID() ProblemID
	GetVersion() int
}

type problemSelectParameter2 struct {
	WorkbookID WorkbookID
	ProblemID  ProblemID
	Version    int
}

func NewProblemSelectParameter2(WorkbookID WorkbookID, problemID ProblemID, version int) (ProblemSelectParameter2, error) {
	m := &problemSelectParameter2{
		WorkbookID: WorkbookID,
		ProblemID:  problemID,
		Version:    version,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *problemSelectParameter2) GetWorkbookID() WorkbookID {
	return p.WorkbookID
}
func (p *problemSelectParameter2) GetProblemID() ProblemID {
	return p.ProblemID
}

func (p *problemSelectParameter2) GetVersion() int {
	return p.Version
}

type ProblemUpdateParameter interface {
	GetProperties() map[string]string
	GetStringProperty(name string) (string, error)
	GetIntProperty(name string) (int, error)
}

type problemUpdateParameter struct {
	Properties map[string]string
}

func NewProblemUpdateParameter(properties map[string]string) (ProblemUpdateParameter, error) {
	m := &problemUpdateParameter{
		Properties: properties,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *problemUpdateParameter) GetProperties() map[string]string {
	return p.Properties
}
func (p *problemUpdateParameter) GetStringProperty(name string) (string, error) {
	s, ok := p.Properties[name]
	if !ok {
		return "", errors.New("key not found")
	}
	return s, nil
}
func (p *problemUpdateParameter) GetIntProperty(name string) (int, error) {
	i, err := strconv.Atoi(p.Properties[name])
	if err != nil {
		return 0, liberrors.Errorf("strconv.Atoi. err: %w", err)
	}
	return i, nil
}

type ProblemPropertyUpdateParameter interface {
	GetKey() string
	GetValue() string
}

type problemPropertyUpdateParameter struct {
	Key   string
	Value string
}

func NewProblemPropertyUpdateParameter(key, value string) (ProblemPropertyUpdateParameter, error) {
	m := &problemPropertyUpdateParameter{
		Key:   key,
		Value: value,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *problemPropertyUpdateParameter) GetKey() string {
	return p.Key
}

func (p *problemPropertyUpdateParameter) GetValue() string {
	return p.Value
}

type ProblemSearchCondition interface {
	GetWorkbookID() WorkbookID
	GetPageNo() int
	GetPageSize() int
	GetKeyword() string
}

type problemSearchCondition struct {
	WorkbookID WorkbookID
	PageNo     int `validate:"required,gte=1"`
	PageSize   int `validate:"required,gte=1,lte=1000"`
	Keyword    string
}

func NewProblemSearchCondition(workbookID WorkbookID, pageNo, pageSize int, keyword string) (ProblemSearchCondition, error) {
	m := &problemSearchCondition{
		WorkbookID: workbookID,
		PageNo:     pageNo,
		PageSize:   pageSize,
		Keyword:    keyword,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (c *problemSearchCondition) GetWorkbookID() WorkbookID {
	return c.WorkbookID
}

func (c *problemSearchCondition) GetPageNo() int {
	return c.PageNo
}

func (c *problemSearchCondition) GetPageSize() int {
	return c.PageSize
}

func (c *problemSearchCondition) GetKeyword() string {
	return c.Keyword
}

type ProblemIDsCondition interface {
	GetWorkbookID() WorkbookID
	GetIDs() []ProblemID
}

type problemIDsCondition struct {
	WorkbookID WorkbookID
	IDs        []ProblemID
}

func NewProblemIDsCondition(workbookID WorkbookID, ids []ProblemID) (ProblemIDsCondition, error) {
	m := &problemIDsCondition{
		WorkbookID: workbookID,
		IDs:        ids,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (c *problemIDsCondition) GetWorkbookID() WorkbookID {
	return c.WorkbookID
}

func (c *problemIDsCondition) GetIDs() []ProblemID {
	return c.IDs
}

type ProblemSearchResult interface {
	GetTotalCount() int
	GetResults() []ProblemModel
}

type problemSearchResult struct {
	TotalCount int
	Results    []ProblemModel
}

func NewProblemSearchResult(totalCount int, results []ProblemModel) (ProblemSearchResult, error) {
	m := &problemSearchResult{
		TotalCount: totalCount,
		Results:    results,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *problemSearchResult) GetTotalCount() int {
	return m.TotalCount
}

func (m *problemSearchResult) GetResults() []ProblemModel {
	return m.Results
}
