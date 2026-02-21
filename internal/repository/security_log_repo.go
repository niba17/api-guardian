package repository

import (
	"api-guardian/internal/domain/security_log"
	"api-guardian/internal/domain/security_log/interfaces"

	"gorm.io/gorm"
)

type gormSecurityLogRepo struct {
	db *gorm.DB
}

// Return-nya sekarang menggunakan interfaces.SecurityLogRepository
func NewSecurityLogRepository(db *gorm.DB) interfaces.SecurityLogRepository {
	return &gormSecurityLogRepo{db: db}
}

func (r *gormSecurityLogRepo) Save(log *security_log.SecurityLog) error {
	return r.db.Create(log).Error
}

func (r *gormSecurityLogRepo) GetStats() (security_log.DashboardStats, error) {
	var stats security_log.DashboardStats
	var total, blocked, unique int64
	var avg float64

	// Query satu-satu untuk akurasi tinggi
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

func (r *gormSecurityLogRepo) GetRecent(limit int) ([]security_log.SecurityLog, error) {
	var logs []security_log.SecurityLog
	err := r.db.Order("timestamp desc").Limit(limit).Find(&logs).Error
	return logs, err
}
