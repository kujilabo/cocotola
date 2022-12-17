package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola/cocotola-api/src/app/config"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/controller"
	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appG "github.com/kujilabo/cocotola/cocotola-api/src/app/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	jobU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/job"
	studentU "github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student"
	authG "github.com/kujilabo/cocotola/cocotola-api/src/auth/gateway"
	authS "github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	authU "github.com/kujilabo/cocotola/cocotola-api/src/auth/usecase"
	english_sentence "github.com/kujilabo/cocotola/cocotola-api/src/data/english_sentence"
	english_word "github.com/kujilabo/cocotola/cocotola-api/src/data/english_word"
	"github.com/kujilabo/cocotola/cocotola-api/src/docs"
	jobG "github.com/kujilabo/cocotola/cocotola-api/src/job/gateway"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	pluginCommonGateway "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/gateway"
	pluginCommonS "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	pluginEnglishDomain "github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	pluginEnglishGateway "github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/gateway"
	pluginEnglishS "github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libconfig "github.com/kujilabo/cocotola/lib/config"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
	"github.com/kujilabo/cocotola/lib/log"
)

const readHeaderTimeout = time.Duration(30) * time.Second

const jobIntervalSec = 5

// type newIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appS.ProblemAddParameterIterator, error)

func getValue(values ...string) string {
	for _, v := range values {
		if len(v) != 0 {
			return v
		}
	}
	return ""
}

func main() {
	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	appEnv := getValue(*env, os.Getenv("APP_ENV"), "local")
	logrus.Infof("env: %s", appEnv)

	liberrors.UseXerrorsErrorf()

	cfg, db, sqlDB, tp, err := initialize(ctx, appEnv)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	if err := initApp1(ctx, db, cfg.App.OwnerPassword); err != nil {
		panic(err)
	}

	synthesizer, err := appG.NewSynthesizerClient(cfg.Synthesizer.Endpoint, cfg.Synthesizer.Username, cfg.Synthesizer.Password, time.Duration(cfg.Synthesizer.TimeoutSec)*time.Second)
	if err != nil {
		panic(err)
	}

	// translator
	connTranslator, err := grpc.Dial(cfg.Translator.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	if err != nil {
		panic(err)
	}
	defer connTranslator.Close()
	// translatorClient := pluginCommonGateway.NewTranslatorHTTPClient(cfg.Translator.Endpoint, cfg.Translator.Username, cfg.Translator.Password, time.Duration(cfg.Translator.TimeoutSec)*time.Second)
	translatorClient := pluginCommonGateway.NewTranslatorGRPCClient(connTranslator, cfg.Translator.Username, cfg.Translator.Password, time.Duration(cfg.Translator.TimeoutSec)*time.Second)

	tatoebaClient := pluginCommonGateway.NewTatoebaClient(cfg.Tatoeba.Endpoint, cfg.Tatoeba.Username, cfg.Tatoeba.Password, time.Duration(cfg.Tatoeba.TimeoutSec)*time.Second)

	pf, problemRepositories, problemImportProcessor := initPf(synthesizer, translatorClient, tatoebaClient)

	newIterator := func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
		processor, ok := problemImportProcessor[problemType]
		if ok {
			return processor.CreateCSVReader(ctx, workbookID, reader)
		}
		return nil, liberrors.Errorf("processor not found. problemType: %s", problemType)
	}

	jobRff := func(ctx context.Context, db *gorm.DB) (jobS.RepositoryFactory, error) {
		return jobG.NewRepositoryFactory(ctx, db)
	}

	userRff := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(ctx, db)
	}

	rff := func(ctx context.Context, db *gorm.DB) (appS.RepositoryFactory, error) {
		return appG.NewRepositoryFactory(ctx, db, cfg.DB.DriverName, jobRff, userRff, pf, problemRepositories)
	}

	jobTransaction, authTransaction, appTransaction, err := initTransaction(db, jobRff, userRff, rff)
	if err != nil {
		panic(err)
	}

	if err := initApp2(ctx, appTransaction); err != nil {
		panic(err)
	}

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	if *env == "local" {
		initLocalEnv(ctx, jobTransaction, appTransaction)
	}

	result := run(context.Background(), cfg, appTransaction, db, pf, rff, authTransaction, jobTransaction, appTransaction, synthesizer, translatorClient, tatoebaClient, newIterator)

	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
	os.Exit(result)
}

func initLocalEnv(ctx context.Context, jobTransaction jobS.Transaction, appTransaction appS.Transaction) {
	jobService, err := jobS.NewJobService(ctx, jobTransaction)
	if err != nil {
		panic(err)
	}

	systemAdminModel := userD.NewSystemAdminModel()

	jobUseCaseStat := jobU.NewJobUsecaseStat(appTransaction, jobService)

	s := gocron.NewScheduler(time.UTC)
	if _, err := s.Every(jobIntervalSec).Seconds().Do(func() {
		if err := jobUseCaseStat.AggregateStudyResultsOfAllUsers(context.Background(), systemAdminModel); err != nil {
			logrus.Errorf("AggregateStudyResultsOfAllUsers. err: %v", err)
		}
	}); err != nil {
		panic(err)
	}
	s.StartAsync()
}

func initTransaction(db *gorm.DB, jobRff jobG.RepositoryFactoryFunc, userRff userG.RepositoryFactoryFunc, rff appG.RepositoryFactoryFunc) (jobS.Transaction, authS.Transaction, appS.Transaction, error) {
	jobTransaction := jobG.NewTransaction(db, jobRff)

	authTransaction, err := authG.NewTransaction(db, userRff)
	if err != nil {
		return nil, nil, nil, err
	}

	appTransaction, err := appG.NewTransaction(db, rff)
	if err != nil {
		return nil, nil, nil, err
	}
	return jobTransaction, authTransaction, appTransaction, nil
}

func run(ctx context.Context, cfg *config.Config, transaction service.Transaction, db *gorm.DB, pf appS.ProcessorFactory, rff appG.RepositoryFactoryFunc, authTransaction authS.Transaction, jobTransaction jobS.Transaction, appTransaction appS.Transaction, synthesizerClient appS.SynthesizerClient, translatorClient pluginCommonS.TranslatorClient, tatoebaClient pluginCommonS.TatoebaClient, newIteratorFunc controller.NewIteratorFunc) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	eg.Go(func() error {
		return appServer(ctx, cfg, db, pf, rff, authTransaction, jobTransaction, appTransaction, synthesizerClient, translatorClient, tatoebaClient, newIteratorFunc)
	})
	eg.Go(func() error {
		return jobServer(ctx, cfg, jobTransaction, appTransaction)
	})
	eg.Go(func() error {
		return libG.MetricsServerProcess(ctx, cfg.App.MetricsPort, cfg.Shutdown.TimeSec1)
	})
	eg.Go(func() error {
		return libG.SignalWatchProcess(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		logrus.Error(err)
		return 1
	}
	return 0
}

func appServer(ctx context.Context, cfg *config.Config, db *gorm.DB, pf appS.ProcessorFactory, rff appG.RepositoryFactoryFunc, authTransaction authS.Transaction, jobTransaction jobS.Transaction, appTransaction appS.Transaction, synthesizerClient appS.SynthesizerClient, translatorClient pluginCommonS.TranslatorClient, tatoebaClient pluginCommonS.TatoebaClient, newIteratorFunc controller.NewIteratorFunc) error {
	// cors
	corsConfig := libconfig.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return liberrors.Errorf("corsConfig.Validate. err: %w", err)
	}

	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	googleAuthClient := authG.NewGoogleAuthClient(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleCallbackURL, time.Duration(cfg.Auth.APITimeoutSec)*time.Second)

	registerAppUserCallback := func(ctx context.Context, organizationName string, appUser userD.AppUserModel) error {
		rf, err := rff(ctx, db)
		if err != nil {
			return err
		}
		userRf, err := rf.NewUserRepositoryFactory(ctx)
		if err != nil {
			return err
		}
		return callback(ctx, cfg.App.TestUserEmail, pf, rf, userRf, organizationName, appUser)
	}

	googleUserUsecase := authU.NewGoogleUserUsecase(authTransaction, googleAuthClient, authTokenManager, registerAppUserCallback)
	guestUserUsecase := authU.NewGuestUserUsecase(authTransaction, authTokenManager)
	studentUsecaseWorkbook := studentU.NewStudentUsecaseWorkbook(appTransaction, pf)
	studentUsecaseProblem := studentU.NewStudentUsecaseProblem(appTransaction, pf)
	studentUseCaseStudy := studentU.NewStudentUsecaseStudy(appTransaction, pf)
	studentUsecaseAudio := studentU.NewStudentUsecaseAudio(appTransaction, pf, synthesizerClient)
	studentUsecaseStat := studentU.NewStudentUsecaseStat(appTransaction, pf)

	publicRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitAuthRouterFunc(googleUserUsecase, guestUserUsecase, authTokenManager),
	}
	privateRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitWorkbookRouterFunc(studentUsecaseWorkbook),
		controller.NewInitProblemRouterFunc(studentUsecaseProblem, newIteratorFunc),
		controller.NewInitStudyRouterFunc(studentUseCaseStudy),
		controller.NewInitAudioRouterFunc(studentUsecaseAudio),
		controller.NewInitStatRouterFunc(studentUsecaseStat),
	}
	pluginRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitTranslationRouterFunc(translatorClient),
		controller.NewInitTatoebaRouterFunc(tatoebaClient),
	}

	router, err := controller.NewAppRouter(publicRouterGroupFunc, privateRouterGroupFunc, pluginRouterGroupFunc, authTokenManager, corsConfig, cfg.App, cfg.Auth, cfg.Debug)
	if err != nil {
		panic(err)
	}

	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Title = cfg.App.Name
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = cfg.Swagger.Host
		docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}
	}

	httpServer := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.App.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logrus.Printf("http server listening at %v", httpServer.Addr)

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logrus.Infof("Server forced to shutdown. err: %v", err)
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}

func jobServer(ctx context.Context, cfg *config.Config, jobTransaction jobS.Transaction, appTransaction appS.Transaction) error {
	jobService, err := jobS.NewJobService(ctx, jobTransaction)
	if err != nil {
		return err
	}
	router, err := controller.NewJobRouter(appTransaction, jobService, cfg.Debug)
	if err != nil {
		return err
	}

	httpServer := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.App.JobPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logrus.Printf("job server listening at %v", httpServer.Addr)

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logrus.Infof("Server forced to shutdown. err: %v", err)
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}

func initPf(synthesizerClient appS.SynthesizerClient, translatorClient pluginCommonS.TranslatorClient, tatoebaClient pluginCommonS.TatoebaClient) (appS.ProcessorFactory, map[string]func(context.Context, *gorm.DB) (appS.ProblemRepository, error), map[string]appS.ProblemImportProcessor) {

	englishWordProblemProcessor := pluginEnglishS.NewEnglishWordProblemProcessor(synthesizerClient, translatorClient, tatoebaClient, pluginEnglishGateway.NewEnglishWordProblemAddParameterCSVReader)
	englishPhraseProblemProcessor := pluginEnglishS.NewEnglishPhraseProblemProcessor(synthesizerClient, translatorClient)
	englishSentenceProblemProcessor := pluginEnglishS.NewEnglishSentenceProblemProcessor(synthesizerClient, translatorClient, pluginEnglishGateway.NewEnglishSentenceProblemAddParameterCSVReader)

	problemAddProcessor := map[string]appS.ProblemAddProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemUpdateProcessor := map[string]appS.ProblemUpdateProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemRemoveProcessor := map[string]appS.ProblemRemoveProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemImportProcessor := map[string]appS.ProblemImportProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemQuotaProcessor := map[string]appS.ProblemQuotaProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}

	englishWordProblemRepositoryFunc := func(ctx context.Context, db *gorm.DB) (appS.ProblemRepository, error) {
		// fmt.Println("-------Word")
		return pluginEnglishGateway.NewEnglishWordProblemRepository(db, synthesizerClient, pluginEnglishDomain.EnglishWordProblemType)
	}
	englishPhraseProblemRepositoryFunc := func(ctx context.Context, db *gorm.DB) (appS.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishPhraseProblemRepository(db, synthesizerClient, pluginEnglishDomain.EnglishPhraseProblemType)
	}
	englishSentenceProblemRepositoryFunc := func(ctx context.Context, db *gorm.DB) (appS.ProblemRepository, error) {
		// fmt.Println("-------Sentence")
		return pluginEnglishGateway.NewEnglishSentenceProblemRepository(db, synthesizerClient, pluginEnglishDomain.EnglishSentenceProblemType)
	}

	pf := appS.NewProcessorFactory(problemAddProcessor, problemUpdateProcessor, problemRemoveProcessor, problemImportProcessor, problemQuotaProcessor)

	problemRepositories := map[string]func(context.Context, *gorm.DB) (appS.ProblemRepository, error){
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemRepositoryFunc,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemRepositoryFunc,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemRepositoryFunc,
	}
	return pf, problemRepositories, problemImportProcessor
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init log
	if err := libconfig.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, err
	}

	// init tracer
	tp, err := libconfig.InitTracerProvider(cfg.App.Name, cfg.Trace)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to InitTracerProvider. err: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := libconfig.InitDB(cfg.DB, sqls.SQL)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to InitDB. err: %w", err)
	}

	// userRff := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
	// 	return userG.NewRepositoryFactory(db)
	// }
	// userS.InitSystemAdmin(userRff)

	return cfg, db, sqlDB, tp, nil
}

func initApp1(ctx context.Context, db *gorm.DB, password string) error {
	logger := log.FromContext(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userG.NewRepositoryFactory(ctx, db)
		if err != nil {
			return err
		}

		systemAdmin, err := userS.NewSystemAdmin(ctx, userRf)
		if err != nil {
			return err
		}

		organization, err := systemAdmin.FindOrganizationByName(ctx, "cocotola")
		if err != nil {
			if !errors.Is(err, userS.ErrOrganizationNotFound) {
				return liberrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			firstOwnerAddParam, err := userS.NewFirstOwnerAddParameter("cocotola-owner", "Owner(cocotola)", password)
			if err != nil {
				return liberrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			organizationAddParameter, err := userS.NewOrganizationAddParameter("cocotola", firstOwnerAddParam)
			if err != nil {
				return liberrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			organizationID, err := systemAdmin.AddOrganization(ctx, organizationAddParameter)
			if err != nil {
				return liberrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			logger.Infof("organizationID: %d", organizationID)
			return nil
		}
		logger.Infof("organization: %d", organization.GetID())
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func initApp2(ctx context.Context, appTransaction service.Transaction) error {
	if err := initApp2_1(ctx, appTransaction); err != nil {
		return liberrors.Errorf("failed to initApp2_1. err: %w", err)
	}

	if err := initApp2_2(ctx, appTransaction); err != nil {
		return liberrors.Errorf("failed to initApp2_2. err: %w", err)
	}

	if err := initApp2_3(ctx, appTransaction); err != nil {
		return liberrors.Errorf("failed to initApp2_3. err: %w", err)
	}

	return nil
}

func initApp2_1(ctx context.Context, appTransaction service.Transaction) error {
	var propertiesSystemStudentID userD.AppUserID

	if err := appTransaction.Do(ctx, func(rf appS.RepositoryFactory) error {
		userRf, err := rf.NewUserRepositoryFactory(ctx)
		if err != nil {
			return err
		}

		systemAdmin, err := userS.NewSystemAdmin(ctx, userRf)
		if err != nil {
			return liberrors.Errorf("NewSystemAdmin. err: %w", err)
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, "cocotola")
		if err != nil {
			return liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemStudent, err := systemOwner.FindAppUserByLoginID(ctx, appS.SystemStudentLoginID)
		if err != nil {
			if !errors.Is(err, userS.ErrAppUserNotFound) {
				return liberrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
			}

			param, err := userS.NewAppUserAddParameter(appS.SystemStudentLoginID, "SystemStudent(cocotola)", []string{}, map[string]string{})
			if err != nil {
				return liberrors.Errorf("failed to NewAppUserAddParameter. err: %w", err)
			}

			systemStudentID, err := systemOwner.AddAppUser(ctx, param)
			if err != nil {
				return liberrors.Errorf("failed to AddAppUser. err: %w", err)
			}

			propertiesSystemStudentID = systemStudentID
		} else {
			propertiesSystemStudentID = userD.AppUserID(systemStudent.GetID())
		}
		return nil
	}); err != nil {
		return err
	}

	appS.SetSystemStudentID(propertiesSystemStudentID)

	return nil
}

func initApp2_2(ctx context.Context, appTransaction service.Transaction) error {
	var propertiesSystemSpaceID userD.SpaceID

	if err := appTransaction.Do(ctx, func(rf appS.RepositoryFactory) error {
		userRf, err := rf.NewUserRepositoryFactory(ctx)
		if err != nil {
			return err
		}

		systemAdmin, err := userS.NewSystemAdmin(ctx, userRf)
		if err != nil {
			return liberrors.Errorf("NewSystemAdmin. err: %w", err)
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, appS.OrganizationName)
		if err != nil {
			return liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemSpace, err := systemOwner.FindSystemSpace(ctx)
		if err != nil {
			if !errors.Is(err, userS.ErrSpaceNotFound) {
				return liberrors.Errorf("failed to FindSystemSpace. err: %w", err)
			}

			spaceID, err := systemOwner.AddSystemSpace(ctx)
			if err != nil {
				return liberrors.Errorf("failed to AddSystemSpace. err: %w", err)
			}

			propertiesSystemSpaceID = spaceID
		} else {
			propertiesSystemSpaceID = userD.SpaceID(systemSpace.GetID())
		}

		return nil
	}); err != nil {
		return err
	}

	appS.SetSystemSpaceID(propertiesSystemSpaceID)

	return nil
}

func initApp2_3(ctx context.Context, appTransaction service.Transaction) error {
	var propertiesTatoebaWorkbookID appD.WorkbookID
	if err := appTransaction.Do(ctx, func(rf appS.RepositoryFactory) error {
		userRf, err := rf.NewUserRepositoryFactory(ctx)
		if err != nil {
			return err
		}

		systemAdmin, err := userS.NewSystemAdmin(ctx, userRf)
		if err != nil {
			return liberrors.Errorf("NewSystemAdmin. err: %w", err)
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, appS.OrganizationName)
		if err != nil {
			return liberrors.Errorf("FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemStudentAppUser, err := systemOwner.FindAppUserByLoginID(ctx, appS.SystemStudentLoginID)
		if err != nil {
			return liberrors.Errorf("FindAppUserByLoginID. err: %w", err)
		}

		systemStudentModel, err := appD.NewSystemStudentModel(systemStudentAppUser)
		if err != nil {
			return liberrors.Errorf("NewSystemStudentModel. err: %w", err)
		}

		systemStudent, err := appS.NewSystemStudent(rf, systemStudentModel)
		if err != nil {
			return liberrors.Errorf("NewSystemStudent. err: %w", err)
		}

		tatoebaWorkbook, err := systemStudent.FindWorkbookFromSystemSpace(ctx, appS.TatoebaWorkbookName)
		if err != nil {
			if !errors.Is(err, appS.ErrWorkbookNotFound) {
				return err
			}

			paramToAddWorkbook, err := appS.NewWorkbookAddParameter(pluginEnglishDomain.EnglishSentenceProblemType, appS.TatoebaWorkbookName, appD.Lang2JA, "", map[string]string{})
			if err != nil {
				return liberrors.Errorf("NewWorkbookAddParameter. err: %w", err)
			}

			tatoebaWorkbookID, err := systemStudent.AddWorkbookToSystemSpace(ctx, paramToAddWorkbook)
			if err != nil {
				return liberrors.Errorf("AddWorkbookToSystemSpace. err: %w", err)
			}

			propertiesTatoebaWorkbookID = tatoebaWorkbookID
		} else {
			propertiesTatoebaWorkbookID = tatoebaWorkbook.GetWorkbookID()
		}

		return nil
	}); err != nil {
		return err
	}

	appS.SetTatoebaWorkbookID(propertiesTatoebaWorkbookID)

	return nil
}

func callback(ctx context.Context, testUserEmail string, pf appS.ProcessorFactory, repo appS.RepositoryFactory, userRepo userS.RepositoryFactory, organizationName string, appUser userD.AppUserModel) error {
	logger := log.FromContext(ctx)
	logger.Infof("callback. loginID: %s", appUser.GetLoginID())

	if appUser.GetLoginID() == testUserEmail {
		studentModel, err := appD.NewStudentModel(appUser)
		if err != nil {
			return liberrors.Errorf("NewStudentModel. err: %w", err)
		}

		student, err := appS.NewStudent(ctx, pf, repo, studentModel)
		if err != nil {
			return liberrors.Errorf("NewStudent. err: %w", err)
		}

		if err := english_word.CreateDemoWorkbook(ctx, student); err != nil {
			return liberrors.Errorf("english_word.CreateDemoWorkbook. err: %w", err)
		}

		if err := english_word.Create20NGSLWorkbook(ctx, student); err != nil {
			return liberrors.Errorf("english_word.Create20NGSLWorkbook. err: %w", err)
		}

		if err := english_sentence.CreateFlushWorkbook(ctx, student); err != nil {
			return liberrors.Errorf("english_sentence.CreateFlushWorkbook. err: %w", err)
		}

		// if err := english_word.Create300NGSLWorkbook(ctx, student); err != nil {
		// 	return liberrors.Errorf("failed to Create300NGSLWorkbook. err: %w", err)
		// }
	}

	return nil
}
