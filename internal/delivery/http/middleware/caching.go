package middleware

import (
	"api-guardian/internal/usecase"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// SmartCache adalah middleware caching agresif untuk GET request
func SmartCache(uc *usecase.CacheUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Caching cuma berlaku untuk Method GET
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// 2. Bikin Key Unik: "cache:/path?query"
		cacheKey := fmt.Sprintf("cache:%s?%s", r.URL.Path, r.URL.RawQuery)

		// 3. Cek Redis Dulu (HIT)
		if val, err := uc.Get(r.Context(), cacheKey); err == nil {
			// Kalau ada, langsung kirim balik! (Backend gak kerja)
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Content-Type", "application/json") // Asumsi JSON
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return
		}

		// 4. Kalau Tidak Ada (MISS) -> Teruskan ke Backend tapi "Rekam" jawabannya
		// Kita pakai Recorder untuk menangkap output dari backend
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		// 5. Cek apakah backend sukses (200 OK)
		// Kita cuma mau cache kalau sukses. Kalau error jangan dicache.
		if recorder.Code == http.StatusOK {
			// Simpan ke Redis (Async biar gak nambah latency user)
			go func() {
				_ = uc.Set(r.Context(), cacheKey, recorder.Body.Bytes())
			}()
		}

		// 6. Salin jawaban backend ke user
		w.Header().Set("X-Cache", "MISS")
		// Copy semua header dari backend
		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}
