package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/config"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/usecase"
	"github.com/kujilabo/cocotola/lib/controller/middleware"
)

type InitRouterGroupFunc func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error

func NewInitAdminRouterFunc(adminUsecase usecase.AdminUsecase) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		admin := parentRouterGroup.Group("admin")
		adminHandler := NewAdminHandler(adminUsecase)
		for _, m := range middleware {
			admin.Use(m)
		}
		admin.POST("find", adminHandler.FindTranslationsByFirstLetter)
		admin.GET("text/:text/pos/:pos", adminHandler.FindTranslationByTextAndPos)
		admin.GET("text/:text", adminHandler.FindTranslationsByText)
		admin.PUT("text/:text/pos/:pos", adminHandler.UpdateTranslation)
		admin.DELETE("text/:text/pos/:pos", adminHandler.RemoveTranslation)
		admin.POST("", adminHandler.AddTranslation)
		admin.POST("export", adminHandler.ExportTranslations)
		return nil
	}
}

func NewInitUserRouterFunc(userUsecase usecase.UserUsecase) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		user := parentRouterGroup.Group("user")
		userHandler := NewUserHandler(userUsecase)
		for _, m := range middleware {
			user.Use(m)
		}
		user.GET("dictionary/lookup", userHandler.DictionaryLookup)
		return nil
	}
}

func NewAppRouter(initPublicRouterFunc []InitRouterGroupFunc, initPrivateRouterFunc []InitRouterGroupFunc, corsConfig cors.Config, appConfig *config.AppConfig, authConfig *config.AuthConfig, debugConfig *config.DebugConfig) (*gin.Engine, error) {
	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	authMiddleware := gin.BasicAuth(gin.Accounts{
		authConfig.Username: authConfig.Password,
	})

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))
		v1.Use(authMiddleware)
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

	return router, nil
}
