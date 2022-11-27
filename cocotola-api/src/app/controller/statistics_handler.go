package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	controllerhelper "github.com/kujilabo/cocotola/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/lib/log"
)

type StatHandler interface {
	GetStat(c *gin.Context)
}

type statHandler struct {
}

func NewStatHandler() StatHandler {
	return &statHandler{}
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
		response := Stat{
			History: History{
				Results: []Result{
					{
						Date:     "2022-11-01",
						Mastered: 10,
						Answered: 20,
					},
					{
						Date:     "2022-11-02",
						Mastered: 10,
						Answered: 20,
					},
				},
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
