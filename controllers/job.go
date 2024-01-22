package controllers

import (
	"net/http"

	"optii/models"
	"optii/services"
	"optii/utils"

	"github.com/gin-gonic/gin"
)

type JobController interface {
	Create(c *gin.Context)
}

type jobController struct {
	JobService services.JobService
}

func NewJobController(service services.JobService) JobController {
	return &jobController{
		JobService: service,
	}
}

// Create Job godoc
// @Summary Create an job
// @Description create new job
// @Tags job
// @Accept  json
// @Produce  json
// @Param job body models.CreateJobRequest true "Create Job"
// @Success 201 {object} models.CreateJobRequest
// @Failure 400 {object} utils.Response
// @Router /job [post]
func (ac *jobController) Create(c *gin.Context) {
	var job models.CreateJobRequest

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
		return
	}

	createdJob, err, httpStatus := ac.JobService.CreateJob(&job)
	if err != nil {
		c.JSON(httpStatus, utils.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdJob)
}
