package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"optii/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock service
type MockJobsService struct {
	mock.Mock
}

func (m *MockJobsService) CreateJob(job *models.CreateJobRequest) (*models.Job, error, int) {
	args := m.Called(job)
	return args.Get(0).(*models.Job), args.Error(1), args.Int(2)
}

func TestCreateJob(t *testing.T) {
	mockService := new(MockJobsService)
	var desc, depart, jobItem = "test", "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		JobItem:     &jobItem,
		Locations:   location,
	}

	job := &models.Job{}
	mockService.On("CreateJob", &body).Return(job, nil, http.StatusCreated)

	controller := NewJobController(mockService)

	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	requestBodyBuffer := bytes.NewBuffer(requestBodyBytes)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("POST", "/jobs", requestBodyBuffer)

	controller.Create(context)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreateFailDepartmentJob(t *testing.T) {
	mockService := new(MockJobsService)
	var desc, jobItem = "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		JobItem:     &jobItem,
		Locations:   location,
	}

	job := &models.Job{}
	mockService.On("CreateJob", &body).Return(job, nil, http.StatusBadRequest)

	controller := NewJobController(mockService)

	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	requestBodyBuffer := bytes.NewBuffer(requestBodyBytes)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("POST", "/jobs", requestBodyBuffer)

	controller.Create(context)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "{\"message\":\"department is required\"}", recorder.Body.String())
}


func TestCreateFailItemJob(t *testing.T) {
	mockService := new(MockJobsService)
	var desc, depart = "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		Locations:   location,
	}

	job := &models.Job{}
	mockService.On("CreateJob", &body).Return(job, nil, http.StatusBadRequest)

	controller := NewJobController(mockService)

	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	requestBodyBuffer := bytes.NewBuffer(requestBodyBytes)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("POST", "/jobs", requestBodyBuffer)

	controller.Create(context)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "{\"message\":\"job_item is required\"}", recorder.Body.String())
}


func TestCreateFailLocationJob(t *testing.T) {
	mockService := new(MockJobsService)
	var desc, depart, jobItem = "test", "test", "test"

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		JobItem:     &jobItem,
	}

	job := &models.Job{}
	mockService.On("CreateJob", &body).Return(job, nil, http.StatusBadRequest)

	controller := NewJobController(mockService)

	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	requestBodyBuffer := bytes.NewBuffer(requestBodyBytes)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("POST", "/jobs", requestBodyBuffer)

	controller.Create(context)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "{\"message\":\"locations is required\"}", recorder.Body.String())
}
