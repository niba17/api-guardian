package main

import (
	"api-guardian/internal/storage"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Pastikan alamat ini SAMA PERSIS dengan di .env atau default config
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	adapter := storage.NewRedisAdapter(rdb)
	ctx := context.Background()

	// KITA TANAM JEJAK PERMANEN
	fmt.Println("📍 Menanam Jejak: 'JEJAK_SI_GUNDUL'")
	// Pass 0 sebagai durasi agar permanen (atau gunakan func Set tanpa expire jika adapter mendukung, tapi 0 usually means keep it or persist depending on impl, wait Set usually takes duration. Let's use a very long duration to be safe if adapter logic differs, but standard Redis Set with 0 expiration is usually not valid for 'keep forever' in some wrappers, better use -1 or just a long time for test)
	// Revisi: Adapter Bos pakai r.Client.Set(..., expiration). Di go-redis, 0 means no expiration (permanent).
	err := adapter.Set(ctx, "JEJAK_SI_GUNDUL", "SAYA_ADA_DI_SINI", 0)
	if err != nil {
		fmt.Printf("❌ Gagal Tanam: %v\n", err)
		return
	}

	fmt.Println("✅ Jejak tertanam! Sekarang cari saya di redis-cli.")
}
