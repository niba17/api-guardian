package interfaces

import "context"

type BanUsecase interface {
	// Kita sepakati namanya ExecuteAutoBan agar konsisten dengan panggilan di Middleware
	ExecuteAutoBan(ctx context.Context, ip string) error
}
