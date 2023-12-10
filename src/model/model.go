package model

import "gorm.io/gorm"

type SaveLog struct {
	gorm.Model
	Request  string `json:"request"`
	Response string `json:"response"`
}
