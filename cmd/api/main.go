package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/korlay37/fleet-monitor/internal/models"
	"github.com/korlay37/fleet-monitor/internal/services"
)

var sugar *zap.SugaredLogger

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar = logger.Sugar()
	services.Sugar = sugar

	err := godotenv.Load()
	if err != nil {
		sugar.Errorw("Error loading .env file, continuing without it", "error", err)
	}
	devicesFile := os.Getenv("DEVICES_FILE")
	if devicesFile == "" {
		sugar.Infow("DEVICES_FILE environment variable not set, defaulting to devices.csv")
		devicesFile = "devices.csv"
	}
	port := os.Getenv("PORT")
	if port == "" {
		sugar.Infow("PORT environment variable not set, defaulting to 6733")
		port = "6733"
	}
	devicesData, err := services.GetDevicesFromFile(devicesFile)
	if err != nil {
		sugar.Errorw("Error reading devices file", "error", err)
		return
	} else {
		sugar.Infow("Loaded devices from file", "count", len(devicesData))
		for _, device := range devicesData {
			services.DevicesMap[device] = models.DeviceData{
				DeviceID: device,
			}
		}
	}
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/devices/:device_id/heartbeat", services.PostDeviceHeartbeat)
		api.POST("/devices/:device_id/stats", services.PostDeviceStats)
		api.GET("/devices/:device_id/stats", services.GetDeviceStats)
	}
	router.Run(":" + port)
	sugar.Infow("API started", "port", port)
}
