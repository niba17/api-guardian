package domain

import "time"

// Tambahkan tag gorm di sini
type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"unique;not null;size:50" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"` // Explicit nama kolom
	Role         string    `gorm:"default:'Viewer';size:20" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
