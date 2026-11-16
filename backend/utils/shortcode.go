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
