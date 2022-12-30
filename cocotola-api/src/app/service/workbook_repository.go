//go:generate mockery --output mock --name WorkbookRepository
//go:generate mockery --output mock --name WorkbookSearchCondition
//go:generate mockery --output mock --name WorkbookSearchResult
//go:generate mockery --output mock --name WorkbookAddParameter
//go:generate mockery --output mock --name WorkbookUpdateParameter
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

var ErrWorkbookNotFound = errors.New("workbook not found")
var ErrWorkbookAlreadyExists = errors.New("workbook already exists")
var ErrWorkbookPermissionDenied = errors.New("permission denied")

type WorkbookRepository interface {
	FindPersonalWorkbooks(ctx context.Context, operator domain.StudentModel, param domain.WorkbookSearchCondition) (domain.WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID) (Workbook, error)

	FindWorkbookByName(ctx context.Context, operator userD.AppUserModel, spaceID userD.SpaceID, name string) (Workbook, error)

	AddWorkbook(ctx context.Context, operator userD.AppUserModel, spaceID userD.SpaceID, param domain.WorkbookAddParameter) (domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, version int, param domain.WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, version int) error
}
