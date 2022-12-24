//go:generate mockery --output mock --name StatModel
package domain

import (
	"time"

	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type StatHistoryResult struct {
	Date     time.Time
	Mastered int
	Answered int
}

type StatHistory struct {
	Results []StatHistoryResult
}

type StatModel interface {
	GetUserID() userD.AppUserID
	GetHistory() StatHistory
}

type statModel struct {
	UserID  userD.AppUserID `validate:"required"`
	History StatHistory     `validate:"required,dive"`
}

func NewStatModel(userID userD.AppUserID, history StatHistory) (StatModel, error) {
	m := &statModel{
		UserID:  userID,
		History: history,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *statModel) GetUserID() userD.AppUserID {
	return m.UserID
}

func (m *statModel) GetHistory() StatHistory {
	return m.History
}
