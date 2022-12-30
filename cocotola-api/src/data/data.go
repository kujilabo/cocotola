package data

import (
	"context"
	"errors"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func CreateWorkbookIfNotExists(ctx context.Context, student appS.Student, workbookName string, param appD.WorkbookAddParameter) (appS.Workbook, error) {
	tmpWorkbook1, err := student.FindWorkbookByName(ctx, workbookName)
	if err == nil {
		return tmpWorkbook1, nil
	}

	if !errors.Is(err, appS.ErrWorkbookNotFound) {
		return nil, liberrors.Errorf("student.FindWorkbookByName. err: %w", err)
	}

	workbookID, err := student.AddWorkbookToPersonalSpace(ctx, param)
	if err != nil {
		return nil, liberrors.Errorf("failed to AddWorkbookToPersonalSpace. err: %w", err)
	}

	tmpWorkbook2, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindWorkbookByID. err: %w", err)
	}

	return tmpWorkbook2, nil
}
