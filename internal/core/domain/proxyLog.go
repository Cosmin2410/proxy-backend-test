package domain

import "gorm.io/gorm"

type ProxyLog struct {
	gorm.Model
	Request  string `json:"request"`
	Response string `json:"response"`
}
