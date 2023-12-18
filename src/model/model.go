package model

import (
	"time"

	"gorm.io/gorm"
)

type RateLimitInfo struct {
	RequestCount int
	StartTime    time.Time
}

type SaveLog struct {
	gorm.Model
	Request  string `json:"request"`
	Response string `json:"response"`
}
