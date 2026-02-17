package maskutil

import "regexp"

// Regex Final: Kita tambahkan "? di akhir untuk menangkap kutip penutup jika ada
var sensitiveFields = regexp.MustCompile(`(?i)"?(password|token|api_key|secret|cvv|card_number|credit_card)"?[\s"':]+([^,"'}\s]+)"?`)

func MaskPII(input string) string {
	if input == "" {
		return ""
	}
	// Menggunakan format standar JSON: "key":"[MASKED]"
	return sensitiveFields.ReplaceAllString(input, `"$1":"[MASKED]"`)
}
