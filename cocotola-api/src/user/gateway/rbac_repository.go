package gateway

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

const conf = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

type rbacRepository struct {
	db *gorm.DB
}

func NewRBACRepository(ctx context.Context, db *gorm.DB) service.RBACRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	return &rbacRepository{
		db: db,
	}
}

func (r *rbacRepository) Init() error {
	a, err := gormadapter.NewAdapterByDB(r.db)
	if err != nil {
		return liberrors.Errorf("gormadapter.NewAdapterByDB. err: %w", err)
	}

	m, err := model.NewModelFromString(conf)
	if err != nil {
		return liberrors.Errorf("model.NewModelFromString. err: %w", err)
	}

	if err := a.SavePolicy(m); err != nil {
		return liberrors.Errorf(". err: %w", err)
	}

	return nil
}

func (r *rbacRepository) initEnforcer() (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(r.db)
	if err != nil {
		return nil, liberrors.Errorf("gormadapter.NewAdapterByDB. err: %w", err)
	}

	m, err := model.NewModelFromString(conf)
	if err != nil {
		return nil, liberrors.Errorf("model.NewModelFromString. err: %w", err)
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, liberrors.Errorf("casbin.NewEnforcer. err: %w", err)
	}

	return e, nil
}

func (r *rbacRepository) AddNamedPolicy(subject domain.RBACRole, object domain.RBACObject, action domain.RBACAction) error {
	e, err := r.initEnforcer()
	if err != nil {
		return liberrors.Errorf("r.initEnforcer. err: %w", err)
	}

	if _, err := e.AddNamedPolicy("p", string(subject), string(object), string(action)); err != nil {
		return liberrors.Errorf("e.AddNamedPolicy. err: %w", err)
	}

	return nil
}

func (r *rbacRepository) AddNamedGroupingPolicy(subject domain.RBACUser, object domain.RBACRole) error {
	e, err := r.initEnforcer()
	if err != nil {
		return liberrors.Errorf("r.initEnforcer. err: %w", err)
	}
	if e == nil {
		return errors.Errorf("Nil")
	}

	if _, err := e.AddNamedGroupingPolicy("g", string(subject), string(object)); err != nil {
		return liberrors.Errorf("e.AddNamedGroupingPolicy. err: %w", err)
	}

	return nil
}

func (r *rbacRepository) NewEnforcerWithRolesAndUsers(roles []domain.RBACRole, users []domain.RBACUser) (*casbin.Enforcer, error) {
	subjects := make([]string, 0)
	for _, s := range roles {
		subjects = append(subjects, string(s))
	}
	for _, s := range users {
		subjects = append(subjects, string(s))
	}
	e, err := r.initEnforcer()
	if err != nil {
		return nil, liberrors.Errorf("r.initEnforcer. err: %w", err)
	}
	if err := e.LoadFilteredPolicy(gormadapter.Filter{V0: subjects}); err != nil {
		return nil, liberrors.Errorf("e.LoadFilteredPolicy. err: %w", err)
	}
	return e, nil
}
