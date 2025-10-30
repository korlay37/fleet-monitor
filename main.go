package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HeartbeatRequest struct {
	SentAt time.Time `json:"sent_at"`
}

type StatsRequest struct {
	SentAt     time.Time `json:"sent_at"`
	UploadTime int       `json:"upload_time"`
}

type DeviceData struct {
	DeviceID    string
	Heartbeats  []time.Time
	UploadTimes []int
}

type StatsResponse struct {
	Uptime        float64 `json:"uptime"`
	AvgUploadTime string  `json:"avg_upload_time"`
}

var sugar *zap.SugaredLogger
var devices = map[string]DeviceData{}
var devicesMutex = sync.RWMutex{}

func cleanDevicesData(lines []string) []string {
	var result []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if i > 0 && trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getDevicesFromFile() []string {
	devices, err := os.ReadFile("devices.csv")
	if err != nil {
		sugar.Errorw("Error reading devices file", "error", err)
		return []string{}
	}
	lines := strings.Split(string(devices), "\n")
	devicesData := cleanDevicesData(lines)
	sugar.Infow("Loaded devices from file", "count", len(devicesData))
	return devicesData
}

func postDeviceHeartbeat(context *gin.Context) {
	id := context.Param("device_id")
	var heartbeatRequest HeartbeatRequest
	sugar.Infow("Received heartbeat request", "device_id", id)
	if err := context.ShouldBindJSON(&heartbeatRequest); err != nil {
		sugar.Errorw("Invalid heartbeat request", "device_id", id, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	devicesMutex.Lock()
	device, exists := devices[id]
	if exists {
		device.Heartbeats = append(device.Heartbeats, heartbeatRequest.SentAt)
		devices[id] = device
		context.JSON(http.StatusNoContent, nil)
	} else {
		sugar.Warnw("Heartbeat request from unknown device", "device_id", id)
		context.JSON(http.StatusNotFound, gin.H{"error": "Device id '" + id + "' not found"})
		return
	}
	devicesMutex.Unlock()
}

func postDeviceStats(context *gin.Context) {
	id := context.Param("device_id")
	var statsRequest StatsRequest
	sugar.Infow("Received Upload Stats request", "device_id", id)
	if err := context.ShouldBindJSON(&statsRequest); err != nil {
		sugar.Errorw("Invalid Upload Stats request", "device_id", id, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	devicesMutex.Lock()
	device, exists := devices[id]
	if exists {
		device.UploadTimes = append(device.UploadTimes, statsRequest.UploadTime)
		devices[id] = device
		context.JSON(http.StatusNoContent, nil)
	} else {
		sugar.Warnw("Upload Stats request from unknown device", "device_id", id)
		context.JSON(http.StatusNotFound, gin.H{"error": "Device id '" + id + "' not found"})
		return
	}
	devicesMutex.Unlock()
}
func getDeviceStats(context *gin.Context) {
	id := context.Param("device_id")
	sugar.Infow("Device Stats requested", "device_id", id)
	devicesMutex.RLock()
	device, exists := devices[id]
	if exists {
		uptime := calculateUptime(device.Heartbeats)
		avgUploadTime := calculateAverageUploadTime(device.UploadTimes)
		context.JSON(http.StatusOK, StatsResponse{
			Uptime:        uptime,
			AvgUploadTime: avgUploadTime,
		})
	} else {
		sugar.Warnw("Device stats requested but not found", "device_id", id)
		context.JSON(http.StatusNotFound, gin.H{"error": "Device id '" + id + "' not found"})
		return
	}
	devicesMutex.RUnlock()
}

func calculateUptime(heartbeats []time.Time) float64 {
	if len(heartbeats) < 2 {
		return 100.0
	}
	uptime := (float64(len(heartbeats)) / heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes()) * 100
	return uptime
}

func calculateAverageUploadTime(uploadTimes []int) string {
	var sum float64
	for _, value := range uploadTimes {
		sum += float64(value)
	}
	timeDuration := sum / float64(len(uploadTimes))
	totalSeconds := timeDuration / 1000000000.0
	minutes := int(totalSeconds / 60.0)
	seconds := totalSeconds - float64(minutes*60)
	return fmt.Sprintf("%dm%.9fs", minutes, seconds)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar = logger.Sugar()

	for _, device := range getDevicesFromFile() {
		devices[device] = DeviceData{
			DeviceID: device,
		}
	}
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/devices/:device_id/heartbeat", postDeviceHeartbeat)
		api.POST("/devices/:device_id/stats", postDeviceStats)
		api.GET("/devices/:device_id/stats", getDeviceStats)
	}
	router.Run("localhost:6733")
	sugar.Infow("API started", "port", "6733")
}
