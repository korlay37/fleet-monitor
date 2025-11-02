package models

import "time"

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
	UploadSum   int
	UploadCount int
}

type StatsResponse struct {
	Uptime        float64 `json:"uptime"`
	AvgUploadTime string  `json:"avg_upload_time"`
}
