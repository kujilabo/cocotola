package gateway

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

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

var (
	MaxNumberOfProblemsToFindAllProblems  = 1000
	MaxNumberOfProblemIDsToFindProblemIDs = 100
)

type englishWordProblemEntity struct {
	ID                uint
	Version           int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CreatedBy         uint
	UpdatedBy         uint
	OrganizationID    uint
	WorkbookID        uint
	Number            int
	AudioID           uint
	Text              string
	Pos               int
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Lang2             string
	Translated        string
	PhraseID1         uint
	PhraseID2         uint
	SentenceID1       uint
	SentenceID2       uint
	// joined columns
	SentenceText1       string `gorm:"->"` // readonly
	SentenceTranslated1 string `gorm:"->"` // readonly
	SentenceNote1       string `gorm:"->"` // readonly
}

func (e *englishWordProblemEntity) TableName() string {
	return "english_word_problem"
}

func (e *englishWordProblemEntity) toProblem(ctx context.Context, synthesizerClient appS.SynthesizerClient) (service.EnglishWordProblem, error) {
	model, err := userD.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	properties := make(map[string]interface{})
	problemModel, err := appD.NewProblemModel(model, e.Number, domain.EnglishWordProblemType, properties)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	problem, err := appS.NewProblem(synthesizerClient, problemModel)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	lang2, err := appD.NewLang2(e.Lang2)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewLang2. lang2: %s, err: %w", e.Lang2, err)
	}

	phrases := make([]domain.EnglishPhraseProblemModel, 0)
	sentences := make([]domain.EnglishWordSentenceProblemModel, 0)
	if e.SentenceID1 != 0 {
		sentence, err := domain.NewEnglishWordProblemSentenceModel(appD.AudioID(0), e.SentenceText1, lang2, e.SentenceTranslated1, e.SentenceNote1)
		if err != nil {
			return nil, liberrors.Errorf(". err: %w", err)
		}
		sentences = append(sentences, sentence)
	}

	englishWordProblemModel, err := domain.NewEnglishWordProblemModel(problemModel, appD.AudioID(e.AudioID), e.Text, e.Pos, e.Phonetic, e.PresentThird, e.PresentParticiple, e.PastTense, e.PastParticiple, lang2, e.Translated, phrases, sentences)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	englishWordProblem, err := service.NewEnglishWordProblem(englishWordProblemModel, problem)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	logger := log.FromContext(ctx)
	// FIXME
	logger.Infof("properties: %+v", englishWordProblem.GetProperties(ctx))

	return englishWordProblem, nil
}

type englishWordProblemAddParemeter struct {
	Number            int `validate:"required"`
	AudioID           uint
	Text              string `validate:"required"`
	Pos               int    `validate:"required"`
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Lang2             string `validate:"required"`
	Translated        string
	PhraseID1         uint
	PhraseID2         uint
	SentenceID1       uint
	SentenceID2       uint
}

func toEnglishWordProblemAddParameter(param appD.ProblemAddParameter) (*englishWordProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["audioId"]; !ok {
		return nil, liberrors.Errorf("audioId is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["pos"]; !ok {
		return nil, liberrors.Errorf("pos is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, liberrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	if err != nil {
		return nil, liberrors.Errorf("audioId is not integer. audioId: %s, err: %w", param.GetProperties()["audioId"], err)
	}

	number := 1
	// number, err := param.GetIntProperty("number")
	// if err != nil {
	// 	return nil, err
	// }

	pos, err := strconv.Atoi(param.GetProperties()["pos"])
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	m := &englishWordProblemAddParemeter{
		Number:     number,
		AudioID:    uint(audioID),
		Lang2:      param.GetProperties()["lang2"],
		Text:       param.GetProperties()["text"],
		Pos:        pos,
		Translated: param.GetProperties()["translated"],
	}
	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

// type englishWordProblemUpdateParemeter struct {
// 	Number            int
// 	AudioID           uint
// 	Text              string `validate:"required"`
// 	Phonetic          string
// 	PresentThird      string
// 	PresentParticiple string
// 	PastTense         string
// 	PastParticiple    string
// 	Translated        string
// 	PhraseID1         uint
// 	PhraseID2         uint
// 	SentenceID1       uint
// 	SentenceID2       uint
// }

func toEnglishWordProblemUpdateParameter(newVersion int, updatedBy uint, param appD.ProblemUpdateParameter) (*englishWordProblemEntity, error) {
	if _, ok := param.GetProperties()[service.EnglishWordProblemUpdatePropertyAudioID]; !ok {
		return nil, liberrors.Errorf("audioId is not defined. err: %w", libD.ErrInvalidArgument)
	}

	text, err := param.GetStringProperty(service.EnglishWordProblemUpdatePropertyText)
	if err != nil {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	audioID, err := param.GetIntProperty(service.EnglishWordProblemUpdatePropertyAudioID)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	sentenceID, err := param.GetIntProperty(service.EnglishWordProblemUpdatePropertySentenceID1)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	m := &englishWordProblemEntity{
		Version:     newVersion,
		UpdatedBy:   updatedBy,
		AudioID:     uint(audioID),
		Text:        text,
		Translated:  param.GetProperties()[service.EnglishWordProblemUpdatePropertyTranslated],
		SentenceID1: uint(sentenceID),
	}
	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil

	// englishWordProblem := englishWordProblemEntity{
	// 	Version:           id.GetVersion() + 1,
	// 	UpdatedBy:         operator.GetID(),
	// 	AudioID:           problemParam.AudioID,
	// 	Number:            problemParam.Number,
	// 	Phonetic:          problemParam.Phonetic,
	// 	PresentThird:      problemParam.PresentThird,
	// 	PresentParticiple: problemParam.PresentParticiple,
	// 	PastTense:         problemParam.PastTense,
	// 	PastParticiple:    problemParam.PastParticiple,
	// 	Translated:        problemParam.Translated,
	// 	SentenceID1:       problemParam.SentenceID1,
	// }

}

type englishWordProblemRepository struct {
	db                *gorm.DB
	synthesizerClient appS.SynthesizerClient
	problemType       appD.ProblemTypeName
}

func NewEnglishWordProblemRepository(db *gorm.DB, synthesizerClient appS.SynthesizerClient, problemType appD.ProblemTypeName) (appS.ProblemRepository, error) {
	return &englishWordProblemRepository{
		db:                db,
		synthesizerClient: synthesizerClient,
		problemType:       problemType,
	}, nil
}

func (r *englishWordProblemRepository) FindProblems(ctx context.Context, operator appD.StudentModel, param appD.ProblemSearchCondition) (appD.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblems")
	defer span.End()

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(param.GetWorkbookID()))
	}

	var problemEntities []englishWordProblemEntity
	if result := where().Order("text, pos").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(ctx, count, problemEntities)
}

func (r *englishWordProblemRepository) FindAllProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (appD.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindAllProblems")
	defer span.End()

	limit := MaxNumberOfProblemsToFindAllProblems

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}

	var problemEntities []englishWordProblemEntity
	if result := where().Order("text, pos").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(ctx, count, problemEntities)
}

func (r *englishWordProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator appD.StudentModel, param appD.ProblemIDsCondition) (appD.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblemsByProblemIDs")
	defer span.End()

	if len(param.GetIDs()) > MaxNumberOfProblemIDsToFindProblemIDs {
		return nil, libD.ErrInvalidArgument
	}

	ids := make([]uint, 0)
	for _, id := range param.GetIDs() {
		ids = append(ids, uint(id))
	}

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID())).
		Where("id in ?", ids)

	var problemEntities []englishWordProblemEntity
	if result := db.Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(ctx, 0, problemEntities)
}

func (r *englishWordProblemRepository) toProblemSearchResult(ctx context.Context, count int64, problemEntities []englishWordProblemEntity) (appD.ProblemSearchResult, error) {
	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(ctx, r.synthesizerClient)
		if err != nil {
			return nil, liberrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	foundProblems, err := appD.NewProblemSearchResult(int(count), problems)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return foundProblems, nil
}

func (r *englishWordProblemRepository) FindProblemByID(ctx context.Context, operator appD.StudentModel, id appD.ProblemSelectParameter1) (appS.Problem, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblemByID")
	defer span.End()

	var problemEntity englishWordProblemEntity

	db := r.db.Table("english_word_problem AS T1").
		Select("T1.*,"+
			"T2.text AS sentence_text1,"+
			"T2.translated AS sentence_translated1,"+
			"T2.note AS sentence_note1").
		Joins("LEFT JOIN english_sentence_problem as T2 ON T1.sentence_id1 = T2.id").
		Where("T1.organization_id = ?", uint(operator.GetOrganizationID())).
		Where("T1.workbook_id = ?", uint(id.GetWorkbookID())).
		Where("T1.id = ?", uint(id.GetProblemID()))

	if result := db.First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, appS.ErrProblemNotFound
		}
		return nil, result.Error
	}

	foundProblem, err := problemEntity.toProblem(ctx, r.synthesizerClient)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return foundProblem, nil
}

func (r *englishWordProblemRepository) FindProblemIDs(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) ([]appD.ProblemID, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblemIDs")
	defer span.End()

	pageNo := 1
	pageSize := 1000
	limit := pageSize

	ids := make([]appD.ProblemID, 0)
	for {
		offset := (pageNo - 1) * pageSize

		where := r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))

		var problemEntities []englishWordProblemEntity
		if result := where.Order("text, pos").
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

func (r *englishWordProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator appD.StudentModel, condition interface{}) ([]appD.ProblemModel, error) {
	return nil, errors.New("not implement")
}

func (r *englishWordProblemRepository) AddProblem(ctx context.Context, operator appD.StudentModel, param appD.ProblemAddParameter) (appD.ProblemID, error) {
	ctx, span := tracer.Start(ctx, "englishWordProblemRepository.AddProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	problemParam, err := toEnglishWordProblemAddParameter(param)
	if err != nil {
		return 0, liberrors.Errorf("failed to toEnglishWordProblemAddParameter. param: %+v, err: %w", param, err)
	}

	englishWordProblem := englishWordProblemEntity{
		Version:           1,
		CreatedBy:         operator.GetID(),
		UpdatedBy:         operator.GetID(),
		OrganizationID:    uint(operator.GetOrganizationID()),
		WorkbookID:        uint(param.GetWorkbookID()),
		AudioID:           problemParam.AudioID,
		Number:            problemParam.Number, //            param.GetNumber(),
		Text:              problemParam.Text,
		Pos:               problemParam.Pos,
		Phonetic:          problemParam.Phonetic,
		PresentThird:      problemParam.PresentThird,
		PresentParticiple: problemParam.PresentParticiple,
		PastTense:         problemParam.PastTense,
		PastParticiple:    problemParam.PastParticiple,
		Lang2:             problemParam.Lang2,
		Translated:        problemParam.Translated,
	}

	logger.Infof("englishWordProblemRepository.AddProblem. text: %s", problemParam.Text)
	if result := r.db.Create(&englishWordProblem); result.Error != nil {
		return 0, liberrors.Errorf("failed to Create. param: %+v, err: %w", param, libG.ConvertDuplicatedError(result.Error, appS.ErrProblemAlreadyExists))
	}

	return appD.ProblemID(englishWordProblem.ID), nil
}

func (r *englishWordProblemRepository) UpdateProblem(ctx context.Context, operator appD.StudentModel, id appD.ProblemSelectParameter2, param appD.ProblemUpdateParameter) error {
	ctx, span := tracer.Start(ctx, "englishWordProblemRepository.UpdateProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	englishWordProblem, err := toEnglishWordProblemUpdateParameter(id.GetVersion()+1, operator.GetID(), param)
	if err != nil {
		return liberrors.Errorf("failed to toEnglishWordProblemUdateParameter. param: %+v, err: %w", param, err)
	}

	// englishWordProblem := englishWordProblemEntity{
	// 	Version:           id.GetVersion() + 1,
	// 	UpdatedBy:         operator.GetID(),
	// 	AudioID:           problemParam.AudioID,
	// 	Number:            problemParam.Number,
	// 	Phonetic:          problemParam.Phonetic,
	// 	PresentThird:      problemParam.PresentThird,
	// 	PresentParticiple: problemParam.PresentParticiple,
	// 	PastTense:         problemParam.PastTense,
	// 	PastParticiple:    problemParam.PastParticiple,
	// 	Translated:        problemParam.Translated,
	// 	SentenceID1:       problemParam.SentenceID1,
	// }

	logger.Infof("englishWordProblemRepository.UpdateProblem. text: %s", englishWordProblem.Text)

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		UpdateColumns(&englishWordProblem)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishWordProblemRepository) UpdateProblemProperty(ctx context.Context, operator appD.StudentModel, id appD.ProblemSelectParameter2, param appD.ProblemUpdateParameter) error {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.UpdateProblem")
	defer span.End()
	return errors.New("not implemented")
}

func (r *englishWordProblemRepository) RemoveProblem(ctx context.Context, operator appD.StudentModel, id appD.ProblemSelectParameter2) error {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.RemoveProblem")
	defer span.End()

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		Delete(&englishWordProblemEntity{})

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishWordProblemRepository) CountProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (int, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.CountProblems")
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
