package helpers

import (
	"fmt"
	"strings"
	"time"
)

func CleanDevicesData(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && trimmed != "device_id" {
			result = append(result, trimmed)
		}
	}
	return result
}

func CalculateUptime(heartbeats []time.Time) float64 {
	if len(heartbeats) == 0 {
		return 0.0
	} else if len(heartbeats) < 2 {
		return 100.0
	}
	// NOTE: Check README.md for more information about the uptime formula.
	// uptime := (float64(len(heartbeats)) / (heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes())) * 100
	uptime := (float64(len(heartbeats)) / (heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes() + 1)) * 100
	return uptime
}

func CalculateAverageUploadTime(uploadSum int, uploadCount int) string {
	if uploadSum == 0 {
		return "0m0.000000000s"
	}
	timeDuration := float64(uploadSum) / float64(uploadCount)
	totalSeconds := timeDuration / 1000000000.0
	minutes := int(totalSeconds / 60.0)
	seconds := totalSeconds - float64(minutes*60)
	return fmt.Sprintf("%dm%.9fs", minutes, seconds)
}
