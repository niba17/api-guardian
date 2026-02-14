package middleware

import (
	"regexp"
)

var (
	// Email: Menangkap b****@gmail.com
	emailRegex = regexp.MustCompile(`(?i)([a-z0-9._%+-])[a-z0-9._%+-]+@([a-z0-9.-]+\.[a-z]{2,})`)

	// Universal Sensor: Menangkap password, pin, secret, dll dalam format JSON atau Teks Biasa
	// Mendukung: "password":"123", password:123, password=123
	piiRegex = regexp.MustCompile(`(?i)(password|passwd|pin|secret|token|api_key|cvv)[\s"']*[:=][\s"']*([^"',\s}]+)[\s"']*`)
)

func MaskPII(input string) string {
	if input == "" {
		return input
	}

	// 1. Sensor Email
	result := emailRegex.ReplaceAllString(input, "$1****@$2")

	// 2. Sensor Kredensial (Password, PIN, dll)
	// $1 = nama field-nya, $2 = value yang mau disensor
	result = piiRegex.ReplaceAllString(result, `"$1":"*****"`)

	return result
}
