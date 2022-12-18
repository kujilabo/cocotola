package service

import (
	"context"
	"errors"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/lib/log"
)

type StudyEventType int

const StudyEventTypeAnswer StudyEventType = 1

type StudyEvent interface {
	GetOrganizationID() userD.OrganizationID
	GetAppUserID() userD.AppUserID
	GetStudyEventType() StudyEventType
	GetProblemType() domain.ProblemTypeName
	GetStudyType() domain.StudyTypeName
	GetProblemID() domain.ProblemID
}

type studyEvent struct {
	OrganizationID userD.OrganizationID
	AppUserID      userD.AppUserID
	StudyEventType StudyEventType
	ProblemType    domain.ProblemTypeName
	StudyType      domain.StudyTypeName
	ProblemID      domain.ProblemID
}

func NewStudyEvent(organizationID userD.OrganizationID, appUserID userD.AppUserID, studyEventType StudyEventType, problemType domain.ProblemTypeName, studyType domain.StudyTypeName, problemID domain.ProblemID) StudyEvent {
	return &studyEvent{
		OrganizationID: organizationID,
		AppUserID:      appUserID,
		StudyEventType: studyEventType,
		ProblemType:    problemType,
		StudyType:      studyType,
		ProblemID:      problemID,
	}
}

func (e *studyEvent) GetOrganizationID() userD.OrganizationID {
	return e.OrganizationID
}

func (e *studyEvent) GetAppUserID() userD.AppUserID {
	return e.AppUserID
}

func (e *studyEvent) GetStudyEventType() StudyEventType {
	return e.StudyEventType
}

func (e *studyEvent) GetProblemType() domain.ProblemTypeName {
	return e.ProblemType
}

func (e *studyEvent) GetStudyType() domain.StudyTypeName {
	return e.StudyType
}

func (e *studyEvent) GetProblemID() domain.ProblemID {
	return e.ProblemID
}

type StudyObserver interface {
	Update(ctx context.Context, studyNotification StudyEvent) error
}

type StudyMonitor interface {
	Attach(observer StudyObserver) error
	Detach(observer StudyObserver) error
	NotifyObservers(ctx context.Context, event StudyEvent) error
}

type studyMonitor struct {
	observers []StudyObserver
}

func NewStudyMonitor() StudyMonitor {
	return &studyMonitor{
		observers: make([]StudyObserver, 0),
	}
}

func (m *studyMonitor) Attach(observer StudyObserver) error {
	for _, o := range m.observers {
		if o == observer {
			return errors.New("observer already exists")
		}
	}
	m.observers = append(m.observers, observer)
	return nil
}

func (m *studyMonitor) Detach(observer StudyObserver) error {
	for i, o := range m.observers {
		if o == observer {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			return nil
		}
	}
	return errors.New("observer not found")
}

func (m *studyMonitor) NotifyObservers(ctx context.Context, event StudyEvent) error {
	for _, o := range m.observers {
		go func(ctx context.Context, o StudyObserver) {
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
