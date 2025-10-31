package helpers_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/korlay37/fleet-monitor/internal/helpers"
)

func TestCleanDevicesData(t *testing.T) {
	type testCase struct {
		lines    []string
		expected []string
	}
	tests := []testCase{
		{
			lines: []string{
				"device_id",
				"60-6b-44-84-dc-64",
				"b4-45-52-a2-f1-3c",
				"",
			},
			expected: []string{
				"60-6b-44-84-dc-64",
				"b4-45-52-a2-f1-3c",
			},
		},
		{
			lines: []string{
				"device_id",
				"",
				"60-6b-44-84-dc-64",
				"",
				"b4-45-52-a2-f1-3c",
				"",
			},
			expected: []string{
				"60-6b-44-84-dc-64",
				"b4-45-52-a2-f1-3c",
			},
		},
		{
			lines: []string{
				"",
				"device_id",
				"60-6b-44-84-dc-64",
				"",
				"",
			},
			expected: []string{
				"60-6b-44-84-dc-64",
			},
		},
	}
	for _, test := range tests {
		result := helpers.CleanDevicesData(test.lines)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestCalculateUptime(t *testing.T) {
	type testCase struct {
		heartbeats []time.Time
		expected   float64
	}
	timeNow := time.Now()
	tests := []testCase{
		{
			heartbeats: []time.Time{},
			expected:   0.0,
		},
		{
			heartbeats: []time.Time{
				timeNow,
			},
			expected: 100.0,
		},
		{
			heartbeats: []time.Time{
				timeNow,
				timeNow.Add((time.Minute) * 2),
			},
			expected: 66.66666666666666,
		},
		{
			heartbeats: []time.Time{
				timeNow,
				timeNow.Add((time.Minute) * 3),
			},
			expected: 50.0,
		},
		{
			heartbeats: []time.Time{
				timeNow,
				timeNow.Add((time.Minute) * 2),
				timeNow.Add((time.Minute) * 4),
			},
			expected: 60.0,
		},
		{
			heartbeats: []time.Time{
				timeNow,
				timeNow.Add((time.Minute) * 1),
				timeNow.Add((time.Minute) * 2),
				timeNow.Add((time.Minute) * 3),
				timeNow.Add((time.Minute) * 4),
			},
			expected: 100.0,
		},
	}
	for _, test := range tests {
		result := helpers.CalculateUptime(test.heartbeats)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}
