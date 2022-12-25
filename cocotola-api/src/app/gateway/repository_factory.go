package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	jobG "github.com/kujilabo/cocotola/cocotola-api/src/job/gateway"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type repositoryFactory struct {
	db                  *gorm.DB
	driverName          string
	location            *time.Location
	userRff             userG.RepositoryFactoryFunc
	pf                  service.ProcessorFactory
	problemRepositories map[domain.ProblemTypeName]func(context.Context, *gorm.DB) (service.ProblemRepository, error)
	problemTypes        ProblemTypes
	studyTypes          StudyTypes
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string, location *time.Location, jobRff jobG.RepositoryFactoryFunc, userRff userG.RepositoryFactoryFunc, pf service.ProcessorFactory, problemRepositories map[domain.ProblemTypeName]func(context.Context, *gorm.DB) (service.ProblemRepository, error)) (service.RepositoryFactory, error) {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	problemTypeRepo := newProblemTypeRepository(db)
	problemTypes, err := problemTypeRepo.FindAllProblemTypes(ctx)
	if err != nil {
		return nil, liberrors.Errorf("problemTypeRepo.FindAllProblemTypes. err: %w", err)
	}

	studyTypeRepo := newStudyTypeRepository(db)
	studyTypes, err := studyTypeRepo.FindAllStudyTypes(ctx)
	if err != nil {
		return nil, liberrors.Errorf("studyTypeRepo.FindAllStudyTypes. err: %w", err)
	}

	return &repositoryFactory{
		db:                  db,
		driverName:          driverName,
		location:            location,
		userRff:             userRff,
		pf:                  pf,
		problemRepositories: problemRepositories,
		problemTypes:        NewProblemTypes(problemTypes),
		studyTypes:          NewStudyTypes(studyTypes),
	}, nil
}

func (f *repositoryFactory) NewWorkbookRepository(ctx context.Context) service.WorkbookRepository {
	return newWorkbookRepository(ctx, f.driverName, f.pf, f.db, f, f.problemTypes)
}

func (f *repositoryFactory) NewProblemRepository(ctx context.Context, problemType domain.ProblemTypeName) (service.ProblemRepository, error) {
	logger := log.FromContext(ctx)
	logger.Infof("problemType: %s", problemType)
	problemRepository, ok := f.problemRepositories[problemType]
	if !ok {
		logger.Errorf("problemTypes: %+v", f.problemRepositories)
		return nil, liberrors.Errorf("problem repository not found. problemType: %s", problemType)
	}
	return problemRepository(ctx, f.db)
}

func (f *repositoryFactory) NewProblemTypeRepository(ctx context.Context) service.ProblemTypeRepository {
	return newProblemTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyTypeRepository(ctx context.Context) service.StudyTypeRepository {
	return newStudyTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyRecordRepository(ctx context.Context) service.StudyRecordRepository {
	return newStudyRecordRepository(ctx, f.db, f, f.problemTypes, f.studyTypes)
}

func (f *repositoryFactory) NewRecordbookRepository(ctx context.Context) service.RecordbookRepository {
	return newRecordbookRepository(ctx, f.db, f, f.problemTypes, f.studyTypes)
}

func (f *repositoryFactory) NewUserQuotaRepository(ctx context.Context) service.UserQuotaRepository {
	return newUserQuotaRepository(f.db, f.location)
}

func (f *repositoryFactory) NewStatRepository(ctx context.Context) service.StatRepository {
	return newStatRepository(ctx, f.db)
}

func (f *repositoryFactory) NewStudyStatRepository(ctx context.Context) service.StudyStatRepository {
	return newStudyStatRepository(ctx, f.db, f)
}

func (f *repositoryFactory) NewJobRepositoryFactory(ctx context.Context) (jobS.RepositoryFactory, error) {
	return jobG.NewRepositoryFactory(ctx, f.db) //nolint:wrapcheck
}

func (f *repositoryFactory) NewUserRepositoryFactory(ctx context.Context) (userS.RepositoryFactory, error) {
	return userG.NewRepositoryFactory(ctx, f.db) //nolint:wrapcheck
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
