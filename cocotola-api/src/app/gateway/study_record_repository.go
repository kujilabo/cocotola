package gateway

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type studyRecordEntity struct {
	AppUserID     uint
	WorkbookID    uint
	ProblemTypeID uint
	StudyTypeID   uint
	ProblemID     uint
	Mastered      bool
	RecordDate    time.Time
}

// type ProblemEntity interface {
// 	ToProblem() domain.Problem
// }

func (e *studyRecordEntity) TableName() string {
	return "study_record"
}

type studyRecordRepository struct {
	rf service.RepositoryFactory
	db *gorm.DB
}

func NewStudyRecordRepository(ctx context.Context, rf service.RepositoryFactory, db *gorm.DB) service.StudyRecordRepository {
	return &studyRecordRepository{
		rf: rf,
		db: db,
	}
}

func (r *studyRecordRepository) AddRecord(ctx context.Context, operator userD.SystemOwnerModel, appUserID userD.AppUserID, workbookID domain.WorkbookID, problemTypeID uint, studyTypeID uint, problemID domain.ProblemID, mastered bool) error {

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	entity := studyRecordEntity{
		AppUserID:     uint(appUserID),
		WorkbookID:    uint(workbookID),
		ProblemTypeID: problemTypeID,
		StudyTypeID:   studyTypeID,
		ProblemID:     uint(problemID),
		Mastered:      mastered,
		RecordDate:    today,
	}

	// Upsert
	if result := r.db.Debug().
		Clauses(clause.OnConflict{
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

func (r *studyRecordRepository) CountAnsweredProblems(ctx context.Context, targetUserID userD.AppUserID, targetDate time.Time) (*service.CountAnsweredResults, error) {
	_, span := tracer.Start(ctx, "recordbookRepository.CountMasteredProblem")
	defer span.End()
	// logger := log.FromContext(ctx)
	// {
	// 	var e []studyRecordEntity
	// 	r.db.Find(&e)
	// 	for _, e2 := range e {
	// 		logger.Infof("%+v", e2)
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
