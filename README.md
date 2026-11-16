# MODUL PROJECT PORTFOLIO
## URL Shortener + Analytics Dashboard
### Stack: Go (Gin + GORM) + React (Vite) + SQLite + JWT

**Untuk:** Persiapan PKL — November 2026  
**Durasi:** 6 minggu  
**Output:** Project live, siap demo, siap dijelasin saat interview  

---

## DAFTAR ISI

- [BAB 0 — Persiapan Lingkungan](#bab-0)
- [BAB 1 — Go Crash Course](#bab-1)
- [BAB 2 — React Crash Course](#bab-2)
- [BAB 3 — Setup Project](#bab-3)
- [BAB 4 — Backend (Go + Gin)](#bab-4)
- [BAB 5 — Frontend (React + Vite)](#bab-5)
- [BAB 6 — Integrasi & Testing](#bab-6)
- [BAB 7 — Deploy](#bab-7)
- [BAB 8 — Interview Prep](#bab-8)

---

## BAB 0 — Persiapan Lingkungan {#bab-0}

### Install Go

```bash
# Download dari https://go.dev/dl/ — ambil versi terbaru 1.27.x
# (Go 1.21 sudah tidak dapat security fix; per Nov 2026 yang didukung
# hanya 1.26 dan 1.27, dan 1.27.4 rilis 10 Nov 2026)
# Linux:
wget https://go.dev/dl/go1.27.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.27.3.linux-amd64.tar.gz

# Tambah ke ~/.bashrc atau ~/.zshrc:
export PATH=$PATH:/usr/local/go/bin

source ~/.bashrc
go version  # harus muncul: go version go1.27.x
```

### Install Node.js

```bash
# Pakai nvm (Node Version Manager) — cara paling aman
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
source ~/.bashrc

nvm install 22
nvm use 22
node -v   # harus v22.x.x
npm -v    # harus 10.x.x
```

### Install VS Code Extensions

Buka VS Code, install ekstensi ini:
- `golang.go` — Go language support
- `dbaeumer.vscode-eslint` — ESLint
- `esbenp.prettier-vscode` — Prettier
- `bradlc.vscode-tailwindcss` — Tailwind IntelliSense
- `PKief.material-icon-theme` — icon biar enak diliat

### Install Tools Tambahan

```bash
# Git (kalau belum ada)
sudo apt install git -y   # Ubuntu/Debian
# atau
sudo dnf install git -y   # Fedora

# GCC (dibutuhkan SQLite driver Go)
sudo apt install gcc -y
# atau
sudo dnf install gcc -y

# Verifikasi semua
go version && node -v && npm -v && git --version && gcc --version
```

---

## BAB 1 — Go Crash Course {#bab-1}

> Tujuan: Cukup buat nulis backend project ini. Bukan buat jadi expert Go dulu.

### 1.1 Program Pertama

```go
package main  // setiap file Go punya package

import "fmt"  // import library standard

func main() { // titik masuk program
    fmt.Println("Hello, World!")
}
```

Jalankan: `go run main.go`

### 1.2 Variabel

```go
// Cara 1: explicit
var nama string = "Sofyan"
var umur int = 17

// Cara 2: shorthand (yang sering dipakai)
nama := "Sofyan"
umur := 17

// Cara 3: const
const PORT = "8080"
```

> Tiga cara di atas adalah ALTERNATIF, bukan urutan. Jangan tulis
> `var umur int = 17` lalu `umur := 17` di fungsi yang sama — itu
> compile error (`umur redeclared`), karena `:=` butuh minimal satu
> variabel baru di sisi kiri.

### 1.3 Tipe Data Penting

```go
string   // teks
int      // bilangan bulat
float64  // bilangan desimal
bool     // true / false
[]string // slice (array dinamis)
map[string]interface{}  // key-value, seperti object JS
```

### 1.4 Fungsi

```go
// Fungsi biasa
func tambah(a int, b int) int {
    return a + b
}

// Fungsi dengan 2 return value (ini idiom Go yang penting!)
// (snippet ini butuh `import "fmt"` di atasnya)
func bagi(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("tidak bisa bagi dengan 0")
    }
    return a / b, nil
}

// Cara pakai:
hasil, err := bagi(10, 2)
if err != nil {
    // tangani error
    fmt.Println("Error:", err)
    return
}
fmt.Println(hasil)
```

> **Pola `if err != nil`** adalah pola paling penting di Go. Hampir semua fungsi return 2 nilai: hasil + error. Kalau error nil artinya sukses.

### 1.5 Struct (Pengganti Class)

```go
// Definisi struct
type User struct {
    ID       uint   
    Username string 
    Email    string 
    Password string 
}

// Buat instance
user := User{
    Username: "sofyan",
    Email:    "sofyan@email.com",
}

// Akses field
fmt.Println(user.Username)

// Method di struct
func (u User) Greet() string {
    return "Halo, " + u.Username
}
```

### 1.6 Error Handling

```go
// Ini pola standar yang SELALU dipakai
hasil, err := someFunction()
if err != nil {
    // log error dan stop / return
    return fmt.Errorf("gagal: %w", err)
}
// lanjut pakai hasil
```

### 1.7 Goroutine (Async)

```go
// Jalankan fungsi di background (async, non-blocking)
go recordClick(urlID, ipAddress)

// Program utama lanjut tanpa nunggu recordClick selesai
// Ini yang kita pakai buat tracking klik tanpa delay redirect
```

### 1.8 HTTP Handler dengan Gin

```go
// Gin itu framework web — mirip Express.js kalau lo pernah denger
r := gin.Default()

// GET endpoint
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
})

// POST endpoint
r.POST("/users", func(c *gin.Context) {
    var body struct {
        Name string `json:"name"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, gin.H{"name": body.Name})
})

// URL parameter
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})

r.Run(":8080")
```

---

## BAB 2 — React Crash Course {#bab-2}

> Tujuan: Cukup buat bikin frontend project ini.

### 2.1 JSX

```jsx
// JSX = JavaScript + HTML. Ditulis di file .jsx
// Setiap komponen = fungsi yang return HTML (JSX)

function Greeting() {
    const name = "Sofyan"
    return (
        <div>
            <h1>Halo, {name}!</h1>  {/* kurung kurawal = ekspresi JS */}
        </div>
    )
}
```

### 2.2 Component & Props

```jsx
// Komponen = blok UI yang bisa dipakai ulang.
//
// Tailwind hanya meng-generate class yang tertulis LENGKAP di kode.
// `bg-${color}-500` TIDAK akan jadi CSS karena Tailwind tidak bisa
// menebak nilai `color` saat build — tombol jadi tidak berwarna.
// Solusinya: tulis semua class lengkap di sebuah map seperti ini.
const colorMap = {
    blue: 'bg-blue-500 hover:bg-blue-600',
    green: 'bg-green-500 hover:bg-green-600',
    red: 'bg-red-500 hover:bg-red-600',
}

function Button({ label, onClick, color = 'blue' }) {
    return (
        <button
            onClick={onClick}
            className={`px-4 py-2 ${colorMap[color]} text-white rounded`}
        >
            {label}
        </button>
    )
}

// Pemakaian:
<Button label="Klik Aku" onClick={() => alert("hi")} color="blue" />
```

### 2.3 useState — State Lokal

```jsx
import { useState } from 'react'

function Counter() {
    const [count, setCount] = useState(0)  // [nilai, setter]

    return (
        <div>
            <p>{count}</p>
            <button onClick={() => setCount(count + 1)}>+</button>
        </div>
    )
}
```

### 2.4 useEffect — Side Effects

```jsx
import { useState, useEffect } from 'react'

function UserList() {
    const [users, setUsers] = useState([])

    // Jalan sekali saat komponen muncul ([] = dependency kosong)
    useEffect(() => {
        fetch('/api/users')
            .then(res => res.json())
            .then(data => setUsers(data))
    }, [])

    return (
        <ul>
            {users.map(user => (
                <li key={user.id}>{user.name}</li>
            ))}
        </ul>
    )
}
```

### 2.5 Axios (HTTP Client)

```jsx
import axios from 'axios'

// GET
const response = await axios.get('/api/urls')
console.log(response.data)

// POST
const res = await axios.post('/api/urls', {
    original_url: 'https://google.com',
    custom_slug: 'ggl'
})

// Dengan auth header
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
```

---

## BAB 3 — Setup Project {#bab-3}

### 3.1 Struktur Folder

```
url-shortener/          <- root project
├── backend/            <- Go backend
│   ├── main.go
│   ├── go.mod
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   ├── user.go
│   │   ├── url.go
│   │   └── click.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── url.go
│   │   └── analytics.go
│   ├── middleware/
│   │   └── auth.go
│   └── utils/
│       └── shortcode.go
└── frontend/           <- React frontend
    ├── index.html
    ├── package.json
    ├── vite.config.js
    └── src/
        ├── main.jsx
        ├── App.jsx
        ├── lib/
        │   └── axios.js
        ├── context/
        │   └── AuthContext.jsx
        ├── pages/
        │   ├── LoginPage.jsx
        │   ├── RegisterPage.jsx
        │   └── DashboardPage.jsx
        └── components/
            ├── Navbar.jsx
            ├── URLCard.jsx
            └── AnalyticsChart.jsx
```

### 3.2 Init Backend

```bash
mkdir -p url-shortener/backend url-shortener/frontend
cd url-shortener/backend

go mod init url-shortener

# Install semua dependency backend
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
```

### 3.3 Init Frontend

```bash
cd ../
npm create vite@latest frontend -- --template react
cd frontend

npm install
npm install axios react-router-dom recharts
npm install -D tailwindcss @tailwindcss/vite
```

### 3.4 File yang tidak boleh masuk git

Buat `backend/.gitignore`:
```
app.db
.env
```

Buat `frontend/.gitignore` (tambahkan ke file yang sudah dibuat Vite):
```
node_modules
.env.local
```

`JWT_SECRET` tidak boleh masuk git — isi lewat environment variables
di Railway, bukan lewat file `.env` yang di-commit. File `app.db`
juga jangan di-commit: isinya data user dan URL lo.

---

## BAB 4 — Backend (Go + Gin) {#bab-4}

> Buat semua file ini sesuai path-nya. Jangan skip.

### 4.1 `backend/config/config.go`

```go
package config

import (
	"os"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	BaseURL   string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "app.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &Config{
		Port:      port,
		DBPath:    dbPath,
		JWTSecret: jwtSecret,
		BaseURL:   baseURL,
	}
}
```

### 4.2 `backend/models/user.go`

```go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // nambah ID, CreatedAt, UpdatedAt, DeletedAt otomatis
	Username   string `gorm:"uniqueIndex;not null" json:"username"`
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	Password   string `gorm:"not null" json:"-"` // json:"-" artinya tidak dikirim ke response
	URLs       []URL  `gorm:"foreignKey:UserID"` // relasi: satu user punya banyak URL
}
```

### 4.3 `backend/models/url.go`

```go
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
```

### 4.4 `backend/models/click.go`

```go
package models

import "gorm.io/gorm"

type Click struct {
	gorm.Model
	URLID     uint   `gorm:"not null" json:"url_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Region    string `json:"region"` // negara dari IP, best-effort (boleh kosong)
}
```

### 4.5 `backend/database/database.go`

```go
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
		log.Fatal("Gagal koneksi database:", err)
	}

	// Auto migrate: buat tabel kalau belum ada
	err = DB.AutoMigrate(&models.User{}, &models.URL{}, &models.Click{})
	if err != nil {
		log.Fatal("Gagal migrate database:", err)
	}

	log.Println("Database siap")
}
```

### 4.6 `backend/utils/shortcode.go`

```go
package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// rng dipakai bersama oleh semua pemanggil. Jangan buat rand.New di
// dalam fungsi tiap kali dipanggil: seed yang sama dalam waktu yang
// berdekatan bisa menghasilkan kode yang sama.
//
// Catatan: kode lama di sini memanggil `rand.New(...)` di func init()
// tapi hasilnya DIBUANG — itu no-op, tidak men-seed apa-apa. Sejak
// Go 1.20, rand global sudah di-seed acak otomatis, jadi `rand.Intn`
// langsung pun aman; rng eksplisit di sini supaya perilakunya jelas.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateShortCode buat kode random 6 karakter
func GenerateShortCode() string {
	length := 6
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rng.Intn(len(charset))]
	}
	return string(code)
}
```

### 4.7 `backend/middleware/auth.go`

```go
package middleware

import (
	"net/http"
	"strings"
	"url-shortener/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ada"})
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "format token salah"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		cfg := config.Load()

		// Parse dan validasi JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Tolak token yang tidak ditandatangani dengan HMAC
			// (mencegah serangan alg confusion).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid"})
			c.Abort()
			return
		}

		// Simpan user_id ke context. Type assertion SELALU pakai
		// bentuk `nilai, ok` — assertion langsung (tanpa ok) akan
		// panic dan mematikan server kalau claim tidak sesuai.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid"})
			c.Abort()
			return
		}
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid"})
			c.Abort()
			return
		}
		c.Set("user_id", uint(userIDFloat))

		c.Next()
	}
}
```

### 4.8 `backend/handlers/auth.go`

```go
package handlers

import (
	"net/http"
	"time"
	"url-shortener/config"
	"url-shortener/database"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func generateToken(userID uint, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // expired 7 hari
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal proses password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username atau email sudah dipakai"})
		return
	}

	cfg := config.Load()
	token, err := generateToken(user.ID, cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal buat token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
		return
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
		return
	}

	cfg := config.Load()
	token, err := generateToken(user.ID, cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal buat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
```

### 4.9 `backend/handlers/url.go`

```go
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

// recordClick: jalankan di goroutine (async) supaya redirect tidak lambat
func recordClick(urlID uint, ip, userAgent string) {
	click := models.Click{
		URLID:     urlID,
		IPAddress: ip,
		UserAgent: userAgent,
		Region:    lookupRegion(ip),
	}
	database.DB.Create(&click)
}

// lookupRegion menebak negara dari IP lewat ip-api.com (gratis, tanpa
// key, limit 45 request/menit). Gagal atau IP lokal/private -> string
// kosong, bukan error: tracking tidak boleh menggagalkan redirect.
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

	// Validasi custom slug: hanya huruf, angka, - dan _. Slug berisi
	// spasi atau "/" akan merusak URL redirect, jadi tolak di awal.
	if input.CustomSlug != "" && !validSlug(input.CustomSlug) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug hanya boleh huruf, angka, - dan _"})
		return
	}

	// Tentukan short code
	shortCode := input.CustomSlug
	if shortCode == "" {
		// Generate random sampai unik, maksimal 10 kali percobaan.
		// Loop tanpa batas berbahaya kalau ada masalah di DB.
		//
		// Catatan: cek-lalu-simpan ini punya race kecil — dua request
		// bersamaan bisa dapat kode yang sama, lalu Create di bawah
		// gagal karena unique index dan user dapat 500. Untuk skala
		// portfolio ini bisa diterima; di production tambahkan retry
		// saat Create gagal karena duplikat.
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
		// Cek apakah custom slug sudah dipakai
		var existing models.URL
		if database.DB.Where("short_code = ?", shortCode).First(&existing).Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "slug sudah dipakai"})
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

func ListURLs(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	cfg := config.Load()

	var urls []models.URL
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&urls)

	// Hitung click count tiap URL
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
		c.JSON(http.StatusNotFound, gin.H{"error": "URL tidak ditemukan"})
		return
	}

	// Hapus klik dulu, baru URL-nya. URL dihapus HARD delete
	// (Unscoped), bukan soft delete: gorm.Model punya DeletedAt, dan
	// baris yang soft-deleted MASIH menempati unique index short_code.
	// Kalau pakai soft delete, slug yang sudah dihapus tidak bisa
	// dipakai lagi (cek di CreateURL lolos, tapi insert gagal 500).
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

	// Catat klik secara async — tidak bikin redirect lambat
	go recordClick(url.ID, c.ClientIP(), c.GetHeader("User-Agent"))

	// 302 (Found), bukan 301. 301 di-cache permanen oleh browser: kalau
	// short URL dihapus lalu kodenya dipakai ulang, browser tetap
	// membuka tujuan lama dari cache. 302 tidak di-cache permanen.
	c.Redirect(http.StatusFound, url.OriginalURL)
}
```

### 4.10 `backend/handlers/analytics.go`

```go
package handlers

import (
	"net/http"
	"strconv"
	"time"
	"url-shortener/database"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
)

type DailyClick struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func GetAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// Ambil semua URL milik user
	var urls []models.URL
	database.DB.Where("user_id = ?", userID).Find(&urls)

	// Kumpulkan ID semua URL
	var urlIDs []uint
	for _, u := range urls {
		urlIDs = append(urlIDs, u.ID)
	}

	// Total klik semua URL
	var totalClicks int64
	if len(urlIDs) > 0 {
		database.DB.Model(&models.Click{}).Where("url_id IN ?", urlIDs).Count(&totalClicks)
	}

	// Klik per hari 7 hari terakhir. SQL hanya mengembalikan hari yang
	// ADA kliknya, jadi hari kosong harus diisi 0 supaya chart 7 hari
	// tampil lengkap, bukan bolong-bolong (lihat fillDailyClicks).
	dailyClicks := fillDailyClicks(urlIDs, 7)

	c.JSON(http.StatusOK, gin.H{
		"total_urls":   len(urls),
		"total_clicks": totalClicks,
		"daily_clicks": dailyClicks,
	})
}

func GetURLAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	urlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	// Pastikan URL milik user yang request
	var url models.URL
	if err := database.DB.Where("id = ? AND user_id = ?", urlID, userID).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL tidak ditemukan"})
		return
	}

	// Total klik URL ini
	var totalClicks int64
	database.DB.Model(&models.Click{}).Where("url_id = ?", url.ID).Count(&totalClicks)

	// Klik per hari (hari kosong diisi 0, lihat fillDailyClicks)
	dailyClicks := fillDailyClicks([]uint{url.ID}, 7)

	// Isi ClickCount supaya field click_count di response "url" tidak 0
	// padahal total_clicks di atas sudah benar.
	url.ClickCount = totalClicks

	c.JSON(http.StatusOK, gin.H{
		"url":          url,
		"total_clicks": totalClicks,
		"daily_clicks": dailyClicks,
	})
}

// fillDailyClicks mengembalikan klik per hari selama `days` hari
// terakhir (termasuk hari ini), dengan hari tanpa klik diisi 0.
func fillDailyClicks(urlIDs []uint, days int) []DailyClick {
	byDate := map[string]int{}
	if len(urlIDs) > 0 {
		var rows []DailyClick
		start := time.Now().AddDate(0, 0, -(days - 1)).Format("2006-01-02")
		database.DB.Raw(`
			SELECT DATE(created_at) as date, COUNT(*) as count
			FROM clicks
			WHERE url_id IN ? AND DATE(created_at) >= ?
			GROUP BY DATE(created_at)
			ORDER BY date ASC
		`, urlIDs, start).Scan(&rows)
		for _, r := range rows {
			byDate[r.Date] = r.Count
		}
	}

	result := make([]DailyClick, 0, days)
	today := time.Now()
	for i := days - 1; i >= 0; i-- {
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		result = append(result, DailyClick{Date: date, Count: byDate[date]})
	}
	return result
}
```

### 4.11 `backend/main.go`

```go
package main

import (
	"log"
	"net/http"
	"os"
	"url-shortener/config"
	"url-shortener/database"
	"url-shortener/handlers"
	"url-shortener/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	r := gin.Default()

	// CORS — izinkan frontend akses backend
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{frontendURL},
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// === PUBLIC ROUTES ===
	// Health check untuk Railway/monitoring. Route statis seperti ini
	// tetap menang atas /:code di bawah, jadi aman didaftarkan duluan.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Redirect handler
	r.GET("/:code", handlers.Redirect)

	// Auth
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	// === PROTECTED ROUTES (butuh JWT) ===
	api := r.Group("/api", middleware.Auth())
	{
		// URL management
		api.POST("/urls", handlers.CreateURL)
		api.GET("/urls", handlers.ListURLs)
		api.DELETE("/urls/:id", handlers.DeleteURL)

		// Analytics
		api.GET("/analytics", handlers.GetAnalytics)
		api.GET("/analytics/:id", handlers.GetURLAnalytics)
	}

	log.Printf("Server jalan di port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Gagal start server:", err)
	}
}
```

### Cek Backend Bisa Jalan

```bash
cd url-shortener/backend
go mod tidy   # download semua dependency yang belum ada
go run .      # jalankan server

# Output yang diharapkan:
# Database siap
# Server jalan di port 8080
# [GIN-debug] ...
```

Test dengan curl:
```bash
# Register user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"sofyan","email":"sofyan@test.com","password":"123456"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"sofyan@test.com","password":"123456"}'

# Buat short URL (ganti TOKEN dengan token dari hasil login)
curl -X POST http://localhost:8080/api/urls \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://google.com","custom_slug":"ggl"}'

# Coba redirect (lihat header Location, harus 302 ke Google)
curl -i http://localhost:8080/ggl

# Health check
curl http://localhost:8080/health
```

---

## BAB 5 — Frontend (React + Vite) {#bab-5}

### 5.1 Setup Tailwind

> Tailwind v4 (yang terinstall kalau `npm install` hari ini) tidak lagi
> pakai `tailwind.config.js` dan tidak punya perintah `npx tailwindcss
> init`. Konfigurasi pindah ke CSS. Tutorial lama yang menyuruh
> `init -p` itu alur v3 dan akan gagal di install baru — ini jebakan
> paling umum saat setup Tailwind di 2026.

Edit `frontend/src/index.css` (hapus semua isi, ganti dengan ini):
```css
@import "tailwindcss";
```

Satu baris itu menggantikan tiga directive lama (`@tailwind base;`
dst). Tidak perlu file `tailwind.config.js` — v4 mendeteksi file
sumber otomatis.

### 5.2 `frontend/vite.config.js`

```js
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: {
      // Forward /api dan /auth ke backend Go saat development
      '/api': 'http://localhost:8080',
      '/auth': 'http://localhost:8080',
    }
  }
})
```

### 5.3 `frontend/src/lib/axios.js`

```js
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '',
})

// Otomatis sisipkan token ke setiap request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Kalau 401, logout dan redirect ke login. Pengecualian: request
// /auth/login dan /auth/register sendiri. Tanpa pengecualian ini,
// login yang gagal (password salah) ikut me-redirect dan me-reload
// halaman, sehingga pesan error tidak pernah tampil.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const url = error.config?.url || ''
    const isAuthRequest = url.startsWith('/auth/')
    if (error.response?.status === 401 && !isAuthRequest) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
```

### 5.4 `frontend/src/context/AuthContext.jsx`

```jsx
import { createContext, useContext, useState } from 'react'

const AuthContext = createContext(null)

// Baca user dari localStorage dengan aman: kalau isinya rusak
// (bukan JSON valid), JSON.parse akan throw dan seluruh aplikasi
// crash saat load. Kembalikan null saja kalau gagal parse.
function loadStoredUser() {
  try {
    return JSON.parse(localStorage.getItem('user') || 'null')
  } catch {
    return null
  }
}

export function AuthProvider({ children }) {
  const [token, setToken] = useState(localStorage.getItem('token'))
  const [user, setUser] = useState(loadStoredUser)

  const login = (tokenValue, userData) => {
    localStorage.setItem('token', tokenValue)
    localStorage.setItem('user', JSON.stringify(userData))
    setToken(tokenValue)
    setUser(userData)
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ token, user, login, logout, isAuth: !!token }}>
      {children}
    </AuthContext.Provider>
  )
}

// Hook custom biar gampang dipakai
export function useAuth() {
  return useContext(AuthContext)
}
```

### 5.5 `frontend/src/App.jsx`

```jsx
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import DashboardPage from './pages/DashboardPage'

// Komponen penjaga route — kalau belum login, kirim ke /login
function ProtectedRoute({ children }) {
  const { isAuth } = useAuth()
  return isAuth ? children : <Navigate to="/login" replace />
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardPage />
              </ProtectedRoute>
            }
          />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
```

### 5.6 `frontend/src/main.jsx`

```jsx
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
```

### 5.7 `frontend/src/pages/LoginPage.jsx`

```jsx
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import api from '../lib/axios'

export default function LoginPage() {
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const res = await api.post('/auth/login', form)
      login(res.data.token, res.data.user)
      navigate('/dashboard')
    } catch (err) {
      setError(err.response?.data?.error || 'Login gagal')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white rounded-xl shadow-md p-8 w-full max-w-md">
        <h1 className="text-2xl font-bold text-gray-800 mb-6">Masuk</h1>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
          >
            {loading ? 'Memproses...' : 'Masuk'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Belum punya akun?{' '}
          <Link to="/register" className="text-blue-600 hover:underline">
            Daftar
          </Link>
        </p>
      </div>
    </div>
  )
}
```

### 5.8 `frontend/src/pages/RegisterPage.jsx`

```jsx
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import api from '../lib/axios'

export default function RegisterPage() {
  const [form, setForm] = useState({ username: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const res = await api.post('/auth/register', form)
      login(res.data.token, res.data.user)
      navigate('/dashboard')
    } catch (err) {
      setError(err.response?.data?.error || 'Registrasi gagal')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white rounded-xl shadow-md p-8 w-full max-w-md">
        <h1 className="text-2xl font-bold text-gray-800 mb-6">Daftar</h1>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Username
            </label>
            <input
              type="text"
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              minLength={3}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              minLength={6}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
          >
            {loading ? 'Memproses...' : 'Daftar'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Sudah punya akun?{' '}
          <Link to="/login" className="text-blue-600 hover:underline">
            Masuk
          </Link>
        </p>
      </div>
    </div>
  )
}
```

### 5.9 `frontend/src/components/Navbar.jsx`

```jsx
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { user, logout } = useAuth()

  return (
    <nav className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="max-w-5xl mx-auto flex justify-between items-center">
        <Link to="/dashboard" className="text-xl font-bold text-blue-600">
          ShortURL
        </Link>
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-600">Halo, {user?.username}</span>
          <button
            onClick={logout}
            className="text-sm text-red-600 hover:underline"
          >
            Keluar
          </button>
        </div>
      </div>
    </nav>
  )
}
```

### 5.10 `frontend/src/components/URLCard.jsx`

```jsx
import { useState } from 'react'

export default function URLCard({ url, onDelete }) {
  const [copied, setCopied] = useState(false)

  // navigator.clipboard mengembalikan Promise: tunggu hasilnya dulu.
  // alert() sebelumnya selalu muncul walau salin gagal, dan alert
  // memblokir seluruh halaman.
  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(url.short_url)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      alert('Gagal menyalin, salin manual dari address bar')
    }
  }

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4 flex justify-between items-start">
      <div className="flex-1 min-w-0">
        <a
          href={url.short_url}
          target="_blank"
          rel="noreferrer"
          className="text-blue-600 font-medium hover:underline text-sm"
        >
          {url.short_url}
        </a>
        <p className="text-gray-500 text-sm truncate mt-1">{url.original_url}</p>
        <div className="flex gap-4 mt-2 text-xs text-gray-400">
          <span>{url.click_count} klik</span>
          <span>{url.created_at}</span>
        </div>
      </div>
      <div className="flex gap-2 ml-4">
        <button
          onClick={copyToClipboard}
          className="text-xs bg-gray-100 hover:bg-gray-200 px-2 py-1 rounded"
        >
          {copied ? 'Disalin' : 'Salin'}
        </button>
        <button
          onClick={() => onDelete(url.id)}
          className="text-xs bg-red-50 hover:bg-red-100 text-red-600 px-2 py-1 rounded"
        >
          Hapus
        </button>
      </div>
    </div>
  )
}
```

### 5.11 `frontend/src/components/AnalyticsChart.jsx`

```jsx
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'

export default function AnalyticsChart({ data }) {
  if (!data || data.length === 0) {
    return (
      <div className="h-48 flex items-center justify-center text-gray-400 text-sm">
        Belum ada data klik
      </div>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={200}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
        <XAxis dataKey="date" tick={{ fontSize: 12 }} />
        <YAxis tick={{ fontSize: 12 }} allowDecimals={false} />
        <Tooltip />
        <Line
          type="monotone"
          dataKey="count"
          stroke="#3b82f6"
          strokeWidth={2}
          dot={{ r: 4 }}
          name="Klik"
        />
      </LineChart>
    </ResponsiveContainer>
  )
}
```

### 5.12 `frontend/src/pages/DashboardPage.jsx`

```jsx
import { useState, useEffect } from 'react'
import Navbar from '../components/Navbar'
import URLCard from '../components/URLCard'
import AnalyticsChart from '../components/AnalyticsChart'
import api from '../lib/axios'

export default function DashboardPage() {
  const [urls, setUrls] = useState([])
  const [analytics, setAnalytics] = useState(null)
  const [form, setForm] = useState({ original_url: '', custom_slug: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  // Load data saat halaman dibuka
  useEffect(() => {
    fetchURLs()
    fetchAnalytics()
  }, [])

  const fetchURLs = async () => {
    try {
      const res = await api.get('/api/urls')
      setUrls(res.data)
    } catch (err) {
      console.error('Gagal load URLs:', err)
    }
  }

  const fetchAnalytics = async () => {
    try {
      const res = await api.get('/api/analytics')
      setAnalytics(res.data)
    } catch (err) {
      console.error('Gagal load analytics:', err)
    }
  }

  const handleShorten = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await api.post('/api/urls', {
        original_url: form.original_url,
        custom_slug: form.custom_slug || undefined,
      })
      setForm({ original_url: '', custom_slug: '' })
      fetchURLs()
      fetchAnalytics()
    } catch (err) {
      setError(err.response?.data?.error || 'Gagal mempersingkat URL')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Hapus URL ini?')) return
    try {
      await api.delete(`/api/urls/${id}`)
      fetchURLs()
      fetchAnalytics()
    } catch (err) {
      alert('Gagal hapus URL')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-5xl mx-auto px-4 py-8 space-y-8">

        {/* Stats */}
        {analytics && (
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-white rounded-xl border border-gray-200 p-6">
              <p className="text-sm text-gray-500">Total URL</p>
              <p className="text-3xl font-bold text-gray-800">{analytics.total_urls}</p>
            </div>
            <div className="bg-white rounded-xl border border-gray-200 p-6">
              <p className="text-sm text-gray-500">Total Klik</p>
              <p className="text-3xl font-bold text-blue-600">{analytics.total_clicks}</p>
            </div>
          </div>
        )}

        {/* Chart */}
        {analytics && (
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <h2 className="text-base font-semibold text-gray-700 mb-4">
              Klik 7 Hari Terakhir
            </h2>
            <AnalyticsChart data={analytics.daily_clicks} />
          </div>
        )}

        {/* Form shorten */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-base font-semibold text-gray-700 mb-4">
            Persingkat URL
          </h2>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4 text-sm">
              {error}
            </div>
          )}

          <form onSubmit={handleShorten} className="flex gap-2 flex-wrap">
            <input
              type="url"
              placeholder="https://contoh.com/url-panjang"
              value={form.original_url}
              onChange={(e) => setForm({ ...form, original_url: e.target.value })}
              className="flex-1 min-w-0 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
            <input
              type="text"
              placeholder="custom-slug (opsional)"
              value={form.custom_slug}
              onChange={(e) => setForm({ ...form, custom_slug: e.target.value })}
              className="w-40 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="submit"
              disabled={loading}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50 font-medium"
            >
              {loading ? '...' : 'Persingkat'}
            </button>
          </form>
        </div>

        {/* Daftar URL */}
        <div>
          <h2 className="text-base font-semibold text-gray-700 mb-3">
            URL Kamu ({urls.length})
          </h2>
          {urls.length === 0 ? (
            <div className="text-center text-gray-400 py-12 text-sm">
              Belum ada URL yang dipersingkat
            </div>
          ) : (
            <div className="space-y-3">
              {urls.map((url) => (
                <URLCard key={url.id} url={url} onDelete={handleDelete} />
              ))}
            </div>
          )}
        </div>

      </main>
    </div>
  )
}
```

---

## BAB 6 — Integrasi & Testing {#bab-6}

### Jalankan Backend

```bash
cd url-shortener/backend
go run .
```

### Jalankan Frontend

```bash
# Terminal baru
cd url-shortener/frontend
npm run dev
```

Buka browser: `http://localhost:5173`

### Checklist Testing Manual

Coba semua ini satu per satu, pastikan tidak ada yang error:

- [ ] Register akun baru → seharusnya redirect ke dashboard
- [ ] Logout → seharusnya redirect ke login
- [ ] Login dengan akun yang baru dibuat
- [ ] Akses `/dashboard` tanpa login → harus redirect ke `/login`
- [ ] Shorten URL `https://google.com`
- [ ] Klik short URL yang muncul → harus redirect ke Google
- [ ] Cek dashboard → click count harus bertambah 1
- [ ] Shorten URL dengan custom slug
- [ ] Coba custom slug yang sama lagi → harus muncul error "slug sudah dipakai"
- [ ] Hapus salah satu URL

---

## BAB 7 — Deploy {#bab-7}

### 7.1 Deploy Backend ke Railway

```bash
# Install Railway CLI
npm install -g @railway/cli
railway login

# Di folder backend
cd url-shortener/backend
railway init
railway up
```

Set environment variables di Railway dashboard:
```
JWT_SECRET=isi-dengan-string-random-panjang
BASE_URL=https://url-yang-railway-berikan.railway.app
FRONTEND_URL=https://url-frontend-vercel.vercel.app
```

> Perhatian: filesystem Railway itu ephemeral — file `app.db` (SQLite)
> HILANG setiap kali service restart atau redeploy. Semua URL dan klik
> ikut hilang. Dua jalan keluar, pilih satu:
>
> 1. Railway Volume (tetap SQLite): di dashboard Railway, tambah Volume
>    lalu mount ke folder tempat `app.db` berada, dan set `DB_PATH` ke
>    path di dalam volume itu. Paling sedikit perubahan.
> 2. Railway Postgres (disarankan untuk deploy): provision Postgres di
>    Railway, ganti driver ke `gorm.io/driver/postgres`, dan baca koneksi
>    dari `DATABASE_URL`. SQLite tetap dipakai untuk development lokal.
>
> Jadi "SQLite sudah cukup" hanya benar untuk lokal — untuk deploy,
> tanpa volume atau Postgres, data demo lo hilang di restart pertama.

### 7.2 Deploy Frontend ke Vercel

```bash
# Install Vercel CLI
npm install -g vercel

# Di folder frontend
cd url-shortener/frontend

# Buat file .env.production
echo "VITE_API_URL=https://url-backend-railway.railway.app" > .env.production

vercel
# Ikuti prompt, pilih defaults
```

### 7.3 Update CORS

Setelah dapat URL Vercel, update `FRONTEND_URL` di Railway environment variables dengan URL Vercel yang sebenarnya, lalu redeploy.

---

## BAB 8 — Interview Prep {#bab-8}

### Cara Jelaskan Project (2 Menit)

> "Saya bikin URL Shortener dengan backend Go pakai framework Gin, dan frontend React dengan Vite. Backend handle registrasi, login JWT, shorten URL, dan tracking klik secara async pakai goroutine supaya redirect tetap cepat. Frontend-nya SPA dengan React Router, ada dashboard yang nampilkan statistik klik per hari pakai Recharts."

### Pertanyaan Teknis yang Mungkin Keluar

**Q: Kenapa pilih Go untuk backend?**
> Go punya performance bagus untuk I/O dan concurrency. Redirect handler perlu cepat, dan Go cocok karena goroutine-nya ringan dibanding thread.

**Q: Kenapa pakai goroutine di redirect handler?**
> Supaya proses record klik tidak bikin user nunggu redirect. Dengan `go recordClick(...)`, user langsung dapat redirect 301, dan click tracking jalan di background.

**Q: Gimana cara JWT bekerja di sini?**
> Saat login/register, server buat token JWT yang berisi user_id dan expiry, ditandatangani dengan secret key. Frontend simpan token di localStorage dan sisipkan di header setiap request. Middleware di Go ekstrak dan validasi token, lalu taruh user_id di context buat dipakai handler.

**Q: Bedanya SQLite sama PostgreSQL?**
> SQLite itu file-based, cocok untuk development dan skala kecil. PostgreSQL lebih cocok untuk production karena support concurrent writes yang lebih baik, scaling horizontal, dan banyak fitur enterprise.

**Q: Apa itu CORS dan kenapa perlu dikonfigurasi?**
> Cross-Origin Resource Sharing. Browser secara default blokir request dari domain berbeda. Backend Go dan frontend React jalan di domain/port berbeda, jadi perlu konfigurasi CORS di backend untuk izinkan request dari URL frontend.

**Q: Gimana React Router bekerja?**
> React Router handle navigasi di client-side tanpa reload halaman. ProtectedRoute bungkus halaman yang butuh auth — kalau token tidak ada di localStorage, redirect ke `/login`.

### Tips Demo

1. Buka dashboard kosong dulu — tunjukin tampilan awalnya
2. Register akun baru live di depan mereka
3. Shorten 2-3 URL, satu pakai custom slug
4. Klik short URL beberapa kali biar ada data
5. Tunjukin chart yang muncul di dashboard
6. Buka kode backend redirect handler dan jelaskan goroutine-nya

---

## Checklist Final Sebelum Interview

- [ ] Project bisa dijalankan lokal tanpa error
- [ ] Semua fitur di BAB 6 testing manual berjalan
- [ ] README.md sudah ditulis (cara install + cara run + penjelasan singkat)
- [ ] Project sudah di-push ke GitHub
- [ ] Live URL sudah bisa diakses (Railway + Vercel)
- [ ] Bisa jelaskan tiap file dan fungsinya
- [ ] Hafal jawaban untuk 6 pertanyaan di atas

---

---

## Catatan Tuning (10 Nov 2026)

Yang diubah dari versi awal, beserta alasannya. Stack tetap Go + React.

**Versi tooling (BAB 0, 3.3, 5.1, 5.2)**
- Go 1.21.13 naik ke 1.27.3 (1.21 sudah tidak dapat security fix; 1.27.4 rilis 10 Nov 2026). Node 20 naik ke 22 (20 sudah EOL). Export `GOPATH` dihapus, sudah default sejak Go 1.8.
- Tailwind v3 naik ke v4: `npx tailwindcss init -p` tidak ada di v4, sekarang pakai `@tailwindcss/vite` + `@import "tailwindcss";`, tanpa `tailwind.config.js`.

**Bug yang diperbaiki**
- Button di BAB 2.2: class dinamis `bg-${color}-500` tidak pernah jadi CSS. Diganti map warna statis.
- `shortcode.go`: `init()` lama memanggil `rand.New(...)` tapi hasilnya dibuang, jadi tidak men-seed apa-apa. Sekarang pakai `rng` package-level.
- Middleware JWT: tambah cek signing method HMAC dan type assertion bentuk `nilai, ok` supaya token aneh tidak bikin server panic.
- `CreateURL`: validasi slug (hanya huruf, angka, `-`, `_`), loop generate dibatasi 10 kali, plus catatan race condition cek-lalu-simpan.
- `DeleteURL`: ganti ke hard delete (`Unscoped`). Soft delete + unique index `short_code` membuat slug yang sudah dihapus tidak bisa dipakai lagi (cek lolos, insert gagal 500).
- `Redirect`: 301 ganti 302 supaya browser tidak meng-cache permanen tujuan lama.
- Interceptor axios 401: login gagal tidak lagi me-reload halaman, jadi pesan error "email atau password salah" akhirnya tampil.
- `AuthContext`: `JSON.parse` dibungkus try/catch supaya localStorage rusak tidak me-crash aplikasi saat load.
- `URLCard`: tombol salin sekarang menunggu hasil clipboard dan menampilkan "Disalin", bukan `alert()` yang selalu muncul walau gagal.

**Fitur dan kelengkapan**
- Analytics: hari tanpa klik diisi 0 supaya chart 7 hari lengkap (`fillDailyClicks`); `ClickCount` diisi di `GetURLAnalytics`; field `Clicks` tidak lagi ikut ke response JSON.
- Click tracking: tambah field `Region` (negara dari IP via ip-api.com, best-effort, gagal tidak masalah).
- Deploy: peringatan SQLite hilang di Railway karena filesystem ephemeral, plus opsi Railway Volume atau Postgres.
- Tambahan kecil: endpoint `/health`, `.gitignore` backend/frontend, `ListURLs` diurutkan terbaru dulu, dan contoh curl untuk shorten, redirect, dan health check.

*Dibuat untuk persiapan PKL November 2026 — Good luck.*
