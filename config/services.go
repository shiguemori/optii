package config

import "optii/services"

func (i *Infra) SetupJobService() services.JobService {
	return services.NewJobService(i.SetupOptiiApi())
}