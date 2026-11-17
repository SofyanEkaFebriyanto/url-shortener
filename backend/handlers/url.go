package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"url-shortener/config"
	"url-shortener/database"
	"url-shortener/models"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

type CreateURLInput struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomSlug  string `json:"custom_slug"`
}

// recordClick: jalankan di goroutine supaya redirect tidak lambat
func recordClick(urlID uint, ip, userAgent string) {
	click := models.Click{
		URLID:     urlID,
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&click)
}

// lookupRegion menebak negara dari IP lewat ip-api.com (gratis tanpa api key, limit 45 request/menit). Gagal atau IP lokal/privat -> string kosong, bukan error: tracking tidak boleh menggagalkan redirect.

func lookupRegion(ip string) string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?fields=country")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var body struct {
		Country string `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return ""
	}
	return body.Country
}

var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validSlug(s string) bool {
	return slugPattern.MatchString(s)
}

func CreateURL(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	cfg := config.Load()

	var input CreateURLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi custom slug: hanya huruf, angka, - dan _. slug berisi spasi atau "/" akan merusak URL redirect, jadi tolak di awal
	if input.CustomSlug != "" && !validSlug(input.CustomSlug) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug hanya boleh huruf, angka, -, dan _"})
		return
	}

	// tentukan short code

	shortCode := input.CustomSlug
	if shortCode == "" {
		// generate random sampai unik, maksimal 10 kali percobaan. loop tanpa batas berbahaya kalau ada masalah di DB.

		// catatan: cek-lalu-simpan ini punya race kecil - dua request bersamaan bisa dapat kode yang sama, lalu create di bawah gagal karena unique index dan user, dapat 500. untuk skala portofolio ini bisa diterima; di production tambahkan retry saat create gagal karena duplikat.
		for i := 0; i < 10; i++ {
			shortCode = utils.GenerateShortCode()
			var existing models.URL
			if database.DB.Where("short_code = ?", shortCode).First(&existing).Error != nil {
				break // tidak ada yang sama, bisa dipakai
			}
			if i == 9 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat kode unik, coba lagi"})
				return
			}
		}
	} else {
		// cek apakah custom slug sudah dipakai
		var existing models.URL
		if database.DB.Where("short_code = ?", shortCode).First(&existing).Error == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "slug sudah dipakai"})
			return
		}
	}

	url := models.URL{
		OriginalURL: input.OriginalURL,
		ShortCode:   shortCode,
		UserID:      userID,
	}

	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           url.ID,
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"short_url":    url.ShortURL(cfg.BaseURL),
		"created_at":   url.CreatedAt,
	})
}

func listURLs(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	cfg := config.Load()

	var urls []models.URL
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&urls)

	// hitung clict count tiap URL
	type URLResponse struct {
		ID          uint   `json:"id"`
		OriginalURL string `json:"original_url"`
		ShortCode   string `json:"short_code"`
		ShortURL    string `json:"short_url"`
		ClickCount  int64  `json:"click_count"`
		CreatedAt   string `json:"created_at"`
	}

	var result []URLResponse
	for _, u := range urls {
		var count int64
		database.DB.Model(&models.Click{}).Where("url_id = ?", u.ID).Count(&count)
		result = append(result, URLResponse{
			ID:          u.ID,
			OriginalURL: u.OriginalURL,
			ShortCode:   u.ShortCode,
			ShortURL:    u.ShortURL(cfg.BaseURL),
			ClickCount:  count,
			CreatedAt:   u.CreatedAt.Format("2006-01-02"),
		})
	}

	if result == nil {
		result = []URLResponse{} // return array kosong, bukan null
	}

	c.JSON(http.StatusOK, result)
}

func DeleteURL(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	urlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	var url models.URL
	if err := database.DB.Where("id = ? AND user_id = ?", urlID, userID).First(&url).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL tidak ditemukan"})
		return
	}

	// hapus click dulu, baru url nya. url dihapus HARD delete (unscope), bukan soft delete: gorm.Model punya DeletedAt, dan baris yang soft-deleted MASIH menempati unique index short_code.
	// kalau pake soft delete, slug yang sudah dihapus tidak bisa dipakai lagi (cek di CreatedURL lolos, tapi insert gagal 500).
	database.DB.Where("url_id = ?", url.ID).Delete(&models.Click{})
	database.DB.Unscoped().Delete(&url)

	c.JSON(http.StatusOK, gin.H{"message": "URL berhasil dihapus"})
}

func Redirect(c *gin.Context) {
	code := c.Param("code")

	var url models.URL
	if err := database.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL tidak ditemukan"})
		return
	}

	// catata klik secara async, tidak bikin redirect lambat
	go recordClick(url.ID, c.ClientIP(), c.GetHeader("User-Agent"))

	// 302 (Found), bukan 301 di cache permanen oleh browser: kalau short URL di hapus lalu kodenya dipakai ullang, browser tetap membuka tujuan lama dari cache. 302 tidak di cache permanen.
	c.Redirect(http.StatusFound, url.OriginalURL)
}
