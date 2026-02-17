package repository

import (
	"api-guardian/internal/domain"

	"gorm.io/gorm"
)

type SecurityLogRepository interface {
	Save(log *domain.SecurityLog) error
	GetStats() (map[string]interface{}, error)
	GetRecent(limit int) ([]domain.SecurityLog, error)
}

type gormSecurityLogRepo struct {
	db *gorm.DB
}

func NewSecurityLogRepository(db *gorm.DB) SecurityLogRepository {
	return &gormSecurityLogRepo{db: db}
}

func (r *gormSecurityLogRepo) Save(log *domain.SecurityLog) error {
	return r.db.Create(log).Error
}

func (r *gormSecurityLogRepo) GetStats() (map[string]interface{}, error) {
	var totalReq, blockedReq, uniqueIPs int64
	var avgLatency float64

	r.db.Model(&domain.SecurityLog{}).Count(&totalReq)
	r.db.Model(&domain.SecurityLog{}).Where("is_blocked = ?", true).Count(&blockedReq)
	r.db.Model(&domain.SecurityLog{}).Distinct("ip").Count(&uniqueIPs)
	r.db.Model(&domain.SecurityLog{}).Select("COALESCE(AVG(latency), 0)").Scan(&avgLatency)

	return map[string]interface{}{
		"total_requests":   totalReq,
		"blocked_requests": blockedReq,
		"unique_ips":       uniqueIPs,
		"avg_latency":      int64(avgLatency),
	}, nil
}

func (r *gormSecurityLogRepo) GetRecent(limit int) ([]domain.SecurityLog, error) {
	var logs []domain.SecurityLog
	err := r.db.Order("timestamp desc").Limit(limit).Find(&logs).Error
	return logs, err
}
