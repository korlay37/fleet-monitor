package services

import (
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/korlay37/fleet-monitor/internal/helpers"
	"github.com/korlay37/fleet-monitor/internal/models"
)

var Sugar *zap.SugaredLogger
var DevicesMap = map[string]models.DeviceData{}
var devicesMapMutex = sync.RWMutex{}

func PostDeviceHeartbeat(context *gin.Context) {
	id := context.Param("device_id")
	var heartbeatRequest models.HeartbeatRequest
	Sugar.Infow("Received heartbeat request", "device_id", id)
	if err := context.ShouldBindJSON(&heartbeatRequest); err != nil {
		Sugar.Errorw("Invalid heartbeat request", "device_id", id, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	devicesMapMutex.Lock()
	device, exists := DevicesMap[id]
	if exists {
		device.Heartbeats = append(device.Heartbeats, heartbeatRequest.SentAt)
		DevicesMap[id] = device
		context.JSON(http.StatusNoContent, nil)
	} else {
		Sugar.Warnw("Heartbeat request from unknown device", "device_id", id)
		context.JSON(http.StatusNotFound, gin.H{"error": "Device id '" + id + "' not found"})
		return
	}
	devicesMapMutex.Unlock()
}

func PostDeviceStats(context *gin.Context) {
	id := context.Param("device_id")
	var statsRequest models.StatsRequest
	Sugar.Infow("Received Upload Stats request", "device_id", id)
	if err := context.ShouldBindJSON(&statsRequest); err != nil {
		Sugar.Errorw("Invalid Upload Stats request", "device_id", id, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	devicesMapMutex.Lock()
	device, exists := DevicesMap[id]
	if exists {
		// Keeping Upload times for future use, if needed.
		device.UploadTimes = append(device.UploadTimes, statsRequest.UploadTime)
		device.UploadSum += statsRequest.UploadTime
		device.UploadCount++
		DevicesMap[id] = device
		context.JSON(http.StatusNoContent, nil)
	} else {
		Sugar.Warnw("Upload Stats request from unknown device", "device_id", id)
		context.JSON(http.StatusNotFound, gin.H{"error": "Device id '" + id + "' not found"})
		return
	}
	devicesMapMutex.Unlock()
}
func GetDeviceStats(context *gin.Context) {
	id := context.Param("device_id")
	Sugar.Infow("Device Stats requested", "device_id", id)
	devicesMapMutex.RLock()
	device, exists := DevicesMap[id]
	if exists {
		uptime := helpers.CalculateUptime(device.Heartbeats)
		avgUploadTime := helpers.CalculateAverageUploadTime(device.UploadSum, device.UploadCount)
		context.JSON(http.StatusOK, models.StatsResponse{
			Uptime:        uptime,
			AvgUploadTime: avgUploadTime,
		})
	} else {
		Sugar.Warnw("Device stats requested but not found", "device_id", id)
		context.JSON(http.StatusNotFound, gin.H{"error": "Device id '" + id + "' not found"})
		return
	}
	devicesMapMutex.RUnlock()
}
