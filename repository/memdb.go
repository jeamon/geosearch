package repository

import (
	"fmt"
	"strings"

	"github.com/jeamon/geosearch/domain"
	"github.com/jeamon/geosearch/helpers"
	"go.uber.org/zap"
)

type JobsDetailsRepository struct {
	logger   *zap.Logger
	Database *domain.JobsDetailsDatabase
}

func NewJobsDetailsRepository(logger *zap.Logger, db *domain.JobsDetailsDatabase) *JobsDetailsRepository {
	return &JobsDetailsRepository{
		logger:   logger,
		Database: db,
	}
}

func (repo *JobsDetailsRepository) Search(title string, userLongitude float64, userLatitude float64, number int, radius float64) ([]domain.Job, error) {
	var distance float64
	var foundJobs []domain.Job
	num := 0
	for _, job := range repo.Database.Jobs {
		if num == number {
			break
		}
		if strings.Contains(job.Title, strings.Title(title)) {
			distance = helpers.DistanceBetween(userLatitude, userLongitude, job.Latitude, job.Longitude)
			if helpers.IsWithinRadius(distance, float64(radius)) {
				foundJobs = append(foundJobs, job)
				num++
			}
		}
	}
	helpers.FillWithEmptyJobs(&foundJobs, 5-len(foundJobs))
	return foundJobs, nil
}

func (repo *JobsDetailsRepository) Find(id string) (domain.Job, error) {
	if job, found := repo.Database.Jobs[id]; found {
		return job, nil
	}
	return domain.Job{}, fmt.Errorf("job does not exist. id : %s", id)
}

func (repo *JobsDetailsRepository) Load(userLongitude float64, userLatitude float64) ([]domain.Job, []domain.Job, error) {
	var availableJobs []domain.Job
	var nearestJobs []domain.Job
	repo.loadRandomJobs(&availableJobs, 5)
	repo.findNearestJobs(&nearestJobs, 3, 5, userLatitude, userLongitude)
	helpers.FillWithEmptyJobs(&availableJobs, 5-len(availableJobs))
	helpers.FillWithEmptyJobs(&nearestJobs, 3-len(nearestJobs))
	return availableJobs, nearestJobs, nil
}

// loadRandomJobs provides a given number of random jobs.
func (repo *JobsDetailsRepository) loadRandomJobs(jobs *[]domain.Job, number int) {
	num := 0
	for _, job := range repo.Database.Jobs {
		if num == number {
			break
		}
		*jobs = append(*jobs, job)
		num++
	}
}

// findNearestJobs provides a given number of nearest jobs from source coordinates and based on radius.
func (repo *JobsDetailsRepository) findNearestJobs(jobs *[]domain.Job, number int, radius float64, slat float64, slng float64) {
	num := 0
	var distance float64
	for _, job := range repo.Database.Jobs {
		if num == number {
			break
		}
		distance = helpers.DistanceBetween(slat, slng, job.Latitude, job.Longitude)
		if helpers.IsWithinRadius(distance, radius) {
			*jobs = append(*jobs, job)
			num++
		}
	}
}
