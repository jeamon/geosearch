package helpers

import (
	"bytes"
	"context"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jeamon/geosearch/domain"
	"go.uber.org/zap"
)

// generateUID provides a random uid to trace a given request or to identify a job details.
func GenerateUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// getIP provides the source IP of the client agent.
func GetIP(r *http.Request) string {
	// Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	// Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP = net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	// Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}

// readBody converts the body into a string.
func ReadBody(reader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	s := buf.String()
	return s, nil
}

// ShutdownServer handles a stop signal and gracefully teardown the server.
func ShutdownServer(logger *zap.Logger, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt, syscall.SIGABRT)
	<-quit
	logger.Info("received exit signal. shutting down the server.")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to gracefully shutting down the server", zap.Error(err))
		if err == context.DeadlineExceeded {
			logger.Info("shutting down deadline of 60 secs exceeded.")
		} else {
			logger.Info("error occured when closing underlying listeners.")
		}
		return
	}
	logger.Info("server gracefully shutdown.")
}

// DistanceBetween computes the distance in kilometers between two coordinates.
// source https://www.geodatasource.com/developers/go
// source https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49
func DistanceBetween(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515 * 1.609344
	return dist
}

func IsWithinRadius(distance float64, radius float64) bool {
	return distance <= radius
}

// FillWithEmptyJobs adds empty job details for UI display.
func FillWithEmptyJobs(jobs *[]domain.Job, gap int) {
	if gap != 0 {
		for i := 0; i < gap; i++ {
			*jobs = append(*jobs, domain.Job{
				ID:    "#",
				Title: "No Job Found",
			})
		}
	}
}
