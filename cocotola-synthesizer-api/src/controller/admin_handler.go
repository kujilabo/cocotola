package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/usecase"
)

type AdminHandler interface {
	FindAudio(c *gin.Context)
}

type adminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) AdminHandler {
	return &adminHandler{adminUsecase: adminUsecase}
}

func (h *adminHandler) FindAudio(c *gin.Context) {

}

// func (h *adminHandler) errorHandle(c *gin.Context, err error) bool {
// 	ctx := c.Request.Context()
// 	logger := log.FromContext(ctx)

// 	logger.Errorf("adminHandler. err: %v", err)
// 	return false
// }
