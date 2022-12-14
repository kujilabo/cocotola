package service

import (
	"context"
	"errors"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/lib/log"
)

type ProblemEventType int

const ProblemEventTypeAdd ProblemEventType = 1
const ProblemEventTypeUpdate ProblemEventType = 2
const ProblemEventTypeRemove ProblemEventType = 3
const observerTimeoutSec = 2

type ProblemEvent interface {
	GetOrganizationID() userD.OrganizationID
	GetAppUserID() userD.AppUserID
	GetProblemEventType() ProblemEventType
	GetProblemType() domain.ProblemTypeName
	GetProblemIDs() []domain.ProblemID
}

type problemEvent struct {
	OrganizationID   userD.OrganizationID
	AppUserID        userD.AppUserID
	ProblemEventType ProblemEventType
	ProblemType      domain.ProblemTypeName
	ProblemIDs       []domain.ProblemID
}

func NewProblemEvent(organizationID userD.OrganizationID, appUserID userD.AppUserID, problemEventType ProblemEventType, problemType domain.ProblemTypeName, problemIDs []domain.ProblemID) ProblemEvent {
	return &problemEvent{
		OrganizationID:   organizationID,
		AppUserID:        appUserID,
		ProblemEventType: problemEventType,
		ProblemType:      problemType,
		ProblemIDs:       problemIDs,
	}
}

func (p *problemEvent) GetOrganizationID() userD.OrganizationID {
	return p.OrganizationID
}

func (p *problemEvent) GetAppUserID() userD.AppUserID {
	return p.AppUserID
}

func (p *problemEvent) GetProblemEventType() ProblemEventType {
	return p.ProblemEventType
}
func (p *problemEvent) GetProblemType() domain.ProblemTypeName {
	return p.ProblemType
}
func (p *problemEvent) GetProblemIDs() []domain.ProblemID {
	return p.ProblemIDs
}

type ProblemObserver interface {
	Update(ctx context.Context, problemNotification ProblemEvent) error
}

type ProblemMonitor interface {
	Attach(observer ProblemObserver) error
	Detach(observer ProblemObserver) error
	NotifyObservers(ctx context.Context, event ProblemEvent) error
}

type problemMonitor struct {
	observers []ProblemObserver
}

func NewProblemMonitor() ProblemMonitor {
	return &problemMonitor{}
}

func (m *problemMonitor) Attach(observer ProblemObserver) error {
	for _, o := range m.observers {
		if o == observer {
			return errors.New("observer already exists")
		}
	}

	m.observers = append(m.observers, observer)
	return nil
}

func (m *problemMonitor) Detach(observer ProblemObserver) error {
	for i, o := range m.observers {
		if o == observer {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			return nil
		}
	}
	return errors.New("observer not found")
}

func (m *problemMonitor) NotifyObservers(ctx context.Context, event ProblemEvent) error {
	for _, o := range m.observers {
		go func(ctx context.Context, o ProblemObserver) {
			ctx, cancel := context.WithTimeout(ctx, observerTimeoutSec*time.Second)
			logger := log.FromContext(ctx)
			defer cancel()
			if err := o.Update(ctx, event); err != nil {
				logger.Errorf("err: %v", err)
			}
		}(context.Background(), o)
	}
	return nil
}
