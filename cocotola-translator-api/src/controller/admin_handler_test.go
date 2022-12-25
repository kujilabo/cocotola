package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/controller"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/usecase"
	usecase_mock "github.com/kujilabo/cocotola/cocotola-translator-api/src/usecase/mock"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func initCrosConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	return corsConfig
}

func initAdminRouter(t *testing.T, adminUsecase usecase.AdminUsecase) *gin.Engine {
	router := gin.New()
	g := router.Group("v1")
	fn := controller.NewInitAdminRouterFunc(adminUsecase)
	err := fn(g)
	require.NoError(t, err)
	return router
}

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

func Test_adminHandler_FindTranslationsByFirstLetter_OK(t *testing.T) {
	// given
	adminUsecase := new(usecase_mock.AdminUsecase)

	apple, err := domain.NewTranslation(1, time.Now(), time.Now(), "apple", domain.PosNoun, domain.Lang2JA, "リンゴ", "mock")
	require.NoError(t, err)
	adminUsecase.On("FindTranslationsByFirstLetter", anythingOfContext, domain.Lang2JA, "a").Return([]domain.Translation{
		apple,
	}, nil)

	r := initAdminRouter(t, adminUsecase)
	w := httptest.NewRecorder()

	// when
	// - letter is 'a'
	body, err := json.Marshal(gin.H{"letter": "a"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/admin/find", bytes.NewBuffer(body))
	req.SetBasicAuth("user", "pass")
	require.NoError(t, err)
	r.ServeHTTP(w, req)

	// then
	resultsExpr := parseExpr(t, "$.results[*]")
	lang2Expr := parseExpr(t, "$.results[*].lang2")

	// - check the status code
	assert.Equal(t, http.StatusOK, w.Code)
	jsonObj := parseJSON(t, w.Body)

	results := resultsExpr.Get(jsonObj)
	assert.Equal(t, 1, len(results))

	lang2 := lang2Expr.Get(jsonObj)
	assert.Equal(t, "ja", lang2[0].(string))
}

func Test_adminHandler_FindTranslationsByFirstLetter_LetterIsNothing(t *testing.T) {
	// given
	adminUsecase := new(usecase_mock.AdminUsecase)

	apple, err := domain.NewTranslation(1, time.Now(), time.Now(), "apple", domain.PosNoun, domain.Lang2JA, "リンゴ", "mock")
	require.NoError(t, err)
	adminUsecase.On("FindTranslationsByFirstLetter", anythingOfContext, domain.Lang2JA, "a").Return([]domain.Translation{
		apple,
	}, nil)

	r := initAdminRouter(t, adminUsecase)

	// when
	// - letter is nothing
	body, err := json.Marshal(gin.H{})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/admin/find", bytes.NewBuffer(body))
	req.SetBasicAuth("user", "pass")
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	bytes, _ := io.ReadAll(w.Body)
	t.Logf("resp: %s", string(bytes))

	// then
	// resultsExpr := parseExpr(t, "$.results[*]")
	// lang2Expr := parseExpr(t, "$.results[*].lang2")

	// - check the status code
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// jsonObj := parseJSON(t, w.Body)

	// results := resultsExpr.Get(jsonObj)
	// assert.Equal(t, 1, len(results))

	// lang2 := lang2Expr.Get(jsonObj)
	// assert.Equal(t, "ja", lang2[0].(string))
}
