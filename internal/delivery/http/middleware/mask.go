package middleware

import (
	"regexp"
)

var (
	// 1. Partial Email Masking: Menangkap value bos@gmail.com menjadi b****@gmail.com (Untuk teks biasa)
	emailRegex = regexp.MustCompile(`(?i)([a-z0-9._%+-])[a-z0-9._%+-]+@([a-z0-9.-]+\.[a-z]{2,})`)

	// 2. Universal Sensor: Menambahkan 'email' ke dalam daftar PII yang disensor TOTAL
	// Mendukung: "password":"123", email:bos@gmail.com, token=123
	piiRegex = regexp.MustCompile(`(?i)(password|passwd|pin|secret|token|api_key|cvv|email)[\s"']*[:=][\s"']*([^"',\s}]+)[\s"']*`)
)

func MaskPII(input string) string {
	if input == "" {
		return input
	}

	// Lapis 1: Sensor format email yang bertebaran bebas di dalam teks
	result := emailRegex.ReplaceAllString(input, "$1****@$2")

	// Lapis 2: Sensor total (*****) untuk key/field yang sifatnya rahasia, TERMASUK field 'email'
	// $1 = nama field-nya, $2 = value yang mau disensor
	result = piiRegex.ReplaceAllString(result, `"$1":"*****"`)

	return result
}
