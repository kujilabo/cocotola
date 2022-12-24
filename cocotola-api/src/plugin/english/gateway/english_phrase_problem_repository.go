package gateway

import (
	"context"
	"errors"
	"math"
	"strconv"

	"gorm.io/gorm"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
	"github.com/kujilabo/cocotola/lib/log"
)

type englishPhraseProblemEntity struct {
	SinmpleModelEntity
	OrganizationID uint
	WorkbookID     uint
	Number         int
	AudioID        uint
	Text           string
	Lang2          string
	Translated     string
}

func (e *englishPhraseProblemEntity) TableName() string {
	return "english_phrase_problem"
}

func (e *englishPhraseProblemEntity) toProblem(synthesizerClient appS.SynthesizerClient) (service.EnglishPhraseProblem, error) {
	model, err := userD.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf("userD.NewModel. err: %w", err)
	}

	properties := make(map[string]interface{})
	for k, v := range toEnglishPhraseProblemProperties(e.Lang2, e.Text, e.Translated) {
		properties[k] = v
	}

	problemModel, err := appD.NewProblemModel(model, e.Number, domain.EnglishPhraseProblemType, properties)
	if err != nil {
		return nil, liberrors.Errorf("appD.NewProblemModel. err: %w", err)
	}

	problem, err := appS.NewProblem(synthesizerClient, problemModel)
	if err != nil {
		return nil, liberrors.Errorf("appD.NewProblem. err: %w", err)
	}

	lang2, err := appD.NewLang2(e.Lang2)
	if err != nil {
		return nil, liberrors.Errorf("appD.NewLang2. err: %w", err)
	}

	englishPhraseProblemModel, err := domain.NewEnglishPhraseProblemModel(problemModel, appD.AudioID(e.AudioID), e.Text, lang2, e.Translated)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewEnglishPhraseProblemModel. err: %w", err)
	}

	phraseProblem, err := service.NewEnglishPhraseProblem(englishPhraseProblemModel, problem)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return phraseProblem, nil
}

func fromEnglishPhraseProblemProperties(properties map[string]string) (string, string, string) {
	return properties["lang2"], properties["text"], properties["translated"]
}

func toEnglishPhraseProblemProperties(lang2, text, translated string) map[string]string {
	return map[string]string{
		"lang2":      lang2,
		"text":       text,
		"translated": translated,
	}
}

type englishPhraseProblemAddParameter struct {
	Number     int
	AudioID    uint
	Lang2      string
	Text       string
	Translated string
}

func toNewEnglishPhraseProblemParam(param appS.ProblemAddParameter) (*englishPhraseProblemAddParameter, error) {
	audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	if err != nil {
		return nil, liberrors.Errorf("strconv.Atoi. err: %w", err)
	}

	number := 1
	// number, err := param.GetIntProperty("number")
	// if err != nil {
	// 	return nil, err
	// }

	lang2, text, translated := fromEnglishPhraseProblemProperties(param.GetProperties())
	m := &englishPhraseProblemAddParameter{
		Number:     number,
		AudioID:    uint(audioID),
		Lang2:      lang2,
		Text:       text,
		Translated: translated,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

type englishPhraseProblemRepository struct {
	db                *gorm.DB
	synthesizerClient appS.SynthesizerClient
	problemType       appD.ProblemTypeName
}

func NewEnglishPhraseProblemRepository(db *gorm.DB, synthesizerClient appS.SynthesizerClient, problemType appD.ProblemTypeName) (appS.ProblemRepository, error) {
	return &englishPhraseProblemRepository{
		db:                db,
		synthesizerClient: synthesizerClient,
		problemType:       problemType,
	}, nil
}

func (r *englishPhraseProblemRepository) FindProblems(ctx context.Context, operator appD.StudentModel, param appS.ProblemSearchCondition) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishPhraseProblemRepository.FindProblems")
	defer span.End()

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	var problemEntities []englishPhraseProblemEntity

	where := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID()))

	if result := where.Order("workbook_id, number, created_at").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, liberrors.Errorf("failed to Find. err: %w", result.Error)
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.synthesizerClient)
		if err != nil {
			return nil, liberrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	var count int64
	if result := where.Model(&englishPhraseProblemEntity{}).Count(&count); result.Error != nil {
		return nil, liberrors.Errorf("failed to Count. err: %w", result.Error)
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	foundProblems, err := appS.NewProblemSearchResult(int(count), problems)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return foundProblems, nil
}

func (r *englishPhraseProblemRepository) FindAllProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishPhraseProblemRepository.FindAllProblems")
	defer span.End()

	limit := 1000
	var problemEntities []englishPhraseProblemEntity

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}
	if result := where().Order("workbook_id, number, text, created_at").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, liberrors.Errorf("failed to Find. err: %w", result.Error)
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.synthesizerClient)
		if err != nil {
			return nil, liberrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	var count int64
	if result := where().Model(&englishPhraseProblemEntity{}).Count(&count); result.Error != nil {
		return nil, liberrors.Errorf("failed to Count. err: %w", result.Error)
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	foundProblems, err := appS.NewProblemSearchResult(int(count), problems)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return foundProblems, nil
}

func (r *englishPhraseProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator appD.StudentModel, param appS.ProblemIDsCondition) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishPhraseProblemRepository.FindProblemsByProblemIDs")
	defer span.End()

	var problemEntities []englishPhraseProblemEntity

	ids := make([]uint, 0)
	for _, id := range param.GetIDs() {
		ids = append(ids, uint(id))
	}

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID())).
		Where("id in ?", ids)

	if result := db.Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.synthesizerClient)
		if err != nil {
			return nil, liberrors.Errorf("e.toProblem. err: %w", err)
		}
		problems[i] = p
	}

	foundProblems, err := appS.NewProblemSearchResult(0, problems)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return foundProblems, nil
}

func (r *englishPhraseProblemRepository) FindProblemByID(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter1) (appS.Problem, error) {
	_, span := tracer.Start(ctx, "englishPhraseProblemRepository.FindProblemByID")
	defer span.End()

	var problemEntity englishPhraseProblemEntity

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID()))

	if result := db.First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, appS.ErrProblemNotFound
		}
		return nil, result.Error
	}

	problem, err := problemEntity.toProblem(r.synthesizerClient)
	if err != nil {
		return nil, err
	}

	return problem, nil
}

func (r *englishPhraseProblemRepository) FindProblemIDs(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) ([]appD.ProblemID, error) {
	_, span := tracer.Start(ctx, "englishPhraseProblemRepository.FindProblemIDs")
	defer span.End()

	pageNo := 1
	pageSize := 1000
	ids := make([]appD.ProblemID, 0)
	for {
		limit := pageSize
		offset := (pageNo - 1) * pageSize
		var problemEntities []englishPhraseProblemEntity

		where := r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))

		if result := where.Order("workbook_id, number, text, created_at").
			Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
			return nil, result.Error
		}

		if len(problemEntities) == 0 {
			break
		}

		for _, r := range problemEntities {
			ids = append(ids, appD.ProblemID(r.ID))
		}

		pageNo++
	}

	return ids, nil
}

func (r *englishPhraseProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator appD.StudentModel, condition interface{}) ([]appD.ProblemModel, error) {
	return nil, errors.New("not implement")
}

func (r *englishPhraseProblemRepository) AddProblem(ctx context.Context, operator appD.StudentModel, param appS.ProblemAddParameter) (appD.ProblemID, error) {
	ctx, span := tracer.Start(ctx, "englishPhraseProblemRepository.AddProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	problemParam, err := toNewEnglishPhraseProblemParam(param)
	if err != nil {
		return 0, liberrors.Errorf("toNewEnglishPhraseProblemParam. err: %w", err)
	}
	englishPhraseProblem := englishPhraseProblemEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
		WorkbookID:     uint(param.GetWorkbookID()),
		AudioID:        problemParam.AudioID,
		Number:         problemParam.Number,
		Text:           problemParam.Text,
		Lang2:          problemParam.Lang2,
		Translated:     problemParam.Translated,
	}

	logger.Infof("englishPhraseProblemRepository.AddProblem. lang2: %s, text: %s", problemParam.Lang2, problemParam.Text)
	if result := r.db.Create(&englishPhraseProblem); result.Error != nil {
		return 0, liberrors.Errorf(". err: %w", libG.ConvertDuplicatedError(result.Error, appS.ErrProblemAlreadyExists))
	}

	return appD.ProblemID(englishPhraseProblem.ID), nil
}

func (r *englishPhraseProblemRepository) UpdateProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) error {
	return errors.New("not implemented")
}

func (r *englishPhraseProblemRepository) UpdateProblemProperty(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) error {
	return errors.New("not implemented")
}

func (r *englishPhraseProblemRepository) RemoveProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2) error {
	ctx, span := tracer.Start(ctx, "englishPhraseProblemRepository.RemoveProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	logger.Infof("englishPhraseProblemRepository.RemoveProblem. text: %d", id.GetProblemID())

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		Delete(&englishPhraseProblemEntity{})

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishPhraseProblemRepository) CountProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (int, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.CountProblems")
	defer span.End()

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}
