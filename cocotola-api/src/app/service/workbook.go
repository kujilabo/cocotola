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
	FindProblems(ctx context.Context, operator domain.StudentModel, param domain.ProblemSearchCondition) (domain.ProblemSearchResult, error)

	FindAllProblems(ctx context.Context, operator domain.StudentModel) (domain.ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param domain.ProblemIDsCondition) (domain.ProblemSearchResult, error)

	FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error)

	// FindProblems searches for problem based on a problem ID
	FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (Problem, error)

	AddProblem(ctx context.Context, operator domain.StudentModel, param domain.ProblemAddParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateProblem(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateProblemProperty(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	RemoveProblem(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter domain.WorkbookUpdateParameter) error

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
	workbookRepo := rf.NewWorkbookRepository(ctx)

	problemRepo, err := rf.NewProblemRepository(ctx, workbookModel.GetProblemType())
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	m := &workbook{
		WorkbookModel: workbookModel,
		rf:            rf,
		pf:            pf,
		workbookRepo:  workbookRepo,
		problemRepo:   problemRepo,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *workbook) GetWorkbookModel() domain.WorkbookModel {
	return m.WorkbookModel
}

func (m *workbook) FindProblems(ctx context.Context, operator domain.StudentModel, param domain.ProblemSearchCondition) (domain.ProblemSearchResult, error) {
	problems, err := m.problemRepo.FindProblems(ctx, operator, param)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return problems, nil
}

func (m *workbook) FindAllProblems(ctx context.Context, operator domain.StudentModel) (domain.ProblemSearchResult, error) {
	problems, err := m.problemRepo.FindAllProblems(ctx, operator, m.GetWorkbookModel().GetWorkbookID())
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return problems, nil
}

func (m *workbook) FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param domain.ProblemIDsCondition) (domain.ProblemSearchResult, error) {
	problems, err := m.problemRepo.FindProblemsByProblemIDs(ctx, operator, param)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return problems, nil
}

func (m *workbook) FindProblemIDs(ctx context.Context, operator domain.StudentModel) ([]domain.ProblemID, error) {
	ids, err := m.problemRepo.FindProblemIDs(ctx, operator, m.GetWorkbookModel().GetWorkbookID())
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return ids, nil
}

func (m *workbook) FindProblemByID(ctx context.Context, operator domain.StudentModel, problemID domain.ProblemID) (Problem, error) {
	id, err := domain.NewProblemSelectParameter1(m.GetWorkbookModel().GetWorkbookID(), problemID)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	problem, err := m.problemRepo.FindProblemByID(ctx, operator, id)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return problem, nil
}

func (m *workbook) AddProblem(ctx context.Context, operator domain.StudentModel, param domain.ProblemAddParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.AddProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemAddProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	added, updated, removed, err := processor.AddProblem(ctx, m.rf, operator, m.GetWorkbookModel(), param)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("c.userClient.DictionaryLookup. err: %w", err)
	}
	return added, updated, removed, nil
}

func (m *workbook) UpdateProblem(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.UpdateProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemUpdateProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	added, updated, removed, err := processor.UpdateProblem(ctx, m.rf, operator, m.GetWorkbookModel(), id, param)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf(". err: %w", err)
	}

	return added, updated, removed, nil
}

func (m *workbook) UpdateProblemProperty(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.UpdateProblemProperty")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemUpdateProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	added, updated, removed, err := processor.UpdateProblemProperty(ctx, m.rf, operator, m.GetWorkbookModel(), id, param)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf(". err: %w", err)
	}
	return added, updated, removed, nil
}

func (m *workbook) RemoveProblem(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbook.RemoveProblem")

	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return nil, nil, nil, errors.New("no update privilege")
	}

	processor, err := m.pf.NewProblemRemoveProcessor(m.GetWorkbookModel().GetProblemType())
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor not found. problemType: %s, err: %w", m.GetWorkbookModel().GetProblemType(), err)
	}

	added, updated, removed, err := processor.RemoveProblem(ctx, m.rf, operator, id)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("processor.RemoveProblem. err: %w", err)
	}

	return added, updated, removed, nil
}

func (m *workbook) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, version int, parameter domain.WorkbookUpdateParameter) error {
	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeUpdate) {
		return ErrWorkbookPermissionDenied
	}

	if err := m.workbookRepo.UpdateWorkbook(ctx, operator, m.GetWorkbookModel().GetWorkbookID(), version, parameter); err != nil {
		return liberrors.Errorf("m.workbookRepo.UpdateWorkbook. err: %w", err)
	}

	return nil
}

func (m *workbook) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, version int) error {
	if !m.GetWorkbookModel().HasPrivilege(domain.PrivilegeRemove) {
		return ErrWorkbookPermissionDenied
	}

	if err := m.workbookRepo.RemoveWorkbook(ctx, operator, m.GetWorkbookModel().GetWorkbookID(), version); err != nil {
		return liberrors.Errorf("m.workbookRepo.RemoveWorkbook. err: %w", err)
	}

	return nil
}

func (m *workbook) CountProblems(ctx context.Context, operator domain.StudentModel) (int, error) {
	num, err := m.problemRepo.CountProblems(ctx, operator, m.GetWorkbookModel().GetWorkbookID())
	if err != nil {
		return 0, liberrors.Errorf("m.problemRepo.CountProblems. err: %w", err)
	}

	return num, nil
}
