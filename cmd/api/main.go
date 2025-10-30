package main

import (
	"github.com/gin-gonic/gin"
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
	devicesData, err := services.GetDevicesFromFile()
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
	router.Run("localhost:6733")
	sugar.Infow("API started", "port", "6733")
}
