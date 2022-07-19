package handlers

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/jeamon/geosearch/helpers"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware logs the metadata and the body of the requests.
func RequestLoggerMiddleware(logger *zap.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		requestId := helpers.GenerateUID()
		buf, _ := ioutil.ReadAll(c.Request.Body)
		data, err := helpers.ReadBody(ioutil.NopCloser(bytes.NewBuffer(buf)))

		logger.Info(
			"received request on:",
			zap.String("url", c.Request.URL.Path),
			zap.String("authorization", c.GetHeader("Authorization")),
			zap.String("method", c.Request.Method),
			zap.String("requestid", requestId),
			zap.String("ip", helpers.GetIP(c.Request)),
			zap.String("agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.String("body", data),
			zap.Error(err),
		)
		c.Set("x-requestid", requestId)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Next()
	}
}
