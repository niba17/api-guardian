package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	foundAny := false

	fmt.Println("🔍 Memulai Pencarian di Semua Database (0-15)...")

	for db := 0; db <= 15; db++ {
		rdb := redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   db,
		})

		keys, _ := rdb.Keys(ctx, "*").Result()
		if len(keys) > 0 {
			foundAny = true
			fmt.Printf("\n📦 [DATABASE %d] Ditemukan %d Keys:\n", db, len(keys))
			for _, key := range keys {
				val, _ := rdb.Get(ctx, key).Result()
				ttl, _ := rdb.TTL(ctx, key).Result()
				fmt.Printf("   🔑 %-20s | 📄 %-10s | ⏳ %v\n", key, val, ttl)
			}
		}
		rdb.Close()
	}

	if !foundAny {
		fmt.Println("\n❌ NIHIL. Tidak ada data di semua database pada 127.0.0.1:6379")
		fmt.Println("💡 Tip: Pastikan server Go Bos masih menyala saat menjalankan script ini.")
	}
}
