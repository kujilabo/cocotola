package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/controller/converter"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student"
	controllerhelper "github.com/kujilabo/cocotola/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/lib/controller/helper"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type AudioHandler interface {
	FindAudioByID(c *gin.Context)
}

type audioHandler struct {
	studentUsecaseAudio studentU.StudentUsecaseAudio
}

func NewAudioHandler(studentUsecaseAudio studentU.StudentUsecaseAudio) AudioHandler {
	return &audioHandler{
		studentUsecaseAudio: studentUsecaseAudio,
	}
}

func (h *audioHandler) FindAudioByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindAudioByID")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := helper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		problemID, err := helper.GetUintFromPath(c, "problemID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		audioID, err := helper.GetUintFromPath(c, "audioID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.studentUsecaseAudio.FindAudioByID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), domain.ProblemID(problemID), domain.AudioID(audioID))
		if err != nil {
			return liberrors.Errorf("domain.AudioID. err: %w", err)
		}

		response, err := converter.ToAudioResponse(ctx, result)
		if err != nil {
			return liberrors.Errorf("converter.ToAudioResponse. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *audioHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if errors.Is(err, service.ErrAudioNotFound) {
		logger.Warnf("audioHandler err: %+v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
		return true
	}
	logger.Errorf("error:%v", err)
	return false
}
