package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// 1. SIMPLE STRINGS (Cek cepat pakai strings.Contains)
var simpleSqlPatterns = []string{
	"UNION SELECT", "DROP TABLE", "INSERT INTO", "DELETE FROM",
	"xp_cmdshell", "exec(", "--", "/*", "*/",
}

// 2. REGEX PATTERNS (Cek pola rumit yang bervariasi spasi/kapital)
// Gunakan `MustCompile` biar kalau regex salah, server langsung panic di awal (fail fast)
var regexSqlPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(union\s+select)`),     // UNION SELECT (case insensitive + spasi bebas)
	regexp.MustCompile(`(?i)(select\s+.*\s+from)`), // SELECT ... FROM
	regexp.MustCompile(`(?i)(insert\s+into)`),      // INSERT INTO
	regexp.MustCompile(`(?i)(update\s+.*\s+set)`),  // UPDATE SET

	// 👇 INI JUARA KITA: Menangkap Variasi OR 1=1
	regexp.MustCompile(`(?i)('\s*or\s*)`),         // ' OR (dengan spasi fleksibel)
	regexp.MustCompile(`(?i)(or\s+1\s*=\s*1)`),    // OR 1=1
	regexp.MustCompile(`(?i)(1\s*=\s*1)`),         // 1=1
	regexp.MustCompile(`(?i)('\s*=\s*')`),         // '='
	regexp.MustCompile(`(?i)(['";])`),             // Karakter pemutus query
	regexp.MustCompile(`(?i)(benchmark|sleep)\(`), // Time-based SQL Injection
}

var xssPatterns = []string{
	"<script>", "javascript:", "onload=", "onerror=", "alert(", "document.cookie",
}

var pathTraversalPatterns = []string{
	"../", "..\\", "/etc/passwd", "C:\\Windows",
}

var badUserAgents = []string{
	"sqlmap", "nikto", "nmap", "python",
}

// BasicWAF Middleware
func BasicWAF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Helper internal untuk mempermudah blocking dengan context
		triggerBlock := func(reason string) {
			// 👇 KUNCI UTAMA: Masukkan alasan ke Context agar dibaca logger.go
			ctx := context.WithValue(r.Context(), "waf_reason", reason)
			blockRequest(w, r.WithContext(ctx), reason)
		}

		// 1. Cek Path
		decodedPath, _ := url.PathUnescape(r.URL.Path)
		if isMalicious(strings.ToLower(decodedPath)) {
			triggerBlock("Path Traversal Detected")
			return
		}

		// 2. Cek User-Agent
		ua := strings.ToLower(r.UserAgent())
		for _, badBot := range badUserAgents {
			if strings.Contains(ua, badBot) {
				triggerBlock("Security Scanner: " + badBot)
				return
			}
		}

		// 3. Cek URL Query (Kita buat lebih spesifik pesannya)
		decodedQuery, err := url.QueryUnescape(r.URL.RawQuery)
		targetQuery := r.URL.RawQuery
		if err == nil {
			targetQuery = decodedQuery
		}

		if isMalicious(targetQuery) {
			reason := "Malicious Query Detected"
			// Deteksi spesifik untuk ThreatDetails
			upperQuery := strings.ToUpper(targetQuery)
			if strings.Contains(upperQuery, "SELECT") || strings.Contains(upperQuery, "OR 1=1") {
				reason = "SQL Injection Attempt"
			} else if strings.Contains(strings.ToLower(targetQuery), "<script>") {
				reason = "XSS Attack Attempt"
			}

			triggerBlock(reason)
			return
		}

		// 4. Cek Body
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				if isMalicious(string(bodyBytes)) {
					triggerBlock("Malicious Payload in Body")
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Fungsi Helper Utama
func isMalicious(input string) bool {
	// Normalisasi input (lowercase) untuk string matching biasa
	lowerInput := strings.ToLower(input)

	// A. Cek SQL Injection (Simple String)
	for _, pattern := range simpleSqlPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}

	// B. Cek SQL Injection (Regex - Case Insensitive sudah dihandle regex (?i))
	for _, regex := range regexSqlPatterns {
		if regex.MatchString(input) { // Pakai input asli (case sensitive matters for regex sometimes)
			return true
		}
	}

	// C. Cek XSS
	for _, pattern := range xssPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}

	// D. Cek Path Traversal
	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

func blockRequest(w http.ResponseWriter, r *http.Request, reason string) {
	// 👇 TITIP PESAN DI HEADER (Agar Logger bisa baca meskipun Context-nya lepas)
	w.Header().Set("X-Guardian-WAF-Reason", reason)

	log.Warn().Str("reason", reason).Msg("WAF BLOCKED REQUEST")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "Security Violation",
		"message": "Request blocked by API Guardian WAF. " + reason,
	})
}
