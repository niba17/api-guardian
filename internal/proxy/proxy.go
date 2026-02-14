package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
)

// LoadBalancer memegang banyak backend proxy & counter
type LoadBalancer struct {
	proxies []*httputil.ReverseProxy
	counter uint64
}

// BreakerTransport "Kurir" yang membawa sekering SPESIFIK untuk 1 backend
type BreakerTransport struct {
	RoundTripper http.RoundTripper
	cb           *gobreaker.CircuitBreaker
}

// RoundTrip menjalankan request melewati sekering
func (c *BreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Execute mengharapkan return interface{}, error
	result, err := c.cb.Execute(func() (interface{}, error) {
		resp, err := c.RoundTripper.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		// Anggap status 500+ sebagai kegagalan sistem backend
		if resp.StatusCode >= 500 {
			// Kita return error supaya dihitung sebagai failure oleh circuit breaker
			return resp, fmt.Errorf("backend error: %d", resp.StatusCode)
		}
		return resp, nil
	})

	if err != nil {
		// Jika errornya dari backend (status 500+), kita tetap kembalikan response aslinya
		// tapi circuit breaker sudah mencatatnya sebagai kegagalan.
		if resp, ok := result.(*http.Response); ok && resp != nil {
			return resp, nil
		}

		if err == gobreaker.ErrOpenState {
			return nil, fmt.Errorf("circuit breaker is OPEN")
		}
		return nil, err
	}

	return result.(*http.Response), nil
}

// Helper untuk membuat Sekering baru
func createCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 1, // Berapa request boleh lolos saat status Half-Open
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second, // Berapa lama sekering mati sebelum coba lagi
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failRatio >= 0.6
		},
	})
}

// NewLoadBalancer inisialisasi semua backend
func NewLoadBalancer(targets []string) (*LoadBalancer, error) {
	var proxies []*httputil.ReverseProxy

	for i, target := range targets {
		if target == "" {
			continue
		}
		targetURL, err := url.Parse(target)
		if err != nil {
			log.Error().Err(err).Str("url", target).Msg("Invalid Backend URL")
			continue
		}

		// 1. Buat Proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 2. Buat Sekering Khusus
		cbName := fmt.Sprintf("CB-Backend-%d", i)
		cb := createCircuitBreaker(cbName)

		// 3. Pasang Transport Custom (Sekering)
		// Penting: Gunakan pointer receiver
		proxy.Transport = &BreakerTransport{
			RoundTripper: http.DefaultTransport,
			cb:           cb,
		}

		// 4. Director (Header Manipulation)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Set Host header agar backend mengenali request (PENTING untuk Google/Cloudflare)
			req.Host = targetURL.Host
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Origin-Host", targetURL.Host)
			if req.Header.Get("User-Agent") == "" {
				req.Header.Set("User-Agent", "API-Guardian/1.0")
			}
		}

		// 5. Error Handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Error().Err(err).Str("backend", target).Msg("Backend Error")
			w.Header().Set("Content-Type", "application/json")

			if err.Error() == "circuit breaker is OPEN" {
				w.WriteHeader(http.StatusServiceUnavailable) // 503
				w.Write([]byte(`{"error": "Backend is resting (Circuit Open)"}`))
			} else {
				w.WriteHeader(http.StatusBadGateway) // 502
				w.Write([]byte(`{"error": "Can't connect to Backend"}`))
			}
		}

		proxies = append(proxies, proxy)
		log.Info().Str("target", target).Msg("Backend addressed on Balancer")
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("No valid backend targets found in configuration")
	}

	return &LoadBalancer{proxies: proxies}, nil
}

// ServeHTTP Implementasi Round-Robin
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddUint64(&lb.counter, 1)
	index := count % uint64(len(lb.proxies))
	lb.proxies[index].ServeHTTP(w, r)
}
