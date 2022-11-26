package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type StudyRecordRepository interface {
	AddRecord(ctx context.Context, operator userD.SystemOwnerModel, appUserID userD.AppUserID, workbookID domain.WorkbookID, problemTypeID uint, studyTypeID uint, problemID domain.ProblemID, mastered bool) error
}
