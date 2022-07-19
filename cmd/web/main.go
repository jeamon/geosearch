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
	logger.Info("starting web server.", zap.String("version", gitTag))

	configData, err := config.NewConfig("./config/config.yaml")
	if err != nil {
		log.Fatal("Failed to initialize web configs", err)
	}

	db, count, err := helpers.LoadJobsDetailsFromCSVFile(configData.Database.CSVFile)
	if err != nil {
		log.Fatalf("Failed to load jobs details: %v", err)
	}
	repo := repository.NewJobsDetailsRepository(logger, db)
	service := service.New(logger, configData, repo)
	handler := handlers.NewWEBHandler(logger, configData, service)
	logger.Info("loaded jobs details and initialized the web server settings", zap.Int("number_jobs", count))

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Static("/assets", "./assets/static")
	router.StaticFile("/favicon.ico", "./assets/static/favicon.ico")
	router.LoadHTMLGlob("assets/templates/*.html")
	handler.InitWEBRoutes(router)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	webServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", configData.WEBServer.Host, configData.WEBServer.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go helpers.ShutdownServer(logger, webServer)

	logger.Info("started web server", zap.String("host", configData.WEBServer.Host), zap.String("port", configData.WEBServer.Port))
	if err := webServer.ListenAndServeTLS(configData.WEBServer.CertsFile, configData.WEBServer.KeyFile); err != nil && err != http.ErrServerClosed {
		logger.Error("failed to bring up the web server",
			zap.String("host", configData.WEBServer.Host),
			zap.String("port", configData.WEBServer.Port),
			zap.Error(err))
	}
}
