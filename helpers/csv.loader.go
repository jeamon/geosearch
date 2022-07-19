package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jeamon/geosearch/domain"
)

func LoadJobsDetailsFromCSVFile(csvFilePath string) (*domain.JobsDetailsDatabase, int, error) {
	var count int
	db := &domain.JobsDetailsDatabase{}
	db.Jobs = make(map[string]domain.Job)
	csvFile, err := os.Open(csvFilePath)
	if err != nil {

		return db, count, fmt.Errorf("failed to load the csv file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	// read all entries into slice of slice of string
	allRecords, err := reader.ReadAll()
	if err != nil {
		return db, count, fmt.Errorf("failed to load the all csv records: %v", err)
	}

	// no need to continue if the file does not have any records.
	if len(allRecords) <= 1 {
		return db, count, fmt.Errorf("csv file has no jobs records")
	}
	// remove headers record.
	allRecords = allRecords[1:]

	// build a job object for each record and add to database.
	for _, record := range allRecords {
		lng, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}
		lt, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}
		job := domain.Job{
			ID:        GenerateUID(),
			Title:     strings.Title(strings.Trim(record[0], "\"")),
			Longitude: lng,
			Latitude:  lt,
		}
		db.Jobs[job.ID] = job
	}
	return db, len(db.Jobs), nil
}
