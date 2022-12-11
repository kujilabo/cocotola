package controller_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/controller"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	service_mock "github.com/kujilabo/cocotola/cocotola-api/src/app/service/mock"
	studentU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student"
	studentU_mock "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student/mock"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func parseJSON(t *testing.T, b *bytes.Buffer) interface{} {
	respBytes, err := io.ReadAll(b)
	require.NoError(t, err)
	obj, err := oj.Parse(respBytes)
	require.NoError(t, err)
	return obj
}

func parseExpr(t *testing.T, v string) jp.Expr {
	expr, err := jp.ParseString(v)
	require.NoError(t, err)
	return expr
}

func newAuthMiddleware(organizationID userD.OrganizationID, userID userD.AppUserID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("OrganizationID", int(organizationID))
		c.Set("AuthorizedUser", int(userID))
	}
}

func initAudioRouter(studentUsecaseAudio studentU.StudentUsecaseAudio, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(authMiddleware)
	g := router.Group("v1")
	fn := controller.NewInitAudioRouterFunc(studentUsecaseAudio)
	fn(g)
	return router
}

func Test_FindAudioByID_OK(t *testing.T) {
	logrus.SetLevel(logrus.WarnLevel)

	// given
	audioModel, _ := domain.NewAudioModel(1, domain.Lang2EN, "apple", "APPLE")
	audio := new(service_mock.Audio)
	audio.On("GetAudioModel").Return(audioModel)

	authMiddleware := newAuthMiddleware(userD.OrganizationID(1), userD.AppUserID(2))
	studentUsecaseAudio := new(studentU_mock.StudentUsecaseAudio)
	studentUsecaseAudio.On("FindAudioByID", anythingOfContext, userD.OrganizationID(1), userD.AppUserID(2), domain.WorkbookID(3), domain.ProblemID(4), domain.AudioID(5)).Return(audio, nil)

	r := initAudioRouter(studentUsecaseAudio, authMiddleware)
	w := httptest.NewRecorder()

	// when
	path := fmt.Sprintf("/v1/workbook/%d/problem/%d/audio/%d", 3, 4, 5)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)

	// then
	idExpr := parseExpr(t, "$.id")
	lang2Expr := parseExpr(t, "$.lang2")
	textExpr := parseExpr(t, "$.text")
	contentExpr := parseExpr(t, "$.content")

	// - check the status code
	assert.Equal(t, http.StatusOK, w.Code)
	jsonObj := parseJSON(t, w.Body)

	id := idExpr.Get(jsonObj)
	assert.Equal(t, int64(1), id[0].(int64))

	lang2 := lang2Expr.Get(jsonObj)
	assert.Equal(t, "en", lang2[0].(string))

	text := textExpr.Get(jsonObj)
	assert.Equal(t, "apple", text[0].(string))

	content := contentExpr.Get(jsonObj)
	assert.Equal(t, "APPLE", content[0].(string))
}
