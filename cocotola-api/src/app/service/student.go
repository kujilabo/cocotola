//go:generate mockery --output mock --name Student
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type Student interface {
	domain.StudentModel

	GetDefaultSpace(ctx context.Context) (userS.Space, error)
	GetPersonalSpace(ctx context.Context) (userS.Space, error)

	FindWorkbooksFromPersonalSpace(ctx context.Context, condition WorkbookSearchCondition) (WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, id domain.WorkbookID) (Workbook, error)

	FindWorkbookByName(ctx context.Context, name string) (Workbook, error)

	AddWorkbookToPersonalSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, workbookID domain.WorkbookID, version int, parameter WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, id domain.WorkbookID, version int) error

	CheckQuota(ctx context.Context, problemType domain.ProblemTypeName, name QuotaName) error

	FindRecordbook(ctx context.Context, workbookID domain.WorkbookID, studyType domain.StudyTypeName) (Recordbook, error)

	FindRecordbookSummary(ctx context.Context, workbookID domain.WorkbookID) (RecordbookSummary, error)

	GetStat(ctx context.Context) (Stat, error)

	FindPreferences(ctx context.Context) (userS.UserPreferences, error)
}

type student struct {
	domain.StudentModel
	rf        RepositoryFactory
	pf        ProcessorFactory
	spaceRepo userS.SpaceRepository
}

func NewStudent(ctx context.Context, pf ProcessorFactory, rf RepositoryFactory, studentModel domain.StudentModel) (Student, error) {
	if pf == nil {
		return nil, liberrors.Errorf("pf is nil. err: %w", libD.ErrInvalidArgument)
	}

	if studentModel == nil {
		return nil, errors.New("studentModel is nil")
	}

	userRf, err := rf.NewUserRepositoryFactory(ctx)
	if err != nil {
		return nil, liberrors.Errorf("rf.NewUserRepositoryFactory. err: %w", libD.ErrInvalidArgument)
	}

	spaceRepo := userRf.NewSpaceRepository(ctx)
	m := &student{
		StudentModel: studentModel,
		pf:           pf,
		rf:           rf,
		spaceRepo:    spaceRepo,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (s *student) GetDefaultSpace(ctx context.Context) (userS.Space, error) {
	space, err := s.spaceRepo.FindDefaultSpace(ctx, s)
	if err != nil {
		return nil, liberrors.Errorf("s.spaceRepo.FindDefaultSpace. err: %w", err)
	}

	return space, nil
}

func (s *student) GetPersonalSpace(ctx context.Context) (userS.Space, error) {
	space, err := s.spaceRepo.FindPersonalSpace(ctx, s)
	if err != nil {
		return nil, liberrors.Errorf("s.spaceRepo.FindPersonalSpace. err: %w", err)
	}

	return space, nil
}

func (s *student) FindWorkbooksFromPersonalSpace(ctx context.Context, condition WorkbookSearchCondition) (WorkbookSearchResult, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return nil, liberrors.Errorf("GetPersonalSpace. err: %w", err)
	}

	// specify space
	newCondition, err := NewWorkbookSearchCondition(condition.GetPageNo(), condition.GetPageSize(), []userD.SpaceID{userD.SpaceID(space.GetID())})
	if err != nil {
		return nil, liberrors.Errorf("NewWorkbookSearchCondition. err: %w", err)
	}

	workbookRepo := s.rf.NewWorkbookRepository(ctx)
	workbooks, err := workbookRepo.FindPersonalWorkbooks(ctx, s, newCondition)
	if err != nil {
		return nil, liberrors.Errorf("workbookRepo.FindPersonalWorkbooks. err: %w", err)
	}

	return workbooks, nil
}

func (s *student) FindWorkbookByID(ctx context.Context, id domain.WorkbookID) (Workbook, error) {
	workbookRepo := s.rf.NewWorkbookRepository(ctx)
	workbook, err := workbookRepo.FindWorkbookByID(ctx, s, id)
	if err != nil {
		return nil, liberrors.Errorf("workbookRepo.FindWorkbookByID. err: %w", err)
	}

	return workbook, nil
}

func (s *student) FindWorkbookByName(ctx context.Context, name string) (Workbook, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return nil, liberrors.Errorf("s.GetPersonalSpace. err: %w", err)
	}

	workbookRepo := s.rf.NewWorkbookRepository(ctx)
	workbook, err := workbookRepo.FindWorkbookByName(ctx, s, userD.SpaceID(space.GetID()), name)
	if err != nil {
		return nil, liberrors.Errorf("workbookRepo.FindWorkbookByName. err: %w", err)
	}

	return workbook, nil
}

func (s *student) AddWorkbookToPersonalSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return 0, liberrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	}

	workbookRepo := s.rf.NewWorkbookRepository(ctx)
	workbookID, err := workbookRepo.AddWorkbook(ctx, s, userD.SpaceID(space.GetID()), parameter)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddWorkbook. err: %w", err)
	}

	return workbookID, nil
}

func (s *student) UpdateWorkbook(ctx context.Context, workbookID domain.WorkbookID, version int, parameter WorkbookUpdateParameter) error {
	workbook, err := s.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return liberrors.Errorf("s.FindWorkbookByID. err: %w", err)
	}

	if err := workbook.UpdateWorkbook(ctx, s, version, parameter); err != nil {
		return liberrors.Errorf("workbook.UpdateWorkbook. err: %w", err)
	}

	return nil
}

func (s *student) RemoveWorkbook(ctx context.Context, workbookID domain.WorkbookID, version int) error {
	workbook, err := s.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return liberrors.Errorf("s.FindWorkbookByID. err: %w", err)
	}

	if err := workbook.RemoveWorkbook(ctx, s, version); err != nil {
		return liberrors.Errorf("workbook.RemoveWorkbook. err: %w", err)
	}

	return nil
}

func (s *student) CheckQuota(ctx context.Context, problemType domain.ProblemTypeName, name QuotaName) error {
	processor, err := s.pf.NewProblemQuotaProcessor(problemType)
	if err != nil {
		return liberrors.Errorf("s.pf.NewProblemQuotaProcessor. err: %w", err)
	}

	userQuotaRepo := s.rf.NewUserQuotaRepository(ctx)

	problemTypeName := (string)(problemType)
	switch name {
	case QuotaNameSize:
		unit := processor.GetUnitForSizeQuota()
		limit := processor.GetLimitForSizeQuota()
		isExceeded, err := userQuotaRepo.IsExceeded(ctx, s.GetOrganizationID(), s.GetAppUserID(), problemTypeName+"_size", unit, limit)
		if err != nil {
			return liberrors.Errorf("userQuotaRepo.IsExceeded(size). err: %w", err)
		}

		if isExceeded {
			return ErrQuotaExceeded
		}

		return nil
	case QuotaNameUpdate:
		unit := processor.GetUnitForUpdateQuota()
		limit := processor.GetLimitForUpdateQuota()
		isExceeded, err := userQuotaRepo.IsExceeded(ctx, s.GetOrganizationID(), s.GetAppUserID(), problemTypeName+"_update", unit, limit)
		if err != nil {
			return liberrors.Errorf("userQuotaRepo.IsExceeded(update). err: %w", err)
		}

		if isExceeded {
			return ErrQuotaExceeded
		}

		return nil
	default:
		return liberrors.Errorf("invalid name. name: %s", name)
	}
}

func (s *student) FindRecordbook(ctx context.Context, workbookID domain.WorkbookID, studyType domain.StudyTypeName) (Recordbook, error) {
	return NewRecordbook(s.rf, s, workbookID, studyType)
}

func (s *student) FindRecordbookSummary(ctx context.Context, workbookID domain.WorkbookID) (RecordbookSummary, error) {
	recordbookSummary, err := NewRecordbookSummary(s.rf, s, workbookID)
	if err != nil {
		return nil, liberrors.Errorf("NewRecordbookSummary. err: %w", err)
	}

	return recordbookSummary, nil
}

func (s *student) GetStat(ctx context.Context) (Stat, error) {
	statRepo := s.rf.NewStatRepository(ctx)
	stat, err := statRepo.FindStat(ctx, s.GetAppUserID())
	if err != nil {
		return nil, liberrors.Errorf("statRepo.FindStat. err: %w", err)
	}

	return stat, nil
}

func (s *student) FindPreferences(ctx context.Context) (userS.UserPreferences, error) {
	return userS.NewUserPreferences, nil
}
