package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Department struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type PageInfo struct {
	TotalCount  int  `json:"totalCount"`
	EndCursor   int  `json:"endCursor"`
	HasNextPage bool `json:"hasNextPage"`
}

type Location struct {
	Id             int              `json:"id"`
	Name           *string           `json:"name,omitempty"`
	DisplayName    *string           `json:"displayName,omitempty"`
	ParentLocation *LocationSimplify `json:"parentLocation,omitempty"`
	LocationType   *LocationType     `json:"locationType,omitempty"`
}

type LocationSimplify struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type LocationType struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

type JobItem struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

type PagedResponse[T any] struct {
	PageInfo PageInfo `json:"pageInfo"`
	Items    []T      `json:"items"`
}

type Departments = PagedResponse[Department]
type Locations = PagedResponse[Location]
type LocationTypes = PagedResponse[LocationType]
type JobItems = PagedResponse[JobItem]
type Jobs = PagedResponse[Job]

type Job struct {
	Id          int        `json:"id,omitempty"`
	DisplayName string     `json:"displayName,omitempty"`
	Type        string     `json:"type,omitempty"`
	Priority    string     `json:"priority"`
	Action      string     `json:"action"`
	Item        Item       `json:"item"`
	Department  Department `json:"department"`
	Role        Roles      `json:"role"`
	Roles       []Roles    `json:"roles,omitempty"`
	Location    []Location `json:"location"`
	Notes       []Notes    `json:"notes"`
	Attachments []string   `json:"attachments,omitempty"`
	Assignee    Assignee   `json:"assignee"`
	DueBy       time.Time  `json:"dueBy"`
}

type Item struct {
	Name string `json:"name"`
}

type Roles struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Notes struct {
	Id   int    `json:"id,omitempty"`
	Note string `json:"note"`
}

type Assignee struct {
	Id         int    `json:"id,omitempty"`
	EmployeeId int    `json:"employeeId"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	AutoAssign bool   `json:"autoAssign"`
}

type CreateJobRequest struct {
	Description *string  `json:"description,omitempty"`
	Department  *string  `json:"department"`
	JobItem     *string  `json:"job_item"`
	Locations   []string `json:"locations"`
}

func (r *CreateJobRequest) UnmarshalJSON(data []byte) error {
	type Alias CreateJobRequest
	aux := &struct {
		Locations []string `json:"locations"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Department == nil {
		return errors.New("department is required")
	}

	r.Department = aux.Department

	if aux.JobItem == nil {
		return errors.New("job_item is required")
	}

	r.JobItem = aux.JobItem

	if aux.Locations == nil {
		return errors.New("locations is required")
	}

	r.Locations = aux.Locations
	r.Description = aux.Description

	return nil
}
