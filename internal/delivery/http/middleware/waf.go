package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// 1. SIMPLE STRINGS
var simpleSqlPatterns = []string{
	"UNION SELECT", "DROP TABLE", "INSERT INTO", "DELETE FROM",
	"xp_cmdshell", "exec(", "--", "/*", "*/",
}

// 2. REGEX PATTERNS
var regexSqlPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(union\s+select)`),
	regexp.MustCompile(`(?i)(select\s+.*\s+from)`),
	regexp.MustCompile(`(?i)(insert\s+into)`),
	regexp.MustCompile(`(?i)(update\s+.*\s+set)`),
	regexp.MustCompile(`(?i)('\s*or\s*)`),
	regexp.MustCompile(`(?i)(or\s+1\s*=\s*1)`),
	regexp.MustCompile(`(?i)(1\s*=\s*1)`),
	regexp.MustCompile(`(?i)('\s*=\s*')`),
	regexp.MustCompile(`(?i)(['";])`),
	regexp.MustCompile(`(?i)(benchmark|sleep)\(`),
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

// 🚀 FUNGSI BANTUAN UNTUK CEK WHITELIST
func isWhitelisted(ip string, whitelist []string) bool {
	for _, allowed := range whitelist {
		if ip == allowed {
			return true
		}
	}
	return false
}

// 🚀 BasicWAF Middleware (Sekarang menerima daftar Whitelist)
func BasicWAF(whitelistIPs []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 👑 1. LOGIKA VVIP BYPASS (Karpet Merah)
			clientIP := GetIP(r)
			if isWhitelisted(clientIP, whitelistIPs) {
				// Kalau IP terdaftar, catat log ringan dan biarkan lewat tanpa periksa isinya!
				log.Debug().Str("ip", clientIP).Msg("👑 VVIP IP Bypassed WAF Inspection")
				next.ServeHTTP(w, r)
				return
			}

			// Pengecualian Path (Misal untuk login)
			if r.URL.Path == "/api/login" {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Cek Path
			decodedPath, _ := url.PathUnescape(r.URL.Path)
			if isMalicious(strings.ToLower(decodedPath)) {
				blockRequest(w, r, "Path Traversal Detected")
				return
			}

			// 3. Cek User-Agent
			ua := strings.ToLower(r.UserAgent())
			for _, badBot := range badUserAgents {
				if strings.Contains(ua, badBot) {
					blockRequest(w, r, "Security Scanner: "+badBot)
					return
				}
			}

			// 4. Cek URL Query
			decodedQuery, err := url.QueryUnescape(r.URL.RawQuery)
			targetQuery := r.URL.RawQuery
			if err == nil {
				targetQuery = decodedQuery
			}

			if isMalicious(targetQuery) {
				reason := "Malicious Query Detected"
				upperQuery := strings.ToUpper(targetQuery)
				if strings.Contains(upperQuery, "SELECT") || strings.Contains(upperQuery, "OR 1=1") {
					reason = "SQL Injection Attempt"
				} else if strings.Contains(strings.ToLower(targetQuery), "<script>") {
					reason = "XSS Attack Attempt"
				}

				blockRequest(w, r, reason)
				return
			}

			// 5. Cek Body (Dengan Limit Memori! 🛡️)
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 1048576))

				if err == nil {
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					if isMalicious(string(bodyBytes)) {
						blockRequest(w, r, "Malicious Payload in Body")
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isMalicious(input string) bool {
	lowerInput := strings.ToLower(input)

	for _, pattern := range simpleSqlPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}

	for _, regex := range regexSqlPatterns {
		if regex.MatchString(input) {
			return true
		}
	}

	for _, pattern := range xssPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}

	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

func blockRequest(w http.ResponseWriter, r *http.Request, reason string) {
	w.Header().Set("X-Guardian-WAF-Reason", reason)
	log.Warn().Str("ip", GetIP(r)).Str("reason", reason).Msg("WAF BLOCKED REQUEST")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "Security Violation",
		"message": "Request blocked by API Guardian WAF. " + reason,
	})
}
