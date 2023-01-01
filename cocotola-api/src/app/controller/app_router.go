package controller

import (
	"context"
	"io"

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

type NewIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType appD.ProblemTypeName, reader io.Reader) (appS.ProblemAddParameterIterator, error)

type InitRouterGroupFunc func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error

func NewInitAuthRouterFunc(googleUserUsecase authU.GoogleUserUsecase, guestUserUsecase authU.GuestUserUsecase, authTokenManager service.AuthTokenManager) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		auth := parentRouterGroup.Group("auth")
		for _, m := range middleware {
			auth.Use(m)
		}
		// googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
		// guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
		authHandler := authH.NewAuthHandler(authTokenManager)
		googleAuthHandler := authH.NewGoogleAuthHandler(googleUserUsecase)
		guestAuthHandler := authH.NewGuestAuthHandler(guestUserUsecase)
		auth.POST("google/authorize", googleAuthHandler.Authorize)
		auth.POST("guest/authorize", guestAuthHandler.Authorize)
		auth.POST("refresh_token", authHandler.RefreshToken)
		return nil
	}
}

func NewInitWorkbookRouterFunc(studentUsecaseWorkbook studentU.StudentUsecaseWorkbook) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		workbook := parentRouterGroup.Group("private/workbook")
		privateWorkbookHandler := NewPrivateWorkbookHandler(studentUsecaseWorkbook)
		for _, m := range middleware {
			workbook.Use(m)
		}
		workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		workbook.POST("", privateWorkbookHandler.AddWorkbook)
		return nil
	}
}

func NewInitProblemRouterFunc(studentUsecaseProblem studentU.StudentUsecaseProblem, newIteratorFunc NewIteratorFunc) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {

		problem := parentRouterGroup.Group("workbook/:workbookID/problem")
		problemHandler, err := NewProblemHandler(studentUsecaseProblem, newIteratorFunc)
		if err != nil {
			return err
		}
		for _, m := range middleware {
			problem.Use(m)
		}
		problem.POST("", problemHandler.AddProblem)
		problem.GET(":problemID", problemHandler.FindProblemByID)
		problem.DELETE(":problemID", problemHandler.RemoveProblem)
		problem.PUT(":problemID", problemHandler.UpdateProblem)
		problem.PUT(":problemID/property", problemHandler.UpdateProblemProperty)
		// v1Problem.GET("problem_ids", problemHandler.FindProblemIDs)
		problem.POST("find", problemHandler.FindProblems)
		problem.POST("find_all", problemHandler.FindAllProblems)
		problem.POST("find_by_ids", problemHandler.FindProblemsByProblemIDs)
		problem.POST("import", problemHandler.ImportProblems)

		return nil
	}
}
func NewInitStudyRouterFunc(studentUsecaseStudy studentU.StudentUsecaseStudy) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		study := parentRouterGroup.Group("study/workbook/:workbookID")
		recordbookHandler, err := NewRecordbookHandler(studentUsecaseStudy)
		if err != nil {
			return err
		}
		for _, m := range middleware {
			study.Use(m)
		}
		study.GET("study_type/:studyType", recordbookHandler.FindRecordbook)
		study.POST("study_type/:studyType/problem/:problemID/record", recordbookHandler.SetStudyResult)
		study.GET("completion_rate", recordbookHandler.GetCompletionRate)

		return nil
	}
}

func NewInitAudioRouterFunc(studentUsecaseAudio studentU.StudentUsecaseAudio) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		audio := parentRouterGroup.Group("workbook/:workbookID/problem/:problemID/audio")
		audioHandler := NewAudioHandler(studentUsecaseAudio)
		for _, m := range middleware {
			audio.Use(m)
		}
		audio.GET(":audioID", audioHandler.FindAudioByID)
		return nil
	}
}

func NewInitStatRouterFunc(studenUsecaseStat studentU.StudentUsecaseStat) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		stat := parentRouterGroup.Group("stat")
		for _, m := range middleware {
			stat.Use(m)
		}
		statHandler := NewStatHandler(studenUsecaseStat)
		stat.GET("", statHandler.GetStat)
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
			if err := fn(plugin); err != nil {
				return nil, err
			}
		}
	}

	return router, nil
}
