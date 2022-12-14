package english_sentence

import (
	"context"
	"errors"
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

	problems, err := workbook.FindAllProblems(ctx, student)
	if err != nil {
		return liberrors.Errorf("createWorkbookIfNotExists. err: %w", err)
	}

	problemMap := make(map[string]struct{})
	for _, problem := range problems.GetResults() {
		properties := problem.GetProperties(ctx)
		text1, ok := properties["text"]
		if !ok {
			continue
		}
		text2, ok := text1.(string)
		if !ok {
			continue
		}
		problemMap[text2] = struct{}{}
	}

	for i, sentence := range sentences {
		// skip if the sentence is already registered
		if _, ok := problemMap[sentence[0]]; ok {
			logger.Infof("Skip %s", sentence[0])
			continue
		}

		properties := map[string]string{
			"number":     strconv.Itoa(i + 1),
			"text":       sentence[0],
			"lang2":      "ja",
			"translated": sentence[1],
		}
		param, err := appD.NewProblemAddParameter(workbook.GetWorkbookID(), properties)
		if err != nil {
			return liberrors.Errorf("NewProblemAddParameter. err: %w", err)
		}

		added, _, _, err := workbook.AddProblem(ctx, student, param)
		if err != nil && !errors.Is(err, appS.ErrProblemAlreadyExists) {
			return liberrors.Errorf("AddProblem. err: %w", err)
		}
		logger.Infof("problemIDs: %v", added)
	}

	logger.Infof("Example %d", workbook.GetWorkbookID())
	return nil
}
