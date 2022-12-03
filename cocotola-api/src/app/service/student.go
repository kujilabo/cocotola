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

	CheckQuota(ctx context.Context, problemType string, name QuotaName) error

	FindRecordbook(ctx context.Context, workbookID domain.WorkbookID, studyType string) (Recordbook, error)

	FindRecordbookSummary(ctx context.Context, workbookID domain.WorkbookID) (RecordbookSummary, error)

	GetStat(ctx context.Context) (Stat, error)
}

type student struct {
	domain.StudentModel
	rf        RepositoryFactory
	pf        ProcessorFactory
	spaceRepo userS.SpaceRepository
	userRf    userS.RepositoryFactory
}

func NewStudent(pf ProcessorFactory, rf RepositoryFactory, userRf userS.RepositoryFactory, studentModel domain.StudentModel) (Student, error) {
	if pf == nil {
		return nil, liberrors.Errorf("pf is nil. err: %w", libD.ErrInvalidArgument)
	}
	if studentModel == nil {
		return nil, errors.New("studentModel is nil")
	}
	spaceRepo, err := userRf.NewSpaceRepository()
	if err != nil {
		return nil, err
	}

	m := &student{
		StudentModel: studentModel,
		pf:           pf,
		rf:           rf,
		spaceRepo:    spaceRepo,
		userRf:       userRf,
	}

	return m, libD.Validator.Struct(m)
}

func (s *student) GetDefaultSpace(ctx context.Context) (userS.Space, error) {
	return s.spaceRepo.FindDefaultSpace(ctx, s)
}

func (s *student) GetPersonalSpace(ctx context.Context) (userS.Space, error) {
	return s.spaceRepo.FindPersonalSpace(ctx, s)
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

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, liberrors.Errorf("NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindPersonalWorkbooks(ctx, s, newCondition)
}

func (s *student) FindWorkbookByID(ctx context.Context, id domain.WorkbookID) (Workbook, error) {
	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, liberrors.Errorf("s.rf.NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindWorkbookByID(ctx, s, id)
}

func (s *student) FindWorkbookByName(ctx context.Context, name string) (Workbook, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return nil, liberrors.Errorf("s.GetPersonalSpace. err: %w", err)
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, liberrors.Errorf("s.rf.NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindWorkbookByName(ctx, s, userD.SpaceID(space.GetID()), name)
}

func (s *student) AddWorkbookToPersonalSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return 0, liberrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

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

	return workbook.UpdateWorkbook(ctx, s, version, parameter)
}

func (s *student) RemoveWorkbook(ctx context.Context, workbookID domain.WorkbookID, version int) error {
	workbook, err := s.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return liberrors.Errorf("s.FindWorkbookByID. err: %w", err)
	}

	return workbook.RemoveWorkbook(ctx, s, version)
}

func (s *student) CheckQuota(ctx context.Context, problemType string, name QuotaName) error {
	processor, err := s.pf.NewProblemQuotaProcessor(problemType)
	if err != nil {
		return liberrors.Errorf("s.pf.NewProblemQuotaProcessor. err: %w", err)
	}

	userQuotaRepo, err := s.rf.NewUserQuotaRepository(ctx)
	if err != nil {
		return liberrors.Errorf("s.rf.NewUserQuotaRepository. err: %w", err)
	}

	switch name {
	case QuotaNameSize:
		unit := processor.GetUnitForSizeQuota()
		limit := processor.GetLimitForSizeQuota()
		isExceeded, err := userQuotaRepo.IsExceeded(ctx, s.GetOrganizationID(), userD.AppUserID(s.GetID()), problemType+"_size", unit, limit)
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
		isExceeded, err := userQuotaRepo.IsExceeded(ctx, s.GetOrganizationID(), userD.AppUserID(s.GetID()), problemType+"_update", unit, limit)
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

func (s *student) FindRecordbook(ctx context.Context, workbookID domain.WorkbookID, studyType string) (Recordbook, error) {
	return NewRecordbook(s.rf, s, workbookID, studyType)
}

func (s *student) FindRecordbookSummary(ctx context.Context, workbookID domain.WorkbookID) (RecordbookSummary, error) {
	return NewRecordbookSummary(s.rf, s, workbookID)
}

func (s *student) GetStat(ctx context.Context) (Stat, error) {
	statRepo, err := s.rf.NewStatRepository(ctx)
	if err != nil {
		return nil, err
	}

	stat, err := statRepo.FindStat(ctx, userD.AppUserID(s.GetID()))
	if err != nil {
		return nil, err
	}

	return stat, nil
}
