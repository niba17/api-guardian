package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreakerTransport adalah decorator untuk http.RoundTripper
type CircuitBreakerTransport struct {
	Transport http.RoundTripper
	CB        *gobreaker.CircuitBreaker
}

// RoundTrip menjalankan request dengan proteksi sekering
func (c *CircuitBreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	result, err := c.CB.Execute(func() (interface{}, error) {
		resp, err := c.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		// Anggap status 500+ sebagai kegagalan sistem backend
		if resp.StatusCode >= 500 {
			return resp, fmt.Errorf("backend failure: %d", resp.StatusCode)
		}
		return resp, nil
	})

	if err != nil {
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

// NewCircuitBreaker factory untuk membuat sekering baru
func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 1,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failRatio >= 0.6
		},
	})
}
