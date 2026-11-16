package models

import "gorm.io/gorm"

type Click struct {
	gorm.Model
	URLID     uint   `gorm:"not null" json:"url_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Region    string `json:"region"` // negara dari IP, best-effort (boleh kosong)
}
