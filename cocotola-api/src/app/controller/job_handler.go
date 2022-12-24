package controller

import (
	"github.com/gin-gonic/gin"

	jobU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/job"
	controllerhelper "github.com/kujilabo/cocotola/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type JobHandler interface {
	AggregateStudyResultsOfAllUsers(c *gin.Context)
}

type jobHandler struct {
	jobUsecaseStat jobU.JobUsecaseStat
}

func NewJobHandler(jobUsecaseStat jobU.JobUsecaseStat) JobHandler {
	return &jobHandler{
		jobUsecaseStat: jobUsecaseStat,
	}
}

func (h *jobHandler) AggregateStudyResultsOfAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	systemAdminModel := userD.NewSystemAdminModel()

	// now := time.Now()
	// yesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		if err := h.jobUsecaseStat.AggregateStudyResultsOfAllUsers(ctx, systemAdminModel); err != nil {
			return liberrors.Errorf(" h.jobUsecaseStat.AggregateStudyResultsOfAllUsers. err: %w", err)
		}
		return nil
	}, h.errorHandle)
}

func (h *jobHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Errorf("jobHandler err: %+v", err)
	return false
}
