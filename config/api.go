package config

import (
	"os"

	"optii/api"
)

func (i *Infra) SetupOptiiApi() api.OptiiApi {
	return api.NewOptiiApi(
		os.Getenv("OPTII_URL"),
		os.Getenv("OPTII_CLIENT_ID"),
		os.Getenv("OPTII_CLIENT_SECRET"),
		os.Getenv("OPTII_AUTHETICATION_URL"),
	)
}
