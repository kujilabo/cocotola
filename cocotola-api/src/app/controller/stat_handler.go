package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student"
	controllerhelper "github.com/kujilabo/cocotola/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type StatHandler interface {
	GetStat(c *gin.Context)
}

type statHandler struct {
	studentUsecaseStat studentU.StudentUsecaseStat
}

func NewStatHandler(studentUsecaseStat studentU.StudentUsecaseStat) StatHandler {
	return &statHandler{
		studentUsecaseStat: studentUsecaseStat,
	}
}

type Result struct {
	Date     string `json:"date"`
	Mastered int    `json:"mastered"`
	Answered int    `json:"answered"`
}

type History struct {
	Results []Result `json:"results"`
}

type Stat struct {
	History History `json:"history"`
}

func (h *statHandler) GetStat(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("GetStat")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		stat, err := h.studentUsecaseStat.GetStat(ctx, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf(". err: %w", err)
		}

		results := make([]Result, 0)
		for _, result := range stat.GetHistory().Results {
			results = append(results, Result{
				Date:     result.Date.Format("2006-01-02"),
				Mastered: result.Mastered,
				Answered: result.Answered,
			})
		}
		response := Stat{
			History: History{
				Results: results,
			},
		}
		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *statHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if errors.Is(err, service.ErrProblemAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"message": "Problem already exists"})
		return true
	} else if errors.Is(err, service.ErrWorkbookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	}
	logger.Errorf("studyHandler error:%v", err)
	return false
}
