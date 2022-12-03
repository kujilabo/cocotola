package gateway

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
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

func NewStatRepository(ctx context.Context, db *gorm.DB) (service.StatRepository, error) {
	return &statRepository{
		db: db,
	}, nil
}

func (r *statRepository) FindStat(ctx context.Context, operatorID userD.AppUserID) (service.Stat, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	dateFormat := "2006-01-02"

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

	m := map[string]statEntity{}
	for i := 0; i < hisotrySize; i++ {
		t := startDate.AddDate(0, 0, i)
		s := t.Format(dateFormat)
		m[s] = statEntity{
			AppUserID: uint(operatorID),
			Answered:  0,
			Mastered:  0,
		}
	}

	for _, entity := range entities {
		m[entity.RecordDate.Format(dateFormat)] = statEntity{
			AppUserID: entity.AppUserID,
			Answered:  entity.Answered,
			Mastered:  entity.Mastered,
		}
	}

	results := make([]domain.StatHistoryResult, hisotrySize)
	for i := 0; i < hisotrySize; i++ {
		t := startDate.AddDate(0, 0, i)
		s := t.Format(dateFormat)
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
