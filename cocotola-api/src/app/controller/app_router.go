package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/config"
	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student"
	authH "github.com/kujilabo/cocotola/cocotola-api/src/auth/controller"
	authM "github.com/kujilabo/cocotola/cocotola-api/src/auth/controller/middleware"
	"github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	authU "github.com/kujilabo/cocotola/cocotola-api/src/auth/usecase"
	pluginCommonController "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/controller"
	pluginCommonService "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	"github.com/kujilabo/cocotola/lib/controller/middleware"
)

type NewIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appS.ProblemAddParameterIterator, error)

type InitRouterGroupFunc func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error

func NewInitAuthRouterFunc(googleUserUsecase authU.GoogleUserUsecase, guestUserUsecase authU.GuestUserUsecase, authTokenManager service.AuthTokenManager) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		v1auth := parentRouterGroup.Group("auth")
		for _, m := range middleware {
			v1auth.Use(m)
		}
		// googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
		// guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
		authHandler := authH.NewAuthHandler(authTokenManager)
		googleAuthHandler := authH.NewGoogleAuthHandler(googleUserUsecase)
		guestAuthHandler := authH.NewGuestAuthHandler(guestUserUsecase)
		v1auth.POST("google/authorize", googleAuthHandler.Authorize)
		v1auth.POST("guest/authorize", guestAuthHandler.Authorize)
		v1auth.POST("refresh_token", authHandler.RefreshToken)
		return nil
	}
}

func NewInitWorkbookRouterFunc(studentUsecaseWorkbook studentU.StudentUsecaseWorkbook) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		v1Workbook := parentRouterGroup.Group("private/workbook")
		privateWorkbookHandler := NewPrivateWorkbookHandler(studentUsecaseWorkbook)
		for _, m := range middleware {
			v1Workbook.Use(m)
		}
		v1Workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		v1Workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		v1Workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		v1Workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		v1Workbook.POST("", privateWorkbookHandler.AddWorkbook)
		return nil
	}
}

func NewInitProblemRouterFunc(studentUsecaseProblem studentU.StudentUsecaseProblem, newIteratorFunc NewIteratorFunc) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {

		v1Problem := parentRouterGroup.Group("workbook/:workbookID/problem")
		problemHandler, err := NewProblemHandler(studentUsecaseProblem, newIteratorFunc)
		if err != nil {
			return err
		}
		for _, m := range middleware {
			v1Problem.Use(m)
		}
		v1Problem.POST("", problemHandler.AddProblem)
		v1Problem.GET(":problemID", problemHandler.FindProblemByID)
		v1Problem.DELETE(":problemID", problemHandler.RemoveProblem)
		v1Problem.PUT(":problemID", problemHandler.UpdateProblem)
		v1Problem.PUT(":problemID/property", problemHandler.UpdateProblemProperty)
		// v1Problem.GET("problem_ids", problemHandler.FindProblemIDs)
		v1Problem.POST("find", problemHandler.FindProblems)
		v1Problem.POST("find_all", problemHandler.FindAllProblems)
		v1Problem.POST("find_by_ids", problemHandler.FindProblemsByProblemIDs)
		v1Problem.POST("import", problemHandler.ImportProblems)

		return nil
	}
}
func NewInitStudyRouterFunc(studentUsecaseStudy studentU.StudentUsecaseStudy) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		v1Study := parentRouterGroup.Group("study/workbook/:workbookID")
		recordbookHandler, err := NewRecordbookHandler(studentUsecaseStudy)
		if err != nil {
			return err
		}
		for _, m := range middleware {
			v1Study.Use(m)
		}
		v1Study.GET("study_type/:studyType", recordbookHandler.FindRecordbook)
		v1Study.POST("study_type/:studyType/problem/:problemID/record", recordbookHandler.SetStudyResult)
		v1Study.GET("completion_rate", recordbookHandler.GetCompletionRate)

		return nil
	}
}
func NewInitAudioRouterFunc(studentUsecaseAudio studentU.StudentUsecaseAudio) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		v1Audio := parentRouterGroup.Group("workbook/:workbookID/problem/:problemID/audio")
		audioHandler := NewAudioHandler(studentUsecaseAudio)
		for _, m := range middleware {
			v1Audio.Use(m)
		}
		v1Audio.GET(":audioID", audioHandler.FindAudioByID)
		return nil
	}
}

func NewInitStatRouterFunc(studenUsecaseStat studentU.StudentUsecaseStat) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		v1Stat := parentRouterGroup.Group("stat")
		for _, m := range middleware {
			v1Stat.Use(m)
		}
		statHandler := NewStatHandler(studenUsecaseStat)
		v1Stat.GET("", statHandler.GetStat)
		return nil
	}
}
func NewInitTranslationRouterFunc(translatorClient pluginCommonService.TranslatorClient) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		pluginTranslation := parentRouterGroup.Group("translation")
		for _, m := range middleware {
			pluginTranslation.Use(m)
		}
		translationHandler := pluginCommonController.NewTranslationHandler(translatorClient)
		pluginTranslation.POST("find", translationHandler.FindTranslations)
		pluginTranslation.GET("text/:text/pos/:pos", translationHandler.FindTranslationByTextAndPos)
		pluginTranslation.GET("text/:text", translationHandler.FindTranslationsByText)
		pluginTranslation.PUT("text/:text/pos/:pos", translationHandler.UpdateTranslation)
		pluginTranslation.DELETE("text/:text/pos/:pos", translationHandler.RemoveTranslation)
		pluginTranslation.POST("", translationHandler.AddTranslation)
		pluginTranslation.POST("export", translationHandler.ExportTranslations)
		return nil
	}
}

func NewInitTatoebaRouterFunc(tatoebaClient pluginCommonService.TatoebaClient) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		pluginTatoeba := parentRouterGroup.Group("tatoeba")
		for _, m := range middleware {
			pluginTatoeba.Use(m)
		}
		tatoebaHandler := pluginCommonController.NewTatoebaHandler(tatoebaClient)
		pluginTatoeba.POST("find", tatoebaHandler.FindSentencePairs)
		pluginTatoeba.POST("sentence/import", tatoebaHandler.ImportSentences)
		pluginTatoeba.POST("link/import", tatoebaHandler.ImportLinks)
		return nil
	}
}

func NewAppRouter(initPublicRouterFunc []InitRouterGroupFunc, initPrivateRouterFunc []InitRouterGroupFunc, initPluginRouterFunc []InitRouterGroupFunc, authTokenManager service.AuthTokenManager, corsConfig cors.Config, appConfig *config.AppConfig, authConfig *config.AuthConfig, debugConfig *config.DebugConfig) (*gin.Engine, error) {
	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	signingKey := []byte(authConfig.SigningKey)
	authMiddleware := authM.NewAuthMiddleware(signingKey)

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))

		for _, fn := range initPublicRouterFunc {
			if err := fn(v1); err != nil {
				return nil, err
			}
		}

		for _, fn := range initPrivateRouterFunc {
			if err := fn(v1, authMiddleware); err != nil {
				return nil, err
			}
		}
	}

	plugin := router.Group("plugin")
	{
		plugin.Use(otelgin.Middleware(appConfig.Name))
		plugin.Use(middleware.NewTraceLogMiddleware(appConfig.Name))
		plugin.Use(authMiddleware)

		for _, fn := range initPluginRouterFunc {
			if err := fn(v1); err != nil {
				return nil, err
			}
		}
	}

	return router, nil
}
