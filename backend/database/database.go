package database

import (
	"log"
	"url-shortener/config"
	"url-shortener/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal koneksi database: ", err)
	}

	// auto migrate: buat label kalo belum ada
	err = DB.AutoMigrate(&models.User{}, &models.URL{}, &models.Click{})
	if err != nil {
		log.Fatal("Gagal migrate database: ", err)
	}

	log.Println("Database siap")
}
