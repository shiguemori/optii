package config

import "optii/controllers"

func (i *Infra) SetupJobController() controllers.JobController {
	return controllers.NewJobController(i.SetupJobService())
}