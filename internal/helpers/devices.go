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
	//NOTE: EVEN THOUGH THIS MATCHES THE EXPECTED RESULTS IN RESULTS.TXT, I THINK IT IS NOT CORRECT...
	// THE FORMULA ACCOUNTS FOR THE MINUTES BETWEEN THE FIRST AND LAST HEARTBEAT, NOT THE TOTAL HEARTBEATS
	// THAT SHOULD HAVE BEEN RECEIVED. EG:
	// FIRST HEARTBEAT:  00:00:00
	// LAST HEARTBEAT: 00:02:00
	// MINUTES BETWEEN: 2 BUT SHOULD HAVE ALSO RECEIVED HEARTBEAT AT 00:01:00.
	// uptime := (float64(len(heartbeats)) / (heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes())) * 100
	// BELOW CODE WOULD TAKE INTO ACCOUNT THE NUMBER OF HEARTBEATS THAT SHOULD HAVE BEEN RECEIVED.
	uptime := (float64(len(heartbeats)) / (heartbeats[len(heartbeats)-1].Sub(heartbeats[0]).Minutes() + 1)) * 100
	return uptime
}

func CalculateAverageUploadTime(uploadTimes []int) string {
	if len(uploadTimes) == 0 {
		return "0m0.000000000s"
	}
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
