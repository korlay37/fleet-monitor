package services

import (
	"os"
	"strings"

	"github.com/korlay37/fleet-monitor/internal/helpers"
)

func GetDevicesFromFile(devicesFile string) ([]string, error) {
	devices, err := os.ReadFile(devicesFile)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(devices), "\n")
	devicesData := helpers.CleanDevicesData(lines)
	return devicesData, nil
}
