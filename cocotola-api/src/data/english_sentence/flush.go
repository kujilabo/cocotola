package english_sentence

import (
	"context"
	"strconv"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/data"
	pluginEnglishDomain "github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

func CreateFlushWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "Flush", [][]string{
		{"This is a good book.", "これは良い本です。"},
	}); err != nil {
		return err
	}
	return nil
}

func CreateWorkbook(ctx context.Context, student appS.Student, workbookName string, sentences [][]string) error {
	logger := log.FromContext(ctx)

	workbookProperties := map[string]string{
		"audioEnabled": "false",
	}
	param, err := appD.NewWorkbookAddParameter(pluginEnglishDomain.EnglishSentenceProblemType, workbookName, appD.Lang2JA, "", workbookProperties)
	if err != nil {
		return liberrors.Errorf("NewWorkbookAddParameter. err: %w", err)
	}

	workbook, err := data.CreateWorkbookIfNotExists(ctx, student, workbookName, param)
	if err != nil {
		return liberrors.Errorf("createWorkbookIfNotExists. err: %w", err)
	}

	for i, sentence := range sentences {
		properties := map[string]string{
			"number":     strconv.Itoa(i + 1),
			"text":       sentence[0],
			"lang2":      "ja",
			"translated": sentence[1],
		}
		param, err := appS.NewProblemAddParameter(workbook.GetWorkbookID(), properties)
		if err != nil {
			return liberrors.Errorf("NewProblemAddParameter. err: %w", err)
		}

		added, _, _, err := workbook.AddProblem(ctx, student, param)
		if err != nil {
			return liberrors.Errorf("workbook.NewProblemAddParameter. err: %w", err)
		}
		logger.Infof("problemIDs: %v", added)
	}

	logger.Infof("Example %d", workbook.GetWorkbookID())
	return nil
}
