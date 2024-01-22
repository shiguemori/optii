package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"optii/models"
)

type OptiiApi interface {
	GetDepartment(id int) (*models.Department, error)
	GetDepartments(displayName string, first, next int) (*models.Departments, error)
	GetLocation(locationId int) (*models.Location, error)
	GetLocations(params map[string]string) (*models.Locations, error)
	GetLocationTypes() (*models.LocationTypes, error)
	GetLocationType(locationTypeId int) (*models.LocationType, error)
	GetJobItem(jobItemId int) (*models.JobItem, error)
	GetJobItems(first int, next int, displayName string) (*models.JobItems, error)
	GetJob(jobId int) (*models.Job, error)
	GetJobs(params map[string]string) (*models.Jobs, error)
	CreateJob(jobData *models.Job) (*models.Job, error)
}

type optiiApi struct {
	httpClient   *http.Client
	url          string
	clientSecret string
	clientId     string
	bearer       string
	authUrl      string
	bearerExpiry time.Time
}

func NewOptiiApi(url, clientId, clientSecret, authUrl string) *optiiApi {
	return &optiiApi{
		httpClient:   &http.Client{Timeout: time.Second * 30},
		url:          url,
		authUrl:      authUrl,
		clientSecret: clientSecret,
		clientId:     clientId,
	}
}

func (s *optiiApi) GetBearer() error {
	data := url.Values{}
	data.Set("client_id", s.clientId)
	data.Set("client_secret", s.clientSecret)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openapi")

	req, err := http.NewRequest("POST", s.authUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get bearer token, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if token, ok := result["access_token"].(string); ok {
		s.bearer = token
	} else {
		return fmt.Errorf("access token missing in response")
	}

	return nil
}

func (s *optiiApi) doRequest(req *http.Request, retry bool) (*http.Response, error) {
	if s.bearer == "" || time.Now().After(s.bearerExpiry) {
		if err := s.GetBearer(); err != nil {
			return nil, err
		}
	}

	req.Header.Add("Authorization", "Bearer "+s.bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusUnauthorized:
		if retry {
			resp.Body.Close()
			return s.doRequest(req, false)
		}
		return nil, fmt.Errorf("autenticação falhou, não é possível repetir a solicitação")
	default:
		resp.Body.Close()
		return nil, fmt.Errorf("resposta falhou, código de status: %d", resp.StatusCode)
	}
}

func (s *optiiApi) GetDepartment(departmentId int) (*models.Department, error) {
	url := fmt.Sprintf("%s/api/v1/departments/%d", s.url, departmentId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var department models.Department
	if err := json.NewDecoder(resp.Body).Decode(&department); err != nil {
		return nil, err
	}

	return &department, nil
}

func (s *optiiApi) GetDepartments(displayName string, first, next int) (*models.Departments, error) {
	queryParams := url.Values{}
	if displayName != "" {
		queryParams.Add("displayName", displayName)
	}
	if first > 0 {
		queryParams.Add("first", strconv.Itoa(first))
	}
	if next > 0 {
		queryParams.Add("next", strconv.Itoa(next))
	}
	fullURL := fmt.Sprintf("%s/api/v1/departments?%s", s.url, queryParams.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response failed, status code: %d", resp.StatusCode)
	}

	var departments models.Departments
	if err := json.NewDecoder(resp.Body).Decode(&departments); err != nil {
		return nil, err
	}

	return &departments, nil
}

func (s *optiiApi) GetLocation(locationId int) (*models.Location, error) {
	url := fmt.Sprintf("%s/api/v1/locations/%d", s.url, locationId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var location models.Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *optiiApi) GetLocations(params map[string]string) (*models.Locations, error) {
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	fullURL := fmt.Sprintf("%s/api/v1/locations?%s", s.url, queryParams.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting locations, status code: %d", resp.StatusCode)
	}

	var locations models.Locations
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, err
	}

	return &locations, nil
}

func (s *optiiApi) GetLocationTypes() (*models.LocationTypes, error) {
	url := fmt.Sprintf("%s/api/v1/locationTypes", s.url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting location types, status code: %d", resp.StatusCode)
	}

	var locationTypes models.LocationTypes
	if err := json.NewDecoder(resp.Body).Decode(&locationTypes); err != nil {
		return nil, err
	}

	return &locationTypes, nil
}

func (s *optiiApi) GetLocationType(locationTypeId int) (*models.LocationType, error) {
	url := fmt.Sprintf("%s/api/v1/locationTypes/%d", s.url, locationTypeId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting location type, status code: %d", resp.StatusCode)
	}

	var locationType models.LocationType
	if err := json.NewDecoder(resp.Body).Decode(&locationType); err != nil {
		return nil, err
	}

	return &locationType, nil
}

func (s *optiiApi) GetJobItem(jobItemId int) (*models.JobItem, error) {
	url := fmt.Sprintf("%s/api/v1/jobitems/%d", s.url, jobItemId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting job item, status code: %d", resp.StatusCode)
	}

	var jobItem models.JobItem
	if err := json.NewDecoder(resp.Body).Decode(&jobItem); err != nil {
		return nil, err
	}

	return &jobItem, nil
}

func (s *optiiApi) GetJobItems(first int, next int, displayName string) (*models.JobItems, error) {
	queryParams := url.Values{}
	if first > 0 {
		queryParams.Add("first", strconv.Itoa(first))
	}
	if next > 0 {
		queryParams.Add("next", strconv.Itoa(next))
	}
	if displayName != "" {
		queryParams.Add("displayName", displayName)
	}
	fullURL := fmt.Sprintf("%s/api/v1/jobitems?%s", s.url, queryParams.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting job items, status code: %d", resp.StatusCode)
	}

	var jobItems models.JobItems
	if err := json.NewDecoder(resp.Body).Decode(&jobItems); err != nil {
		return nil, err
	}

	return &jobItems, nil
}

func (s *optiiApi) GetJob(jobId int) (*models.Job, error) {
	url := fmt.Sprintf("%s/api/v1/jobs/%d", s.url, jobId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var job models.Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *optiiApi) GetJobs(params map[string]string) (*models.Jobs, error) {
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	fullURL := fmt.Sprintf("%s/api/v1/jobs?%s", s.url, queryParams.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobs models.Jobs
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return &jobs, nil
}

func (s *optiiApi) CreateJob(jobData *models.Job) (*models.Job, error) {
	jsonData, err := json.Marshal(jobData)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/jobs", s.url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := s.doRequest(req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var job models.Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}
