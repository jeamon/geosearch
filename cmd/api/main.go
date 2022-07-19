package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jeamon/geosearch/config"
	"github.com/jeamon/geosearch/handlers"
	"github.com/jeamon/geosearch/helpers"
	"github.com/jeamon/geosearch/repository"
	"github.com/jeamon/geosearch/service"
	"go.uber.org/zap"
)

var (
	gitCommit string
	gitTag    string
	buildTime string
)

func main() {
	// Setup the logger with default fields.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize zap logger", err)
	}
	defer logger.Sync()
	logger = logger.With(zap.String("git_commit", gitCommit), zap.String("git_tag", gitTag), zap.String("build_time", buildTime))
	logger.Info("starting api server.", zap.String("version", gitTag))

	configData, err := config.NewConfig("./config/config.yaml")
	if err != nil {
		log.Fatal("Failed to initialize api configs", err)
	}

	db, count, err := helpers.LoadJobsDetailsFromCSVFile(configData.Database.CSVFile)
	if err != nil {
		log.Fatalf("Failed to load jobs details: %v", err)
	}
	repo := repository.NewJobsDetailsRepository(logger, db)
	service := service.New(logger, configData, repo)
	handler := handlers.NewAPIHandler(logger, configData, service)
	logger.Info("loaded jobs details and initialized the api server settings", zap.Int("number_jobs", count))

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginzap.RecoveryWithZap(logger, true))
	handler.InitAPIRoutes(router)
	// simple route to quickly check server availability.
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	apiServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", configData.APIServer.Host, configData.APIServer.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go helpers.ShutdownServer(logger, apiServer)

	logger.Info("started api server", zap.String("host", configData.APIServer.Host), zap.String("port", configData.APIServer.Port))
	if err := apiServer.ListenAndServeTLS(configData.APIServer.CertsFile, configData.APIServer.KeyFile); err != nil && err != http.ErrServerClosed {
		logger.Error("failed to bring up the api server",
			zap.String("host", configData.APIServer.Host),
			zap.String("port", configData.APIServer.Port),
			zap.Error(err))
	}
}
