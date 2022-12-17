package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	jobG "github.com/kujilabo/cocotola/cocotola-api/src/job/gateway"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type repositoryFactory struct {
	db                  *gorm.DB
	driverName          string
	userRff             userG.RepositoryFactoryFunc
	pf                  service.ProcessorFactory
	problemRepositories map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error)
	problemTypes        []domain.ProblemType
	studyTypes          []domain.StudyType
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string, jobRff jobG.RepositoryFactoryFunc, userRff userG.RepositoryFactoryFunc, pf service.ProcessorFactory, problemTypes []domain.ProblemType, studyTypes []domain.StudyType, problemRepositories map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error)) (service.RepositoryFactory, error) {
	if db == nil {
		return nil, libD.ErrInvalidArgument
	}

	return &repositoryFactory{
		db:                  db,
		driverName:          driverName,
		userRff:             userRff,
		pf:                  pf,
		problemRepositories: problemRepositories,
		problemTypes:        problemTypes,
		studyTypes:          studyTypes,
	}, nil
}

func (f *repositoryFactory) NewWorkbookRepository(ctx context.Context) (service.WorkbookRepository, error) {
	return NewWorkbookRepository(ctx, f.driverName, f, f.pf, f.db, f.problemTypes), nil
}

func (f *repositoryFactory) NewProblemRepository(ctx context.Context, problemType string) (service.ProblemRepository, error) {
	logger := log.FromContext(ctx)
	logger.Infof("problemType: %s", problemType)
	problemRepository, ok := f.problemRepositories[problemType]
	if !ok {
		logger.Errorf("problemTypes: %+v", f.problemRepositories)
		return nil, liberrors.Errorf("problem repository not found. problemType: %s", problemType)
	}
	return problemRepository(ctx, f.db)
}

func (f *repositoryFactory) NewProblemTypeRepository(ctx context.Context) (service.ProblemTypeRepository, error) {
	return NewProblemTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyTypeRepository(ctx context.Context) (service.StudyTypeRepository, error) {
	return NewStudyTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyRecordRepository(ctx context.Context) (service.StudyRecordRepository, error) {
	return NewStudyRecordRepository(ctx, f, f.db)
}

func (f *repositoryFactory) NewRecordbookRepository(ctx context.Context) (service.RecordbookRepository, error) {
	return NewRecordbookRepository(ctx, f, f.db, f.problemTypes, f.studyTypes)
}

func (f *repositoryFactory) NewUserQuotaRepository(ctx context.Context) (service.UserQuotaRepository, error) {
	return NewUserQuotaRepository(f.db)
}

func (f *repositoryFactory) NewStatRepository(ctx context.Context) (service.StatRepository, error) {
	return NewStatRepository(ctx, f.db)
}

func (f *repositoryFactory) NewStudyStatRepository(ctx context.Context) (service.StudyStatRepository, error) {
	return NewStudyStatRepository(ctx, f.db, f)
}

func (f *repositoryFactory) NewJobRepositoryFactory(ctx context.Context) (jobS.RepositoryFactory, error) {
	return jobG.NewRepositoryFactory(ctx, f.db)
}

func (f *repositoryFactory) NewUserRepositoryFactory(ctx context.Context) (userS.RepositoryFactory, error) {
	return userG.NewRepositoryFactory(ctx, f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)

type transaction struct {
	db  *gorm.DB
	rff RepositoryFactoryFunc
}

func NewTransaction(db *gorm.DB, rff RepositoryFactoryFunc) (service.Transaction, error) {
	return &transaction{
		db:  db,
		rff: rff,
	}, nil
}

func (t *transaction) Do(ctx context.Context, fn func(rf service.RepositoryFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		rf, err := t.rff(ctx, tx)
		if err != nil {
			return err
		}
		return fn(rf)
	})
}
