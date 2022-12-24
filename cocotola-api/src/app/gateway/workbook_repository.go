package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/gateway/casbinquery"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	libG "github.com/kujilabo/cocotola/lib/gateway"
	"github.com/kujilabo/cocotola/lib/log"

	// casbinquery "github.com/pecolynx/casbin-query"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type workbookEntity struct {
	ID             uint
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	SpaceID        uint
	OwnerID        uint
	Name           string
	Lang2          string
	ProblemTypeID  uint `gorm:"column:problem_type_id"`
	QuestionText   string
	Properties     string
}

func (e *workbookEntity) TableName() string {
	return "workbook"
}

func jsonToStringMap(s string) (map[string]string, error) {
	var m map[string]string
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, liberrors.Errorf("json.Unmarshal. err: %w", err)
	}
	return m, nil
}

func stringMapToJSON(m map[string]string) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", liberrors.Errorf("json.Marshal. err: %w", err)
	}
	return string(b), nil
}

func (e *workbookEntity) toWorkbookModel(rf service.RepositoryFactory, pf service.ProcessorFactory, operator userD.AppUserModel, problemType domain.ProblemTypeName, privs userD.Privileges) (domain.WorkbookModel, error) {
	model, err := userD.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf("userD.NewModel. err: %w", err)
	}

	properties, err := jsonToStringMap(e.Properties)
	if err != nil {
		return nil, liberrors.Errorf("failed to jsonToStringMap. err: %w ", err)
	}

	lang2, err := domain.NewLang2(e.Lang2)
	if err != nil {
		return nil, liberrors.Errorf("invalid lang2. lang2: %s, err: %w", e.Lang2, err)
	}

	workbook, err := domain.NewWorkbookModel(model, userD.SpaceID(e.SpaceID), userD.AppUserID(e.OwnerID), privs, e.Name, lang2, problemType, e.QuestionText, properties)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewWorkbook. entity: %+v, err: %w", e, err)
	}
	return workbook, nil
}

type workbookRepository struct {
	driverName   string
	db           *gorm.DB
	rf           service.RepositoryFactory
	pf           service.ProcessorFactory
	problemTypes ProblemTypes
}

func newWorkbookRepository(ctx context.Context, driverName string, pf service.ProcessorFactory, db *gorm.DB, rf service.RepositoryFactory, problemTypes ProblemTypes) service.WorkbookRepository {
	return &workbookRepository{
		driverName:   driverName,
		db:           db,
		rf:           rf,
		pf:           pf,
		problemTypes: problemTypes,
	}
}

func (r *workbookRepository) FindPersonalWorkbooks(ctx context.Context, operator domain.StudentModel, param service.WorkbookSearchCondition) (service.WorkbookSearchResult, error) {
	ctx, span := tracer.Start(ctx, "workbookRepository.FindPersonalWorkbooks")
	defer span.End()

	logger := log.FromContext(ctx)
	logger.Debugf("workbookRepository.FindPersonalWorkbooks. OperatorID: %d", operator.GetID())

	if param == nil {
		return nil, libD.ErrInvalidArgument
	}

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	workbookEntities := []workbookEntity{}

	objectColumnName := "name"
	subQuery, err := casbinquery.QueryObject(r.db, r.driverName, domain.WorkbookObjectPrefix, objectColumnName, "user_"+strconv.Itoa(int(operator.GetID())), "read")
	if err != nil {
		return nil, liberrors.Errorf("casbinquery.QueryObject. err: %w", err)
	}

	if result := r.db.Model(&workbookEntity{}).
		Joins("inner join (?) AS t3 ON `workbook`.`id`= t3."+objectColumnName, subQuery).
		Order("`workbook`.`name`").Limit(limit).Offset(offset).
		Scan(&workbookEntities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.WorkbookModel, len(workbookEntities))
	priv := userD.NewPrivileges([]userD.RBACAction{domain.PrivilegeRead})
	for i, e := range workbookEntities {
		problemType, err := r.problemTypes.ToProblemType(e.ProblemTypeID)
		if err != nil {
			return nil, liberrors.Errorf("r.problemTypes.ToProblemType. err: %w", err)
		}
		w, err := e.toWorkbookModel(r.rf, r.pf, operator, problemType, priv)
		if err != nil {
			return nil, liberrors.Errorf("toWorkbookModel. err: %w", err)
		}
		results[i] = w
	}

	var count int64
	rows, err := r.db.Raw("select count(*) from workbook inner join (?) AS t3 ON `workbook`.`id`= t3."+objectColumnName, subQuery).Rows()
	if err != nil {
		return nil, liberrors.Errorf("r.db.Raw. err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c int64
		if err := rows.Scan(&c); err != nil {
			return nil, liberrors.Errorf("rows.Scan. err: %w", err)
		}
		count += c
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	workbooks, err := service.NewWorkbookSearchResult(int(count), results)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return workbooks, nil
}

func (r *workbookRepository) getAllWorkbookRoles(workbookID domain.WorkbookID) []userD.RBACRole {
	return []userD.RBACRole{domain.NewWorkbookWriter(workbookID), domain.NewWorkbookReader(workbookID)}
}

func (r *workbookRepository) getAllWorkbookPrivileges() []userD.RBACAction {
	return []userD.RBACAction{domain.PrivilegeRead, domain.PrivilegeUpdate, domain.PrivilegeRemove}
}

func (r *workbookRepository) checkPrivileges(e *casbin.Enforcer, userObject userD.RBACUser, workbookObject userD.RBACObject, privs []userD.RBACAction) (userD.Privileges, error) {
	actions := make([]userD.RBACAction, 0)
	for _, priv := range privs {
		ok, err := e.Enforce(string(userObject), string(workbookObject), string(priv))
		if err != nil {
			return nil, liberrors.Errorf("e.Enforce. err: %w", err)
		}
		if ok {
			actions = append(actions, priv)
		}
	}
	return userD.NewPrivileges(actions), nil
}

// func (r *workbookRepository) canReadWorkbook(operator userD.AppUser, workbookID domain.WorkbookID) error {
// 	objectColumnName := "name"
// 	object := domain.WorkbookObjectPrefix + strconv.Itoa(int(uint(workbookID)))
// 	subject := "user_" + strconv.Itoa(int(operator.GetID()))
// 	casbinQuery, err := casbinquery.FindObject(r.db, r.driverName, object, objectColumnName, subject, "read")
// 	if err != nil {
// 		return err
// 	}
// 	var name string
// 	if result := casbinQuery.First(&name); result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return domain.ErrWorkbookPermissionDenied
// 		}
// 		return result.Error
// 	}
// 	return nil
// }

func (r *workbookRepository) FindWorkbookByID(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (service.Workbook, error) {
	ctx, span := tracer.Start(ctx, "workbookRepository.FindWorkbookByID")
	defer span.End()

	workbookEntity := workbookEntity{}
	if result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("id = ?", uint(workbookID)).
		First(&workbookEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrWorkbookNotFound
		}
		return nil, result.Error
	}

	priv, err := r.getPrivileges(ctx, operator, domain.WorkbookID(workbookEntity.ID))
	if err != nil {
		return nil, liberrors.Errorf("getPrivileges. err: %w", err)
	}
	if !priv.HasPrivilege(domain.PrivilegeRead) {
		return nil, liberrors.Errorf("AppUser(%d) has not privilege(read). err: %w", uint(operator.GetOrganizationID()), service.ErrWorkbookPermissionDenied)
	}

	logger := log.FromContext(ctx)
	logger.Infof("ownerId: %d, operatorId: %d", workbookEntity.OwnerID, operator.GetID())

	problemType, err := r.problemTypes.ToProblemType(workbookEntity.ProblemTypeID)
	if err != nil {
		return nil, liberrors.Errorf("r.problemTypes.ToProblemType. err: %w", err)
	}

	workbookModel, err := workbookEntity.toWorkbookModel(r.rf, r.pf, operator, problemType, priv)
	if err != nil {
		return nil, liberrors.Errorf("workbookEntity.toWorkbookModel. err: %w", err)
	}

	workbook, err := service.NewWorkbook(ctx, r.rf, r.pf, workbookModel)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return workbook, nil
}

func (r *workbookRepository) FindWorkbookByName(ctx context.Context, operator userD.AppUserModel, spaceID userD.SpaceID, name string) (service.Workbook, error) {
	ctx, span := tracer.Start(ctx, "workbookRepository.FindWorkbookByName")
	defer span.End()

	workbookEntity := workbookEntity{}
	if result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("space_id = ?", uint(spaceID)).
		Where("name = ?", name).
		First(&workbookEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrWorkbookNotFound
		}
		return nil, result.Error
	}

	var priv userD.Privileges
	if spaceID == service.GetSystemSpaceID() {
		priv = userD.NewPrivileges([]userD.RBACAction{domain.PrivilegeRead})
	} else {
		privTmp, err := r.getPrivileges(ctx, operator, domain.WorkbookID(workbookEntity.ID))
		if err != nil {
			return nil, liberrors.Errorf("failed to checkPrivileges. err: %w", err)
		}
		if !privTmp.HasPrivilege(domain.PrivilegeRead) {
			return nil, service.ErrWorkbookPermissionDenied
		}
		priv = privTmp
	}

	logger := log.FromContext(ctx)
	logger.Infof("ownerId: %d, operatorId: %d", workbookEntity.OwnerID, operator.GetID())

	problemType, err := r.problemTypes.ToProblemType(workbookEntity.ProblemTypeID)
	if err != nil {
		return nil, liberrors.Errorf("r.problemTypes.ToProblemType. err: %w", err)
	}

	workbookModel, err := workbookEntity.toWorkbookModel(r.rf, r.pf, operator, problemType, priv)
	if err != nil {
		return nil, liberrors.Errorf("workbookEntity.toWorkbookModel. err: %w", err)
	}

	workbook, err := service.NewWorkbook(ctx, r.rf, r.pf, workbookModel)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return workbook, nil
}

func (r *workbookRepository) getPrivileges(ctx context.Context, operator userD.AppUserModel, workbookID domain.WorkbookID) (userD.Privileges, error) {
	userRf, err := r.rf.NewUserRepositoryFactory(ctx)
	if err != nil {
		return nil, liberrors.Errorf("r.rf.NewUserRepositoryFactory. err: %w", err)
	}

	rbacRepo := userRf.NewRBACRepository(ctx)

	workbookRoles := r.getAllWorkbookRoles(workbookID)
	userObject := userD.NewUserObject(userD.AppUserID(operator.GetID()))
	e, err := rbacRepo.NewEnforcerWithRolesAndUsers(workbookRoles, []userD.RBACUser{userObject})
	if err != nil {
		return nil, liberrors.Errorf("failed to NewEnforcerWithRolesAndUsers. err: %w", err)
	}
	workbookObject := domain.NewWorkbookObject(workbookID)
	privs := r.getAllWorkbookPrivileges()
	return r.checkPrivileges(e, userObject, workbookObject, privs)
}

func (r *workbookRepository) AddWorkbook(ctx context.Context, operator userD.AppUserModel, spaceID userD.SpaceID, param service.WorkbookAddParameter) (domain.WorkbookID, error) {
	_, span := tracer.Start(ctx, "workbookRepository.AddWorkbook")
	defer span.End()

	problemTypeID, err := r.problemTypes.ToProblemTypeID(param.GetProblemType())
	if err != nil {
		return 0, liberrors.Errorf("unsupported problemType. problemType: %s", param.GetProblemType())
	}
	propertiesJSON, err := stringMapToJSON(param.GetProperties())
	if err != nil {
		return 0, liberrors.Errorf("stringMapToJSON. err: %w", err)
	}
	workbook := workbookEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		SpaceID:        uint(spaceID),
		OwnerID:        operator.GetID(),
		ProblemTypeID:  problemTypeID,
		Name:           param.GetName(),
		Lang2:          param.GetLang2().String(),
		QuestionText:   param.GetQuestionText(),
		Properties:     propertiesJSON,
	}
	if result := r.db.Create(&workbook); result.Error != nil {
		return 0, liberrors.Errorf(". err: %w", libG.ConvertDuplicatedError(result.Error, service.ErrWorkbookAlreadyExists))
	}

	workbookID := domain.WorkbookID(workbook.ID)

	userRf, err := r.rf.NewUserRepositoryFactory(ctx)
	if err != nil {
		return 0, liberrors.Errorf("r.rf.NewUserRepositoryFactory. err: %w", err)
	}

	rbacRepo := userRf.NewRBACRepository(ctx)
	userObject := userD.NewUserObject(userD.AppUserID(operator.GetID()))
	workbookObject := domain.NewWorkbookObject(workbookID)
	workbookWriter := domain.NewWorkbookWriter(workbookID)

	// the workbookWriter role can read, update, remove
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeRead); err != nil {
		return 0, liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeUpdate); err != nil {
		return 0, liberrors.Errorf("Failed to AddNamedPolicy. priv: update, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeRemove); err != nil {
		return 0, liberrors.Errorf("Failed to AddNamedPolicy. priv: remove, err: %w", err)
	}

	// user is assigned the workbookWriter role
	if err := rbacRepo.AddNamedGroupingPolicy(userObject, workbookWriter); err != nil {
		return 0, liberrors.Errorf("Failed to AddNamedGroupingPolicy. err: %w", err)
	}

	// rbacRepo.NewEnforcerWithRolesAndUsers([]userD.RBACRole{workbookWriter}, []userD.RBACUser{userObject})

	return workbookID, nil
}

func (r *workbookRepository) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID, version int) error {
	_, span := tracer.Start(ctx, "workbookRepository.RemoveWorkbook")
	defer span.End()

	workbook := workbookEntity{}
	if result := r.db.Where("organization_id = ? and id = ? and version = ?", operator.GetOrganizationID(), id, version).Delete(&workbook); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return service.ErrWorkbookNotFound
		}

		return result.Error
	}

	return nil
}

func (r *workbookRepository) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID, version int, param service.WorkbookUpdateParameter) error {
	_, span := tracer.Start(ctx, "workbookRepository.UpdateWorkbook")
	defer span.End()

	if result := r.db.Model(&workbookEntity{}).
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("id = ?", uint(id)).
		Where("version = ?", version).
		Updates(map[string]interface{}{
			"name":          param.GetName(),
			"question_text": param.GetQuestionText(),
			"version":       gorm.Expr("version + 1"),
		}); result.Error != nil {
		return liberrors.Errorf(". err: %w", libG.ConvertDuplicatedError(result.Error, service.ErrWorkbookAlreadyExists))
	}

	return nil
}

// func (r *workbookRepository) ChangeSpace(ctx context.Context, operator domain.AbstractStudent, id uint, spaceID uint) error {
// 	result := r.db.Model(&workbookEntity{}).Where(workbookEntity{
// 		OrganizationID: operator.OrganizationID(),
// 		ID:             id,
// 	}).Update(workbookEntity{
// 		SpaceID: spaceID,
// 	})
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return domain.NewWorkbookNotFoundError(id)
// 	}

// 	return nil
// }
