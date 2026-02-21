package usecase

import (
	"context"
	"time"

	// 👈 Import kontrak CacheRepository yang sudah kita buat
	// Sesuaikan path-nya jika Bos menaruhnya di folder lain
	cacheIntf "api-guardian/internal/domain/cache/interfaces"
)

// CacheUsecase menangani logika bisnis caching
type CacheUsecase struct {
	cacheRepo cacheIntf.CacheRepository // 👈 Ganti *redis.Client dengan Interface
	TTL       time.Duration
}

// Constructor sekarang meminta cacheRepo, bukan rdb mentah
func NewCacheUsecase(repo cacheIntf.CacheRepository, ttl time.Duration) *CacheUsecase {
	return &CacheUsecase{
		cacheRepo: repo,
		TTL:       ttl,
	}
}

// Get mengambil data cache berdasarkan key
func (u *CacheUsecase) Get(ctx context.Context, key string) (string, error) {
	// 🚀 Langsung oper ke Repository
	return u.cacheRepo.Get(ctx, key)
}

// Set menyimpan data ke cache dengan durasi TTL
func (u *CacheUsecase) Set(ctx context.Context, key string, value []byte) error {
	// 🚀 Langsung oper ke Repository
	return u.cacheRepo.Set(ctx, key, value, u.TTL)
}
