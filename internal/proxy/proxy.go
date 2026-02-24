package proxy

import (
	"api-guardian/internal/config" // 👈 Import config
	"api-guardian/internal/delivery/http/middleware"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

type LoadBalancer struct {
	proxies []*httputil.ReverseProxy
	counter uint64
}

// 🚀 Ubah parameter dari targets []string menjadi cfg *config.AppConfig
func NewLoadBalance(cfg *config.AppConfig) (*LoadBalancer, error) {
	targets := cfg.TargetURLs

	// 🛡️ PROTEKSI 1: Cek apakah target kosong di awal
	if len(targets) == 0 {
		return nil, fmt.Errorf("load balancer error: no valid backend targets provided")
	}

	var proxies []*httputil.ReverseProxy

	for i, target := range targets {
		if target == "" {
			continue
		}

		// 🛡️ PROTEKSI 2: Tangkap error parsing URL
		targetURL, err := url.Parse(target)
		if err != nil {
			log.Warn().Err(err).Str("target", target).Msg("Skipping invalid target URL")
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 🚀 INI DIA PERBAIKANNYA: Panggil dengan 6 Parameter dari Config!
		cb := middleware.NewCircuitBreak(
			fmt.Sprintf("Backend-%d", i),
			cfg.CBMaxRequests,
			cfg.CBIntervalSec,
			cfg.CBTimeoutSec,
			cfg.CBMinRequests,
			cfg.CBFailRatio,
		)

		proxy.Transport = &middleware.CircuitBreakTransport{
			Transport: http.DefaultTransport,
			CB:        cb,
		}

		setupProxyCallbacks(proxy, targetURL, target)

		proxies = append(proxies, proxy)
	}

	// 🛡️ PROTEKSI 3: Pastikan ada minimal 1 proxy yang sukses dibuat
	if len(proxies) == 0 {
		return nil, fmt.Errorf("load balancer error: all provided targets were invalid")
	}

	return &LoadBalancer{proxies: proxies}, nil
}

func setupProxyCallbacks(proxy *httputil.ReverseProxy, targetURL *url.URL, rawTarget string) {
	// 1. MODIFIKASI REQUEST (Berangkat ke Backend)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Host = targetURL.Host
		req.Header.Set("X-Origin-Host", targetURL.Host)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

		// 🚀 Meneruskan IP asli ke backend
		clientIP := middleware.GetIP(req)
		req.Header.Set("X-Real-IP", clientIP)

		existingXFF := req.Header.Get("X-Forwarded-For")
		if existingXFF != "" {
			req.Header.Set("X-Forwarded-For", existingXFF+", "+clientIP)
		} else {
			req.Header.Set("X-Forwarded-For", clientIP)
		}
	}

	// 2. MODIFIKASI RESPONSE (Pulang dari Backend)
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Server")
		resp.Header.Del("X-Powered-By")

		resp.Header.Set("Server", "api-guardian")
		resp.Header.Set("X-Content-Type-Options", "nosniff")
		resp.Header.Set("X-Frame-Options", "SAMEORIGIN")
		resp.Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return nil
	}

	// 3. PENANGANAN ERROR (Backend Mati/Timeout)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Str("backend", rawTarget).Msg("Proxy Error")

		w.Header().Set("Content-Type", "application/json")

		if err != nil && err.Error() == "circuit breaker is OPEN" {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "Service Unavailable", "message": "Circuit Breaker is OPEN"}`))
			return
		}

		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error": "Bad Gateway", "message": "Backend server is unreachable"}`))
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddUint64(&lb.counter, 1)
	index := count % uint64(len(lb.proxies))
	lb.proxies[index].ServeHTTP(w, r)
}
