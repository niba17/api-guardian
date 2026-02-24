package dto

// HealthResponse adalah format balasan resmi untuk status server
type HealthResponse struct {
	System          string `json:"system"`
	RedisConnection string `json:"redis_connection"`
	CircuitBreaker  string `json:"circuit_breaker"`
	Time            string `json:"time"`
}
