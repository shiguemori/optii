package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"optii/api"
	"optii/models"
)

type JobService interface {
	CreateJob(job *models.CreateJobRequest) (*models.Job, error, int)
}

type jobService struct {
	api api.OptiiApi
}

func NewJobService(api api.OptiiApi) JobService {
	return &jobService{
		api: api,
	}
}

func (s *jobService) departmentExists(department string) (bool, error) {
	dep, err := s.api.GetDepartments(department, 0, 0)
	if err != nil {
		return false, err
	}

	return dep != nil, nil
}

func (s *jobService) jobItemExists(jobItem string) (bool, error) {
	j, err := s.api.GetJobItems(0, 0, jobItem)
	if err != nil {
		return false, err
	}

	return j != nil, nil
}

func (s *jobService) locationsExist(locations []string) (bool, error) {
	for i := range locations {
		tempMap := make(map[string]string)
		tempMap["displayName"] = locations[i]
		location, err := s.api.GetLocations(tempMap)
		if err != nil {
			return false, err
		}
		if location == nil {
			return false, nil
		}
	}

	return true, nil
}

// CreateJob creates a new job in Optii.
// It returns the created job if successful and an error (along with the HTTP status code) if there's any issue.
// This function uses goroutines to make API calls asynchronously.
// Currently, a hardcoded newJob *models.Job is created in the "Optii" for testing purposes using Swagger example, but the API may not be functioning correctly.
func (s *jobService) CreateJob(job *models.CreateJobRequest) (*models.Job, error, int) {
	departmentResultChan := make(chan bool)
	jobItemResultChan := make(chan bool)
	locationsResultChan := make(chan bool)

	doAsyncQuery := func(existsFunc func() (bool, error), resultChan chan bool) {
		go func() {
			exists, err := existsFunc()
			if err != nil || !exists {
				resultChan <- false
			} else {
				resultChan <- true
			}
		}()
	}

	if job.Department != nil {
		doAsyncQuery(func() (bool, error) {
			return s.departmentExists(*job.Department)
		}, departmentResultChan)
	}

	if job.JobItem != nil {
		doAsyncQuery(func() (bool, error) {
			return s.jobItemExists(*job.JobItem)
		}, jobItemResultChan)
	}

	if job.Locations != nil {
		doAsyncQuery(func() (bool, error) {
			return s.locationsExist(job.Locations)
		}, locationsResultChan)
	}

	departmentExists := <-departmentResultChan
	jobItemExists := <-jobItemResultChan
	locationsExist := <-locationsResultChan

	if !departmentExists {
		return nil, fmt.Errorf("invalid department"), http.StatusBadRequest
	}

	if !jobItemExists {
		return nil, fmt.Errorf("invalid job item"), http.StatusBadRequest
	}

	if !locationsExist {
		return nil, fmt.Errorf("invalid location"), http.StatusBadRequest
	}

	var newJob *models.Job
	newJob = &models.Job{
		Item: models.Item{
			Name: "drink",
		},
		Priority: "highest",
		Department: models.Department{
			Id: 7,
		},
		Role: models.Roles{
			Id: 10,
		},
		Location: []models.Location{
			{
				Id: 10,
			},
		},
		Action: "deliver",
		Notes: []models.Notes{
			{
				Note: "test",
			},
		},
		Attachments: nil,
		Assignee: models.Assignee{
			EmployeeId: 1,
			Username:   "test",
			AutoAssign: false,
		},
		DueBy: time.Now().Add(time.Hour * 24),
	}

	switch *job.Department {
	case "Housekeeping":
		validJobItems := []string{"Blanket", "Sheets", "Mattress"}
		if job.JobItem != nil && contains(validJobItems, *job.JobItem) {
			if checkHousekeepingLocation(job.Locations) {
				var roomLocations []string
				for _, loc := range job.Locations {
					if strings.Contains(strings.ToLower(loc), "room") {
						roomLocations = append(roomLocations, loc)
					}
				}

				var floorLocations []string
				for _, loc := range job.Locations {
					if strings.Contains(strings.ToLower(loc), "floor") {
						floorLocations = append(floorLocations, loc)
					}
				}

				// TODO: The swagger documentation and the exercise documentation are confusing
			} else {
				return nil, fmt.Errorf("invalid location"), http.StatusBadRequest
			}
		}
	case "Engineering":
		if job.JobItem != nil {
			if job.Locations != nil && len(job.Locations) > 0 {
				// TODO: The swagger documentation and the exercise documentation are confusing
			} else {
				return nil, fmt.Errorf("at least one location is required for Engineering department"), http.StatusBadRequest
			}
		} else {
			return nil, fmt.Errorf("job item is required for Engineering department"), http.StatusBadRequest
		}
	case "Room Service":
		if job.JobItem != nil {
			if job.Locations != nil && len(job.Locations) > 0 {
				// TODO: The swagger documentation and the exercise documentation are confusing
			} else {
				return nil, fmt.Errorf("at least one location is required for Room Service department"), http.StatusBadRequest
			}
		} else {
			return nil, fmt.Errorf("job item is required for Room Service department"), http.StatusBadRequest
		}
	}

	resp, err := s.api.CreateJob(newJob)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return resp, nil, http.StatusCreated
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func checkHousekeepingLocation(locations []string) bool {
	for _, loc := range locations {
		if strings.Contains(strings.ToLower(loc), "floor") || strings.Contains(strings.ToLower(loc), "room") {
			return true
		}
	}
	return false
}
