package controller

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/controller/converter"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/controller/entity"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	controllerhelper "github.com/kujilabo/cocotola/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/lib/controller/helper"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type TranslationHandler interface {
	FindTranslations(c *gin.Context)
	FindTranslationByTextAndPos(c *gin.Context)
	FindTranslationsByText(c *gin.Context)
	AddTranslation(c *gin.Context)
	UpdateTranslation(c *gin.Context)
	RemoveTranslation(c *gin.Context)
	ExportTranslations(c *gin.Context)
}

type translationHandler struct {
	translatorClient service.TranslatorClient
}

func NewTranslationHandler(translatorClient service.TranslatorClient) TranslationHandler {
	return &translationHandler{translatorClient: translatorClient}
}

func (h *translationHandler) FindTranslations(c *gin.Context) {
	ctx := c.Request.Context()

	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		param := entity.TranslationFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		lang2, err := appD.NewLang2(param.Lang2)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.translatorClient.FindTranslationsByFirstLetter(ctx, lang2, param.Letter)
		if err != nil {
			return liberrors.Errorf("h.translatorClient.FindTranslationsByFirstLetter. err: %w", err)
		}

		response, err := converter.ToTranslationFindResposne(ctx, result)
		if err != nil {
			return liberrors.Errorf("converter.ToTranslationFindResposne. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) FindTranslationByTextAndPos(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslationByTextAndPos")

	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		text := helper.GetStringFromPath(c, "text")

		pos, err := helper.GetIntFromPath(c, "pos")
		if err != nil {
			return liberrors.Errorf("helper.GetIntFromPath. err: %w", err)
		}

		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return liberrors.Errorf("domain.NewWordPos. err: %w", err)
		}

		result, err := h.translatorClient.FindTranslationByTextAndPos(ctx, appD.Lang2JA, text, wordPos)
		if err != nil {
			return liberrors.Errorf("h.translatorClient.FindTranslationByTextAndPos. err: %w", err)
		}

		response, err := converter.ToTranslationResposne(ctx, result)
		if err != nil {
			return liberrors.Errorf("converter.ToTranslationResposne. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) FindTranslationsByText(c *gin.Context) {
	ctx := c.Request.Context()

	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		text := helper.GetStringFromPath(c, "text")
		results, err := h.translatorClient.FindTranslationsByText(ctx, appD.Lang2JA, text)
		if err != nil {
			return liberrors.Errorf("h.translatorClient.FindTranslationsByText. err: %w", err)
		}

		response, err := converter.ToTranslationListResposne(ctx, results)
		if err != nil {
			return liberrors.Errorf("converter.ToTranslationListResposne. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) AddTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		param := entity.TranslationAddParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTranslationAddParameter(ctx, &param)
		if err != nil {
			return liberrors.Errorf("converter.ToTranslationAddParameter. err: %w", err)
		}

		if err := h.translatorClient.AddTranslation(ctx, parameter); err != nil {
			return liberrors.Errorf("h.translatorClient.AddTranslation. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) UpdateTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		text := helper.GetStringFromPath(c, "text")

		pos, err := helper.GetIntFromPath(c, "pos")
		if err != nil {
			return liberrors.Errorf("helper.GetIntFromPath. path: %s, err: %w", "pos", err)
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return liberrors.Errorf("domain.NewWordPos. err: %w", err)
		}

		param := entity.TranslationUpdateParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTranslationUpdateParameter(ctx, &param)
		if err != nil {
			return liberrors.Errorf("converter.ToTranslationUpdateParameter. err: %w", err)
		}

		if err := h.translatorClient.UpdateTranslation(ctx, appD.Lang2JA, text, wordPos, parameter); err != nil {
			return liberrors.Errorf("h.translatorClient.UpdateTranslation. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) RemoveTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		text := helper.GetStringFromPath(c, "text")

		pos, err := helper.GetIntFromPath(c, "pos")
		if err != nil {
			return liberrors.Errorf("helper.GetIntFromPath. path: %s, err: %w", "pos", err)
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return liberrors.Errorf("domain.NewWordPos. err: %w", err)
		}

		if err := h.translatorClient.RemoveTranslation(ctx, appD.Lang2JA, text, wordPos); err != nil {
			return liberrors.Errorf("h.translatorClient.RemoveTranslation. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) ExportTranslations(c *gin.Context) {
	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		csvStruct := [][]string{
			{"name", "address", "phone"},
			{"Ram", "Tokyo", "1236524"},
			{"Shaym", "Beijing", "8575675484"},
		}
		b := new(bytes.Buffer)
		w := csv.NewWriter(b)
		if err := w.WriteAll(csvStruct); err != nil {
			return liberrors.Errorf("w.WriteAll. err: %w", err)
		}
		if _, err := c.Writer.Write(b.Bytes()); err != nil {
			return liberrors.Errorf("c.Writer.Write. err: %w", err)
		}
		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	if errors.Is(err, service.ErrTranslationAlreadyExists) {
		logger.Warnf("translationHandler. err: %v", err)
		c.JSON(http.StatusConflict, gin.H{"message": "Translation already exists"})
		return true
	}
	logger.Errorf("translationHandler. err: %v", err)
	return false
}
