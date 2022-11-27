//go:generate mockery --output mock --name StudyRecordRepository
package service

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type StudyRecordRepository interface {
	AddRecord(ctx context.Context, operator userD.SystemOwnerModel, appUserID userD.AppUserID, workbookID domain.WorkbookID, problemTypeID uint, studyTypeID uint, problemID domain.ProblemID, mastered bool) error

	CountAnsweredProblems(ctx context.Context, targetUserID userD.AppUserID, targetDate time.Time) (*CountAnsweredResults, error)
}
