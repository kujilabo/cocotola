package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type studyRecordEntity struct {
	ID             string
	OrganizationID uint
	AppUserID      uint
	WorkbookID     uint
	ProblemTypeID  uint
	StudyTypeID    uint
	ProblemID      uint
	Mastered       bool
	RecordDate     time.Time
}

// type ProblemEntity interface {
// 	ToProblem() domain.Problem
// }

func (e *studyRecordEntity) TableName() string {
	return "study_record"
}

type studyRecordRepository struct {
	db           *gorm.DB
	rf           service.RepositoryFactory
	problemTypes ProblemTypes
	studyTypes   StudyTypes
}

func newStudyRecordRepository(ctx context.Context, db *gorm.DB, rf service.RepositoryFactory, problemTypes ProblemTypes, studyTypes StudyTypes) service.StudyRecordRepository {
	return &studyRecordRepository{
		db:           db,
		rf:           rf,
		problemTypes: problemTypes,
		studyTypes:   studyTypes,
	}
}

func (r *studyRecordRepository) AddRecord(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, problemType domain.ProblemTypeName, studyType domain.StudyTypeName, problemID domain.ProblemID, mastered bool) error {

	problemTypeID, err := r.problemTypes.ToProblemTypeID(problemType)
	if err != nil {
		return liberrors.Errorf("r.problemTypes.ToProblemTypeID. err: %w", err)
	}
	studyTypeID, err := r.studyTypes.ToStudyTypeID(studyType)
	if err != nil {
		return liberrors.Errorf("r.studyTypes.ToStudyTypeID. err: %w", err)
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	entity := studyRecordEntity{
		ID:             libD.NewULID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		AppUserID:      uint(operator.GetAppUserID()),
		WorkbookID:     uint(workbookID),
		ProblemTypeID:  problemTypeID,
		StudyTypeID:    studyTypeID,
		ProblemID:      uint(problemID),
		Mastered:       mastered,
		RecordDate:     today,
	}

	// Upsert
	if result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "app_user_id"},
			{Name: "workbook_id"},
			{Name: "study_type_id"},
			{Name: "problem_type_id"},
			{Name: "problem_id"},
			{Name: "record_date"},
		}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{
			"mastered",
		}), // column needed to be updated
	}).Create(&entity); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *studyRecordRepository) CountAnsweredProblems(ctx context.Context, operator userD.SystemOwnerModel, targetUserID userD.AppUserID, targetDate time.Time) (*service.CountAnsweredResults, error) {
	_, span := tracer.Start(ctx, "recordbookRepository.CountMasteredProblem")
	defer span.End()
	// logger := log.FromContext(ctx)
	// {
	// 	var e []studyRecordEntity
	// 	r.db.Find(&e)
	// 	for _, e2 := range e {
	// 		logger.Infof("%+v", e2)
	// 		logger.Infof("%s", e2.RecordDate.Format(time.RFC3339))
	// 	}
	// }

	type countEntity struct {
		WorkbookID    uint
		ProblemTypeID uint
		StudyTypeID   uint
		Answered      int
		Mastered      int
	}
	var entities []countEntity

	if result := r.db.Select("count(*) as answered, sum(mastered) as mastered, workbook_id, problem_type_id, study_type_id").
		Model(&studyRecordEntity{}).
		Where("app_user_id = ?", uint(targetUserID)).
		Where("record_date = ?", targetDate).
		Group("workbook_id").
		Group("problem_type_id").
		Group("study_type_id").
		Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]service.CountAnsweredResult, len(entities))
	for i, entity := range entities {
		results[i] = service.CountAnsweredResult{
			WorkbookID:    entity.WorkbookID,
			ProblemTypeID: entity.ProblemTypeID,
			StudyTypeID:   entity.StudyTypeID,
			Answered:      entity.Answered,
			Mastered:      entity.Mastered,
		}
	}

	return &service.CountAnsweredResults{
		Results: results,
	}, nil
}
