package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type studyStatEntity struct {
	ID             string
	OrganizationID uint
	AppUserID      uint
	WorkbookID     uint
	ProblemTypeID  uint
	StudyTypeID    uint
	Answered       int
	Mastered       int
	RecordDate     time.Time
}

func (e *studyStatEntity) TableName() string {
	return "study_history"
}

type studyStatRepository struct {
	db *gorm.DB
	rf service.RepositoryFactory
}

func NewStudyStatRepository(ctx context.Context, db *gorm.DB, rf service.RepositoryFactory) (service.StudyStatRepository, error) {
	return &studyStatRepository{
		db: db,
		rf: rf,
	}, nil
}

func (r *studyStatRepository) AggregateResults(ctx context.Context, operator domain.SystemOwnerModel, targetDate time.Time, userID domain.AppUserID) error {
	studyRecordRepo, err := r.rf.NewStudyRecordRepository(ctx)
	if err != nil {
		return err
	}

	results, err := studyRecordRepo.CountAnsweredProblems(ctx, userID, targetDate)
	if err != nil {
		return err
	}

	for _, result := range results.Results {
		entity := studyStatEntity{
			ID:             libD.NewULID(),
			OrganizationID: uint(operator.GetOrganizationID()),
			AppUserID:      uint(userID),
			WorkbookID:     result.WorkbookID,
			ProblemTypeID:  result.ProblemTypeID,
			StudyTypeID:    result.StudyTypeID,
			Answered:       result.Answered,
			Mastered:       result.Mastered,
		}
		// Upsert
		if result := r.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "app_user_id"},
				{Name: "record_date"},
			}, // key colume
			DoUpdates: clause.AssignmentColumns([]string{
				"answered",
				"mastered",
			}), // column needed to be updated
		}).Create(&entity); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (r *studyStatRepository) CleanStudyStats(ctx context.Context, operator domain.SystemOwnerModel, expirationDate time.Time) error {

	studyStatEntity := studyStatEntity{}
	if result := r.db.
		Where("organization_id = ?", operator.GetOrganizationID()).
		Where("started_at < ?", expirationDate).
		Delete(&studyStatEntity); result.Error != nil {
		return result.Error
	}

	return nil
}
