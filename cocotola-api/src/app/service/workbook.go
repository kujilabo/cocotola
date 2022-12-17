//go:generate mockery --output mock --name Workbook
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type Workbook interface {
	domain.WorkbookModel

	// FindProblems searches for problems based on search condition
	FindProblems(ctx context.Context, operator domain.StudentModel, param ProblemSearchCondition) (ProblemSearchResult, error)

	FindAllProblems(ctx context.Context, operator domain.StudentModel) (ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param ProblemIDsCondition) (ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error)

	// FindProblems searches for problem based on a problem ID
	FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (Problem, error)

	AddProblem(ctx context.Context, operator domain.StudentModel, param ProblemAddParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateProblemProperty(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	RemoveProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator domain.StudentModel, version int) error

	CountProblems(ctx context.Context, operator domain.StudentModel) (int, error)
}

type workbook struct {
	domain.WorkbookModel
	rf           RepositoryFactory
	pf           ProcessorFactory
	workbookRepo WorkbookRepository
	problemRepo  ProblemRepository
}

func NewWorkbook(ctx context.Context, rf RepositoryFactory, pf ProcessorFactory, workbookModel domain.WorkbookModel) (Workbook, error) {
	workbookRepo, err := rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, err
	}

	problemRepo, err := rf.NewProblemRepository(ctx, workbookModel.GetProblemType())
	if err != nil {
		return nil, err
	}

	m := &workbook{
		WorkbookModel: workbookModel,
		rf:            rf,
		pf:            pf,
		workbookRepo:  workbookRepo,
		problemRepo:   problemRepo,
	}

	return m, libD.Validator.Struct(m)
}

func (m *workbook) GetWorkbookModel() domain.WorkbookModel {
	return m.WorkbookModel
}

func (m *workbook) FindProblems(ctx context.Context, operator domain.StudentModel, param ProblemSearchCondition) (ProblemSearchResult, error) {
	return m.problemRepo.FindProblems(ctx, operator, param)
}

func (m *workbook) FindAllProblems(ctx context.Context, operator domain.StudentModel) (ProblemSearchResult, error) {
	return m.problemRepo.FindAllProblems(ctx, operator, m.GetWorkbookModel().GetWorkbookID())
}

func (m *workbook) FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param ProblemIDsCondition) (ProblemSearchResult, error) {
	return m.problemRepo.FindProblemsByProblemIDs(ctx, operator, param)
}

func (m *workbook) FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error) {
	return m.problemRepo.FindProblemIDs(ctx, operator, m.GetWorkbookModel().GetWorkbookID())
}

func (m *workbook) FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (Problem, error) {
	id, err := NewProblemSelectParameter1(m.GetWorkbookModel().GetWorkbookID(), problemID)
	if err != nil {
		return nil, err
	}
	return m.problemRepo.FindProblemByID(ctx, operator, id)
}

func (m *workbook) AddProblem(ctx context.Context, operator domain.StudentModel, param ProblemAddParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.AddProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemAddProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	logger.Infof("processor.AddProblem")
	return processor.AddProblem(ctx, m.rf, operator, m.GetWorkbookModel(), param)
}

func (m *workbook) UpdateProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.UpdateProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemUpdateProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	return processor.UpdateProblem(ctx, m.rf, operator, m.GetWorkbookModel(), id, param)
}

func (m *workbook) UpdateProblemProperty(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2, param ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.UpdateProblemProperty")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemUpdateProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	return processor.UpdateProblemProperty(ctx, m.rf, operator, m.GetWorkbookModel(), id, param)
}

func (m *workbook) RemoveProblem(ctx context.Context, operator domain.StudentModel, id ProblemSelectParameter2) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.RemoveProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemRemoveProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	return processor.RemoveProblem(ctx, m.rf, operator, id)
}

func (m *workbook) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter WorkbookUpdateParameter) error {
	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return ErrWorkbookPermissionDenied
	}

	return m.workbookRepo.UpdateWorkbook(ctx, operator, m.GetWorkbookModel().GetWorkbookID(), version, parameter)
}

func (m *workbook) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, version int) error {
	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeRemove) {
		return ErrWorkbookPermissionDenied
	}

	return m.workbookRepo.RemoveWorkbook(ctx, operator, m.GetWorkbookModel().GetWorkbookID(), version)
}

func (m *workbook) CountProblems(ctx context.Context, operator domain.StudentModel) (int, error) {
	return m.problemRepo.CountProblems(ctx, operator, m.GetWorkbookModel().GetWorkbookID())
}
