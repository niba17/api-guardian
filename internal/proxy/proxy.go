package proxy

import (
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

func NewLoadBalancer(targets []string) (*LoadBalancer, error) {
	var proxies []*httputil.ReverseProxy

	for i, target := range targets {
		if target == "" {
			continue
		}
		targetURL, _ := url.Parse(target)

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// ⚡ GUNAKAN CIRCUIT BREAKER DARI MIDDLEWARE
		cb := middleware.NewCircuitBreaker(fmt.Sprintf("Backend-%d", i))
		proxy.Transport = &middleware.CircuitBreakerTransport{
			Transport: http.DefaultTransport,
			CB:        cb,
		}

		// Director & ErrorHandler tetap sama...
		setupProxyCallbacks(proxy, targetURL, target)

		proxies = append(proxies, proxy)
	}

	return &LoadBalancer{proxies: proxies}, nil
}

// Helper untuk merapikan NewLoadBalancer
func setupProxyCallbacks(proxy *httputil.ReverseProxy, targetURL *url.URL, rawTarget string) {
	// 1. Simpan director asli bawaan Go
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		// 2. WAJIB: Jalankan director asli agar skema HTTPS tidak hilang
		originalDirector(req)

		// 3. Tambahkan modifikasi kita
		req.Host = targetURL.Host
		req.Header.Set("X-Origin-Host", targetURL.Host)
		// Supaya Google/Backend tahu ini request dari Proxy
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	}

	// ErrorHandler tetap aman...
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Str("backend", rawTarget).Msg("Proxy Error")
		if err.Error() == "circuit breaker is OPEN" {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "Circuit Breaker Open"}`))
			return
		}
		w.WriteHeader(http.StatusBadGateway)
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddUint64(&lb.counter, 1)
	index := count % uint64(len(lb.proxies))
	lb.proxies[index].ServeHTTP(w, r)
}
