package main

import (
	"log/slog"

	"optii/config"
	"optii/docs"
	_ "optii/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Optii API
// @description This is a Optii API server.
// @BasePath /api/v1
// @version v1
// @host localhost:8080
func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	r := gin.Default()

	infra := config.NewInfra()
	controller := infra.SetupJobController()

	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	jobs := r.Group("/jobs")
	jobs.POST("", controller.Create)

	r.Run()
}
