package services

import (
	"net/http"
	"testing"

	"optii/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type JobRepositoryMock struct {
	mock.Mock
}

func (m *JobRepositoryMock) GetDepartment(id int) (*models.Department, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *JobRepositoryMock) GetDepartments(displayName string, first, next int) (*models.Departments, error) {
	args := m.Called(displayName, first, next)
	return args.Get(0).(*models.Departments), args.Error(1)
}

func (m *JobRepositoryMock) GetLocation(locationId int) (*models.Location, error) {
	args := m.Called(locationId)
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *JobRepositoryMock) GetLocations(params map[string]string) (*models.Locations, error) {
	args := m.Called(params)
	return args.Get(0).(*models.Locations), args.Error(1)
}

func (m *JobRepositoryMock) GetLocationTypes() (*models.LocationTypes, error) {
	args := m.Called()
	return args.Get(0).(*models.LocationTypes), args.Error(1)
}

func (m *JobRepositoryMock) GetLocationType(locationTypeId int) (*models.LocationType, error) {
	args := m.Called(locationTypeId)
	return args.Get(0).(*models.LocationType), args.Error(1)
}

func (m *JobRepositoryMock) GetJobItem(jobItemId int) (*models.JobItem, error) {
	args := m.Called(jobItemId)
	return args.Get(0).(*models.JobItem), args.Error(1)
}

func (m *JobRepositoryMock) GetJobItems(first int, next int, displayName string) (*models.JobItems, error) {
	args := m.Called(first, next, displayName)
	return args.Get(0).(*models.JobItems), args.Error(1)
}

func (m *JobRepositoryMock) GetJob(jobId int) (*models.Job, error) {
	args := m.Called(jobId)
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *JobRepositoryMock) GetJobs(params map[string]string) (*models.Jobs, error) {
	args := m.Called(params)
	return args.Get(0).(*models.Jobs), args.Error(1)
}

func (m *JobRepositoryMock) CreateJob(jobData *models.Job) (*models.Job, error) {
	args := m.Called(jobData)
	return args.Get(0).(*models.Job), args.Error(1)
}

func TestCreateJob(t *testing.T) {
	mockRepo := new(JobRepositoryMock)
	service := jobService{api: mockRepo}

	var desc, depart, jobItem = "test", "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		JobItem:     &jobItem,
		Locations:   location,
	}

	dep := &models.Departments{}
	mockRepo.On("GetDepartments", mock.Anything, mock.Anything, mock.Anything).Return(dep, nil).Once()

	loc := &models.Locations{}
	mockRepo.On("GetLocations", mock.Anything).Return(loc, nil).Once()

	item := &models.JobItems{}
	mockRepo.On("GetJobItems", mock.Anything, mock.Anything, mock.Anything).Return(item, nil).Once()

	job := &models.Job{}
	mockRepo.On("CreateJob", mock.Anything).Return(job, nil).Once()

	result, err, httpStatusCode := service.CreateJob(&body)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusCreated, httpStatusCode)

	mockRepo.AssertExpectations(t)
}

func TestCreateJobWithNilDepartment(t *testing.T) {
	mockRepo := new(JobRepositoryMock)
	service := jobService{api: mockRepo}

	var desc, depart, jobItem = "test", "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		JobItem:     &jobItem,
		Locations:   location,
	}

	var dep *models.Departments
	mockRepo.On("GetDepartments", mock.Anything, mock.Anything, mock.Anything).Return(dep, nil).Once()

	loc := &models.Locations{}
	mockRepo.On("GetLocations", mock.Anything).Return(loc, nil).Once()

	item := &models.JobItems{}
	mockRepo.On("GetJobItems", mock.Anything, mock.Anything, mock.Anything).Return(item, nil).Once()

	_, err, httpStatusCode := service.CreateJob(&body)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, httpStatusCode)

	mockRepo.AssertExpectations(t)
}

func TestCreateJobWithNilLocations(t *testing.T) {
	mockRepo := new(JobRepositoryMock)
	service := jobService{api: mockRepo}

	var desc, depart, jobItem = "test", "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		JobItem:     &jobItem,
		Locations:   location,
	}

	dep := &models.Departments{}
	mockRepo.On("GetDepartments", mock.Anything, mock.Anything, mock.Anything).Return(dep, nil).Once()

	var loc *models.Locations
	mockRepo.On("GetLocations", mock.Anything).Return(loc, nil).Once()

	item := &models.JobItems{}
	mockRepo.On("GetJobItems", mock.Anything, mock.Anything, mock.Anything).Return(item, nil).Once()

	_, err, httpStatusCode := service.CreateJob(&body)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, httpStatusCode)

	mockRepo.AssertExpectations(t)
}

func TestCreateJobWithNilJobItem(t *testing.T) {
	mockRepo := new(JobRepositoryMock)
	service := jobService{api: mockRepo}

	var desc, depart, jobItem = "test", "test", "test"
	var location []string
	location = append(location, "test")

	body := models.CreateJobRequest{
		Description: &desc,
		Department:  &depart,
		JobItem:     &jobItem,
		Locations:   location,
	}

	dep := &models.Departments{}
	mockRepo.On("GetDepartments", mock.Anything, mock.Anything, mock.Anything).Return(dep, nil).Once()

	loc := &models.Locations{}
	mockRepo.On("GetLocations", mock.Anything).Return(loc, nil).Once()

	var item *models.JobItems
	mockRepo.On("GetJobItems", mock.Anything, mock.Anything, mock.Anything).Return(item, nil).Once()

	_, err, httpStatusCode := service.CreateJob(&body)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, httpStatusCode)

	mockRepo.AssertExpectations(t)
}