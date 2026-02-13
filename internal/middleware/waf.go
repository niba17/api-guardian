package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url" // <-- TAMBAHAN BARU
	"strings"

	"github.com/rs/zerolog/log"
)

// Daftar Kata Terlarang (Blacklist)
var sqlInjectionPatterns = []string{
	"UNION SELECT", "OR 1=1", "--", "/*", "DROP TABLE", "INSERT INTO", "DELETE FROM",
	"xp_cmdshell", "exec(",
}

var xssPatterns = []string{
	"<script>", "javascript:", "onload=", "onerror=", "alert(", "document.cookie",
}

var pathTraversalPatterns = []string{
	"../", "..\\", "/etc/passwd", "C:\\Windows",
}

// User Agent yang mencurigakan (Tools Hacking)
var badUserAgents = []string{
	"sqlmap", "nikto", "nmap", "python", // Hapus "curl" biar Bos gampang testing
}

// BasicWAF adalah middleware untuk mendeteksi payload berbahaya
func BasicWAF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Cek User-Agent (Blokir Bot Nakal)
		ua := strings.ToLower(r.UserAgent())
		for _, badBot := range badUserAgents {
			if strings.Contains(ua, badBot) {
				blockRequest(w, r, "Bad Bot Detected: "+badBot)
				return
			}
		}

		// 2. Cek URL Query Parameters (GET) - IMPROVED! 🛡️
		// Kita decode dulu: "%27+OR+1=1" -> "' OR 1=1"
		decodedQuery, err := url.QueryUnescape(r.URL.RawQuery)
		if err == nil {
			// Kalau berhasil decode, pakai yang decoded. Kalau gagal, pakai raw.
			// Ubah ke lowercase biar case-insensitive
			lowerQuery := strings.ToLower(decodedQuery)
			if isMalicious(lowerQuery) {
				blockRequest(w, r, "Malicious Query Parameter Detected")
				return
			}
		} else {
			// Fallback cek raw query kalau decode gagal
			if isMalicious(strings.ToLower(r.URL.RawQuery)) {
				blockRequest(w, r, "Malicious Query Parameter Detected")
				return
			}
		}

		// 3. Cek Request Body (POST/PUT)
		// Kita baca body, cek isinya, lalu kembalikan lagi (Restore) supaya bisa dibaca Backend
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			// Baca body ke memory
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Kembalikan body ke tempat asalnya (PENTING! Kalau lupa, backend bakal terima body kosong)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Ubah ke string lowercase buat dicek
			bodyString := strings.ToLower(string(bodyBytes))
			if isMalicious(bodyString) {
				blockRequest(w, r, "Malicious Payload in Body Detected")
				return
			}
		}

		// Kalau aman, silakan lewat
		next.ServeHTTP(w, r)
	})
}

// Fungsi Helper untuk cek pattern
func isMalicious(input string) bool {
	// Cek SQLi
	for _, pattern := range sqlInjectionPatterns {
		if strings.Contains(input, strings.ToLower(pattern)) {
			return true
		}
	}
	// Cek XSS
	for _, pattern := range xssPatterns {
		if strings.Contains(input, strings.ToLower(pattern)) {
			return true
		}
	}
	// Cek Path Traversal
	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(input, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func blockRequest(w http.ResponseWriter, r *http.Request, reason string) {
	// Catat log merah
	log.Warn().
		Str("ip", r.RemoteAddr).
		Str("reason", reason).
		Str("path", r.URL.Path).
		Msg("⛔ WAF BLOCKED REQUEST")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "Security Violation",
		"message": "Request blocked by API Guardian WAF. " + reason,
	})
}
