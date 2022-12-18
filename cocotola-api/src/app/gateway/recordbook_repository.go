package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type recordbookEntity struct {
	OrganizationID uint
	AppUserID      uint
	WorkbookID     uint
	ProblemTypeID  uint
	StudyTypeID    uint
	ProblemID      uint
	ResultPrev1    *bool
	ResultPrev2    *bool
	ResultPrev3    *bool
	Level          int
	Mastered       bool
	LastAnsweredAt time.Time
}

// type ProblemEntity interface {
// 	ToProblem() domain.Problem
// }

func (e *recordbookEntity) TableName() string {
	return "recordbook"
}

type recordbookRepository struct {
	db           *gorm.DB
	rf           service.RepositoryFactory
	problemTypes ProblemTypes
	studyTypes   StudyTypes
}

func newRecordbookRepository(ctx context.Context, db *gorm.DB, rf service.RepositoryFactory, problemTypes ProblemTypes, studyTypes StudyTypes) service.RecordbookRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}
	return &recordbookRepository{
		db:           db,
		rf:           rf,
		problemTypes: problemTypes,
		studyTypes:   studyTypes,
	}
}

func (r *recordbookRepository) FindStudyRecords(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType domain.StudyTypeName) (map[domain.ProblemID]domain.StudyRecord, error) {
	_, span := tracer.Start(ctx, "recordbookRepository.FindStudyResults")
	defer span.End()

	studyTypeID, err := r.studyTypes.ToStudyTypeID(studyType)
	if err != nil {
		return nil, liberrors.Errorf("unsupported studyType. studyType: %s", studyType)
	}

	var entities []recordbookEntity
	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("app_user_id = ?", operator.GetID()).
		Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make(map[domain.ProblemID]domain.StudyRecord)
	for _, e := range entities {
		results[domain.ProblemID(e.ProblemID)] = domain.StudyRecord{
			Level:          e.Level,
			ResultPrev1:    *e.ResultPrev1,
			Mastered:       e.Mastered,
			LastAnsweredAt: &e.LastAnsweredAt,
		}
	}

	return results, nil
}

func (r *recordbookRepository) SetResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType domain.StudyTypeName, problemType domain.ProblemTypeName, problemID domain.ProblemID, studyResult, mastered bool) error {
	ctx, span := tracer.Start(ctx, "recordbookRepository.SetResult")
	defer span.End()

	studyTypeID, err := r.studyTypes.ToStudyTypeID(studyType)
	if err != nil {
		return liberrors.Errorf("unsupported studyType. studyType: %s, err: %w", studyType, err)
	}

	problemTypeID, err := r.problemTypes.ToProblemTypeID(problemType)
	if err != nil {
		return liberrors.Errorf("unsupported problemType. problemType: %s, err:%w", problemType, err)
	}

	if mastered {
		return r.setMastered(ctx, operator, workbookID, studyTypeID, problemTypeID, problemID)
	}

	return r.setResult(ctx, operator, workbookID, studyTypeID, problemTypeID, problemID, studyResult)
}

func (r *recordbookRepository) setResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyTypeID uint, problemTypeID uint, problemID domain.ProblemID, studyResult bool) error {
	logger := log.FromContext(ctx)
	var entity recordbookEntity
	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Debugf("workbook_id = %d and study_type_id = %d and problem_id = %d", uint(workbookID), studyTypeID, uint(problemID))

			prev := false
			level := 0
			if studyResult {
				prev = true
				level = 1
			}
			entity = recordbookEntity{
				OrganizationID: uint(operator.GetOrganizationID()),
				AppUserID:      operator.GetID(),
				WorkbookID:     uint(workbookID),
				ProblemTypeID:  problemTypeID,
				StudyTypeID:    studyTypeID,
				ProblemID:      uint(problemID),
				ResultPrev1:    &prev,
				ResultPrev2:    nil,
				ResultPrev3:    nil,
				Level:          level,
				LastAnsweredAt: time.Now(),
			}
			if result := r.db.Create(&entity); result.Error != nil {
				return result.Error
			}
			return nil
		}
		return result.Error
	}

	if studyResult {
		if entity.Level < domain.StudyMaxLevel {
			entity.Level++
		}
	} else {
		if entity.Level > domain.StudyMinLevel {
			entity.Level--
		}
	}

	if entity.ResultPrev2 != nil {
		b := *entity.ResultPrev2
		entity.ResultPrev3 = &b
	}
	if entity.ResultPrev1 != nil {
		b := *entity.ResultPrev1
		entity.ResultPrev2 = &b
	}
	*entity.ResultPrev1 = studyResult
	entity.LastAnsweredAt = time.Now()

	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		Updates(&entity); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *recordbookRepository) setMastered(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyTypeID uint, problemTypeID uint, problemID domain.ProblemID) error {
	logger := log.FromContext(ctx)

	var entity recordbookEntity
	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Debugf("workbook_id = %d and study_type_id = %d and problem_id = %d", uint(workbookID), studyTypeID, uint(problemID))
			entity = recordbookEntity{
				AppUserID:      operator.GetID(),
				WorkbookID:     uint(workbookID),
				ProblemTypeID:  problemTypeID,
				StudyTypeID:    studyTypeID,
				ProblemID:      uint(problemID),
				ResultPrev1:    nil,
				ResultPrev2:    nil,
				ResultPrev3:    nil,
				Level:          0,
				Mastered:       true,
				LastAnsweredAt: time.Now(),
			}
			if result := r.db.Create(&entity); result.Error != nil {
				return result.Error
			}
			return nil
		}
		return result.Error
	}

	entity.Mastered = true
	entity.LastAnsweredAt = time.Now()

	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		Updates(&entity); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *recordbookRepository) CountMasteredProblems(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (map[domain.StudyTypeName]int, error) {
	_, span := tracer.Start(ctx, "recordbookRepository.CountMasteredProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	type studyTypeCountMap struct {
		StudyTypeID int
		Count       int
	}

	var results []studyTypeCountMap
	if result := r.db.Select("study_type_id, count(*) as count").
		Model(&recordbookEntity{}).
		Where("workbook_id = ?", uint(workbookID)).
		Where("app_user_id = ?", operator.GetID()).
		Where("mastered = ?", true).
		Group("study_type_id").
		Find(&results); result.Error != nil {
		return nil, result.Error
	}

	resultMap := make(map[domain.StudyTypeName]int)
	for _, studyType1 := range r.studyTypes.Values() {
		resultMap[studyType1.GetName()] = 0
		for _, result := range results {
			studyType2, err := r.studyTypes.ToStudyType(uint(result.StudyTypeID))
			if err != nil {
				return nil, err
			}
			if studyType1.GetName() == studyType2 {
				resultMap[studyType2] = result.Count
				break
			}
		}

	}

	logger.Debugf("CountMemorizedProblem. map: %+v", resultMap)

	return resultMap, nil
}
