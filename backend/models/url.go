package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	OriginalURL string  `gorm:"not null" json:"original_url"`
	ShortCode   string  `gorm:"uniqueIndex;not null" json:"short_code"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	Clicks      []Click `gorm:"foreignKey:URLID" json:"-"` // relasi, jangan ikut ke response JSON
	ClickCount  int64   `gorm:"-" json:"click_count"`      // gorm:"-" = tidak disimpan di DB, dihitung manual
}

// ShortURL: buat endpoint nanti
func (u *URL) ShortURL(baseURL string) string {
	return baseURL + "/" + u.ShortCode
}
