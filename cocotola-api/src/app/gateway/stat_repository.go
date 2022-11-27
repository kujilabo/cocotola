package gateway

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const hisotrySize = 7

type statEntity struct {
	AppUserID  uint
	RecordDate time.Time
	Answered   int
	Mastered   int
}

type statRepository struct {
	db *gorm.DB
}

func NewStatRepository(ctx context.Context, db *gorm.DB) service.StatRepository {
	return &statRepository{
		db: db,
	}
}
func (r *statRepository) FindStat(ctx context.Context, operatorID userD.AppUserID) (service.Stat, error) {

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	startDate := today.AddDate(0, 0, -hisotrySize)
	endDate := today.AddDate(0, 0, -1)
	var entities []statEntity
	if result := r.db.Debug().Select("app_user_id, record_date, sum(answered) as answered, sum(mastered) as mastered").
		Table("study_stat").
		Where("app_user_id = ?", uint(operatorID)).
		Where("? <= record_date", startDate).
		Where("record_date <= ? ", endDate).
		Group("app_user_id, record_date").
		Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	logrus.Warnf("entities: %+v", entities)

	m := map[string]statEntity{}
	for i := 0; i < hisotrySize; i++ {
		t := startDate.AddDate(0, 0, i)
		s := t.Format(time.RFC3339)
		m[s] = statEntity{
			AppUserID: uint(operatorID),
			Answered:  0,
			Mastered:  0,
		}
	}

	for _, entity := range entities {
		m[entity.RecordDate.Format(time.RFC3339)] = statEntity{
			AppUserID: entity.AppUserID,
			Answered:  entity.Answered,
			Mastered:  entity.Mastered,
		}
	}

	results := make([]domain.StatHistoryResult, hisotrySize)
	for i := 0; i < hisotrySize; i++ {
		t := startDate.AddDate(0, 0, i)
		s := t.Format(time.RFC3339)
		results[i] = domain.StatHistoryResult{
			Date:     t,
			Mastered: m[s].Mastered,
			Answered: m[s].Answered,
		}
	}
	history := domain.StatHistory{
		Results: results,
	}

	model, err := domain.NewStatModel(operatorID, history)
	if err != nil {
		return nil, err
	}

	return model, nil
}
