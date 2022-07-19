package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jeamon/geosearch/config"
	"github.com/jeamon/geosearch/service"
	"go.uber.org/zap"
)

type WEBHandler struct {
	logger  *zap.Logger
	config  *config.Config
	service service.Service
}

// NewWEBHandler provides a new instance of WEBHandler.
func NewWEBHandler(logger *zap.Logger, config *config.Config, s service.Service) *WEBHandler {
	return &WEBHandler{logger: logger, config: config, service: s}
}

func (web *WEBHandler) InitWEBRoutes(router *gin.Engine) *gin.Engine {
	router.Use(RequestLoggerMiddleware(web.logger))
	router.NoRoute(web.NotFound)
	router.GET("/", web.Index)
	router.GET("/error", web.ErrorPage)
	router.GET("/jobs", web.Home)
	router.GET("/jobs/view/:id", web.Get)
	router.GET("/jobs/search", web.Search)
	return router
}

// NotFound responds to requests towards non implemented routes.
func (web *WEBHandler) NotFound(c *gin.Context) {
	web.logger.Error("url does not exist", zap.String("requestid", c.GetString("x-requestid")))
	c.HTML(http.StatusBadRequest, "error.html", gin.H{
		"Status": 404,
		"Error":  "Not Found",
		"Infos":  "The URL requested does not exist.",
	})
}

// Index handles request to root domain.
func (web *WEBHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// Home handles and displays jobs based on location.
func (web *WEBHandler) Home(c *gin.Context) {
	lat := c.DefaultQuery("lat", "")
	lng := c.DefaultQuery("lng", "")
	if lat == "" || lng == "" {
		web.logger.Error("empty geolocation data", zap.String("requestid", c.GetString("x-requestid")), zap.Error(nil))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "The URL submitted is malformatted. Geolocation data is required.",
		})
		return
	}

	latitude, err1 := strconv.ParseFloat(lat, 64)
	longitude, err2 := strconv.ParseFloat(lng, 64)
	if err1 != nil || err2 != nil {
		web.logger.Error("no valid geolocation data", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err1), zap.Error(err2))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "Geolocation data submitted is invalid.",
		})
		return
	}
	openings, nearest, err := web.service.Load(longitude, latitude)
	if err != nil {
		web.logger.Error("unable to load random jobs based on location", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Status": 503,
			"Error":  "Error Occured During Processing",
			"Infos":  "The platform failed to process the geolocation data. Try again later.",
		})
		return
	}

	c.HTML(http.StatusOK, "jobs.html", gin.H{
		"nearestJobs":      nearest,
		"availableJobs":    openings[:len(openings)-1],
		"lastAvailableJob": openings[len(openings)-1],
		"userLatitude":     latitude,
		"userLongitude":    longitude,
		"title":            "",
		"isSearch":         false,
		"mapsAPIKEY":       web.config.MapsAPIKEY,
	})
}

func (web *WEBHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		web.logger.Error("bad request. invalid job id", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "The Job ID provided is invalid.",
		})
		return
	}

	job, err := web.service.Get(id)
	if err != nil {
		web.logger.Error("unable to find requested job", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Status": 404,
			"Error":  "Not Found",
			"Infos":  "The Job ID requested does not exist.",
		})
		return
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"job": job,
	})
}

func (web *WEBHandler) Search(c *gin.Context) {
	text := c.DefaultQuery("title", "")
	title := strings.Join(strings.Fields(text), " ")
	lat := c.DefaultQuery("lat", "")
	lng := c.DefaultQuery("lng", "")
	if title == "" || lat == "" || lng == "" {
		web.logger.Error("empty data into the query string", zap.String("requestid", c.GetString("x-requestid")), zap.Error(nil))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "The URL submitted is malformatted. Geolocation data and Job Title are required.",
		})
		return
	}

	latitude, err1 := strconv.ParseFloat(lat, 64)
	longitude, err2 := strconv.ParseFloat(lng, 64)
	if err1 != nil || err2 != nil {
		web.logger.Error("no valid geolocation data", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err1), zap.Error(err2))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "Geolocation data submitted is invalid.",
		})
		return
	}

	jobs, err := web.service.Search(title, longitude, latitude)
	if err != nil {
		web.logger.Error("unable to search for requested jobs title", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Status": 503,
			"Error":  "Internal Server Error",
			"Infos":  "Failed to search for requested jobs title. Try again later.",
		})
		return
	}

	c.HTML(http.StatusOK, "jobs.html", gin.H{
		"nearestJobs":      jobs[0:3],
		"availableJobs":    jobs[:len(jobs)-1],
		"lastAvailableJob": jobs[len(jobs)-1],
		"userLatitude":     latitude,
		"userLongitude":    longitude,
		"title":            title,
		"isSearch":         true,
		"mapsAPIKEY":       web.config.MapsAPIKEY,
	})
}

// ErrorPage displays page for common errors.
func (web *WEBHandler) ErrorPage(c *gin.Context) {
	issue := c.DefaultQuery("issue", "")
	if issue == "" {
		web.logger.Error("bad request. no error code provided", zap.String("requestid", c.GetString("x-requestid")), zap.Error(nil))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "The URL requested is invalid. Missing issue code.",
		})
		return
	}
	var infos string
	switch issue {
	case "location-no-support":
		infos = "Sorry, your browser does not support geolocation service."
	case "location-perm-denied":
		infos = "Request for geolocation denied. You must allow geolocation to continue using this website. Retry."
	case "location-pos-unvailable":
		infos = "Sorry, location information is unvailable. Try again later."
	case "location-timeout":
		infos = "Sorry, the request to get your current location timed out. Reload the page."
	case "location-unknown":
		infos = "Sorry, an unknown error occured. Try again later."
	}

	if infos == "" {
		web.logger.Error("bad request. invalid error code provided", zap.String("requestid", c.GetString("x-requestid")), zap.Error(nil))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Status": 400,
			"Error":  "Bad Request",
			"Infos":  "The URL requested is invalid. Unknown issue code.",
		})
		return
	}
	c.HTML(http.StatusForbidden, "error.html", gin.H{
		"Status": 403,
		"Error":  "Failed to Process",
		"Infos":  infos,
	})
}
