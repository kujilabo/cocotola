package gateway

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const userFetchSize = 10

type studyStatEntity struct {
	AppUserID     uint
	WorkbookID    uint
	ProblemTypeID uint
	StudyTypeID   uint
	Answered      int
	Mastered      int
	RecordDate    time.Time
}

func (e *studyStatEntity) TableName() string {
	return "study_history"
}

type studyStatRepository struct {
	db     *gorm.DB
	rf     service.RepositoryFactory
	userRf userS.RepositoryFactory
}

func NewStudyStatRepository(ctx context.Context, rf service.RepositoryFactory, db *gorm.DB, userRf userS.RepositoryFactory) (service.StudyStatRepository, error) {
	return &studyStatRepository{
		rf:     rf,
		db:     db,
		userRf: userRf,
	}, nil
}

func (r *studyStatRepository) AggregateResultsOfAllUsers(ctx context.Context, operator domain.SystemOwnerModel, targetDate time.Time) error {
	userRepo, err := r.userRf.NewAppUserRepository()
	if err != nil {
		return err
	}
	studyRecordRepo, err := r.rf.NewStudyRecordRepository(ctx)
	if err != nil {
		return err
	}

	pageNo := 1
	pageSize := userFetchSize
	for {
		userIDs, err := userRepo.FindAppUserIDs(ctx, operator, pageNo, pageSize)
		if err != nil {
			return err
		}
		if len(userIDs) == 0 {
			break
		}
		for _, userID := range userIDs {
			results, err := studyRecordRepo.CountAnsweredProblems(ctx, userID, targetDate)
			if err != nil {
				return err
			}

			for _, result := range results.Results {
				entity := studyStatEntity{
					AppUserID:     uint(userID),
					WorkbookID:    result.WorkbookID,
					ProblemTypeID: result.ProblemTypeID,
					StudyTypeID:   result.StudyTypeID,
					Answered:      result.Answered,
					Mastered:      result.Mastered,
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
		}

		pageNo++
	}

	return nil
}
