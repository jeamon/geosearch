package handlers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jeamon/geosearch/config"
	"github.com/jeamon/geosearch/domain"
	"github.com/jeamon/geosearch/service"
	"go.uber.org/zap"
)

type errResponse struct {
	RequestId        string `json:"request_id"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
}

type APIHandler struct {
	logger  *zap.Logger
	config  *config.Config
	service service.Service
}

// NewAPIHandler provides a new instance of APIHandler.
func NewAPIHandler(logger *zap.Logger, config *config.Config, s service.Service) *APIHandler {
	return &APIHandler{logger: logger, config: config, service: s}
}

func (api *APIHandler) InitAPIRoutes(router *gin.Engine) *gin.Engine {
	router.Use(RequestLoggerMiddleware(api.logger))
	router.NoRoute(api.NotFound)
	apiRouter := router.Group("/api/v1")
	apiRouter.Use(cors.New(cors.Config{
		AllowWildcard:    true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Accept", "content-type", "User-Agent", "Accept-Language", "Referer", "DNT", "Connection", "Pragma", "Cache-Control", "TE", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiRouter.POST("/jobs/search", api.Search)
	apiRouter.GET("/jobs/:id", api.Get)
	return router
}

// NotFound responds to requests towards non implemented routes.
func (api *APIHandler) NotFound(c *gin.Context) {
	api.logger.Error("route does not exist", zap.String("requestid", c.GetString("x-requestid")))
	c.JSON(404, errResponse{
		RequestId:        c.GetString("x-requestid"),
		Message:          "invalid request. make sure to use the exact endpoint.",
		DeveloperMessage: "endpoint called with that method does not exist.",
	})
}

func (api *APIHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		api.logger.Error("bad request. invalid job id", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.JSON(400, errResponse{
			RequestId:        c.GetString("x-requestid"),
			Message:          "bad request. malformatted job id",
			DeveloperMessage: err.Error(),
		})
		return
	}

	job, err := api.service.Get(id)
	if err != nil {
		api.logger.Error("unable to find requested job", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.JSON(404, errResponse{
			RequestId:        c.GetString("x-requestid"),
			Message:          "invalid request. unexpected request body data format",
			DeveloperMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.GetJobResponse{
		RequestId: c.GetString("x-requestid"),
		Message:   "found",
		Job:       job,
	})
}

// Search process request to find a configured number of nearest jobs.
func (api *APIHandler) Search(c *gin.Context) {
	var req domain.JobSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.logger.Error("unable to bind store request input", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.JSON(400, errResponse{
			RequestId:        c.GetString("x-requestid"),
			Message:          "invalid request. unexpected request body data format",
			DeveloperMessage: err.Error(),
		})
		return
	}

	jobs, err := api.service.Search(req.Title, req.Longitude, req.Latitude)
	if err != nil {
		api.logger.Error("unable to search for requested jobs", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.JSON(500, errResponse{
			RequestId:        c.GetString("x-requestid"),
			Message:          "an error occured while searching for requested jobs",
			DeveloperMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SearchJobsResponse{
		RequestId: c.GetString("x-requestid"),
		Message:   "search successfully",
		Jobs:      jobs,
	})
}
