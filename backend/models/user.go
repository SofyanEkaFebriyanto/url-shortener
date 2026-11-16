package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // nambah ID, CreatedAt, UpdatedAt, DeletedAt otomatis
	Username   string `gorm:"uniqueIndex;not null" json:"username"`
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	Password   string `gorm:"not null" json:"-"` //json:"-" artinya tidak dikirim ke response
	URLs       []URL  `gorm:"foreignKey:UserID"` // relasi:satu user punya banyak URL
}
