package repository

import (
	"api-guardian/internal/domain/dashboard"
	"api-guardian/internal/domain/security_log"
	"api-guardian/internal/domain/security_log/interfaces"

	"gorm.io/gorm"
)

// 🚀 1. GANTI: Jadi huruf kecil (Private)
type securityLogRepository struct {
	db *gorm.DB
}

// 🚀 2. CONSTRUCTOR: Return-nya menggunakan Interface dari Domain
func NewSecurityLogRepository(db *gorm.DB) interfaces.SecurityLogRepository {
	return &securityLogRepository{db: db}
}

// 🚀 3. METHOD: Tetap Besar (Harus sesuai dengan Interface di Domain)
func (r *securityLogRepository) Save(log *security_log.SecurityLog) error {
	return r.db.Create(log).Error
}

func (r *securityLogRepository) GetStats() (dashboard.Dashboard, error) {
	var stats dashboard.Dashboard
	var total, blocked, unique int64
	var avg float64

	r.db.Model(&security_log.SecurityLog{}).Count(&total)
	r.db.Model(&security_log.SecurityLog{}).Where("is_blocked = ?", true).Count(&blocked)
	r.db.Model(&security_log.SecurityLog{}).Distinct("ip").Count(&unique)
	r.db.Model(&security_log.SecurityLog{}).Select("COALESCE(AVG(latency), 0)").Scan(&avg)

	stats.TotalRequests = int(total)
	stats.TotalBlocked = int(blocked)
	stats.TotalSuccess = int(total - blocked)
	stats.UniqueIPs = int(unique)
	stats.AvgLatency = int64(avg)

	return stats, nil
}

func (r *securityLogRepository) GetRecent(limit int) ([]security_log.SecurityLog, error) {
	var logs []security_log.SecurityLog
	err := r.db.Order("timestamp desc").Limit(limit).Find(&logs).Error
	return logs, err
}
