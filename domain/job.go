package domain

type Job struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type JobsDetailsDatabase struct {
	Jobs map[string]Job
}

type JobSearchRequest struct {
	Title     string  `json:"title" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

type SearchJobsResponse struct {
	RequestId string `json:"request_id"`
	Message   string `json:"message"`
	Jobs      []Job  `json:"jobs"`
}

type GetJobResponse struct {
	RequestId string `json:"request_id"`
	Message   string `json:"message"`
	Job       Job    `json:"job"`
}
