package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/config"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	usecaseJ "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/job"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	"github.com/kujilabo/cocotola/lib/controller/middleware"
	ginlog "github.com/onrik/logrus/gin"
)

func NewJobRouter(transaction service.Transaction, jobService jobS.JobService, debugConfig *config.DebugConfig) (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	jobUseCaseStat := usecaseJ.NewJobUsecaseStat(transaction, jobService)
	jobHandler := NewJobHandler(jobUseCaseStat)
	router.GET("aggregate_results", jobHandler.AggregateStudyResultsOfAllUsers)
	return router, nil
}
