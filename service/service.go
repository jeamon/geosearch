package service

import (
	"github.com/jeamon/geosearch/config"
	"github.com/jeamon/geosearch/domain"
	"go.uber.org/zap"
)

type Repository interface {
	Search(title string, userLongitude float64, userLatitude float64, number int, radius float64) ([]domain.Job, error)
	Load(userLongitude float64, userLatitude float64) ([]domain.Job, []domain.Job, error)
	Find(id string) (domain.Job, error)
}

type Service interface {
	Search(jobTitle string, userLongitude float64, userLatitude float64) ([]domain.Job, error)
	Load(userLongitude float64, userLatitude float64) ([]domain.Job, []domain.Job, error)
	Get(id string) (domain.Job, error)
}

type service struct {
	logger *zap.Logger
	config *config.Config
	repo   Repository
}

func New(logger *zap.Logger, config *config.Config, repo Repository) Service {
	return &service{logger: logger, config: config, repo: repo}
}

func (s *service) Search(title string, userLongitude float64, userLatitude float64) ([]domain.Job, error) {
	return s.repo.Search(title, userLongitude, userLatitude, s.config.SearchResults.Number, s.config.SearchResults.Radius)
}

func (s *service) Get(id string) (domain.Job, error) {
	return s.repo.Find(id)
}

func (s *service) Load(userLongitude float64, userLatitude float64) ([]domain.Job, []domain.Job, error) {
	return s.repo.Load(userLongitude, userLatitude)
}
