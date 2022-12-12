//go:generate mockery --output mock --name StudentUsecaseAudio
package student

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/usecase"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type StudentUsecaseAudio interface {
	FindAudioByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (service.Audio, error)
}

type studentUsecaseAudio struct {
	transaction       service.Transaction
	pf                service.ProcessorFactory
	synthesizerClient service.SynthesizerClient
}

func NewStudentUsecaseAudio(transaction service.Transaction, pf service.ProcessorFactory, synthesizerClient service.SynthesizerClient) StudentUsecaseAudio {
	return &studentUsecaseAudio{
		transaction:       transaction,
		pf:                pf,
		synthesizerClient: synthesizerClient,
	}
}

func (s *studentUsecaseAudio) FindAudioByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (service.Audio, error) {
	var result service.Audio
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, userRf, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}

		problem, err := workbook.FindProblemByID(ctx, student, problemID)
		if err != nil {
			return err
		}

		savedAudioID, ok := (problem.GetProperties(ctx)["audioId"]).(domain.AudioID)
		if !ok {
			return errors.New("mismatch")
		}

		logger := log.FromContext(ctx)
		if audioID != savedAudioID {
			logger.Debugf("properties: %+v", problem.GetProperties(ctx))
			logger.Warnf("audioID: %d, %s", audioID, problem.GetProperties(ctx)["audioId"])
			message := "invalid audio id"
			return domain.NewPluginError(domain.ErrorType(domain.ErrorTypeClient), message, []string{message}, libD.ErrInvalidArgument)
		}

		tmpResult, err := s.synthesizerClient.FindAudioByAudioID(ctx, audioID)
		if err != nil {
			return err
		}

		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *studentUsecaseAudio) findStudentAndWorkbook(ctx context.Context, rf service.RepositoryFactory, userRf userS.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (service.Student, service.Workbook, error) {
	studentService, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
	if err != nil {
		return nil, nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}
	workbookService, err := studentService.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, err
	}
	return studentService, workbookService, nil
}
