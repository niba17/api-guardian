package middleware

import (
	"regexp"
)

var (
	// Email: Tetap aman
	emailRegex = regexp.MustCompile(`(?i)([a-z0-9._%+-])[a-z0-9._%+-]+@([a-z0-9.-]+\.[a-z]{2,})`)

	// Password: Lebih fleksibel (bisa handle "password":"xxx", password:xxx, atau password=xxx)
	// Kita cari kata password, lalu ambil karakter setelah titik dua/sama dengan sampai ketemu koma/spasi/tutup kurung
	passRegex = regexp.MustCompile(`(?i)(password\s*[:=]\s*["']?)([^"'\s,}]+)(["']?)`)
)

func MaskPII(input string) string {
	if input == "" {
		return input
	}

	// 1. Sensor Email
	result := emailRegex.ReplaceAllString(input, "$1****@$2")

	// 2. Sensor Password (Ditingkatkan!)
	// $1 = password: , $2 = value-nya, $3 = penutup kutip (jika ada)
	result = passRegex.ReplaceAllString(result, `$1******$3`)

	return result
}
