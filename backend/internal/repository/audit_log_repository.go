package repository

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

// AuditLogRepository defines the interface for audit log data operations
type AuditLogRepository interface {
	Create(auditLog *model.AuditLog) error
	List(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error)
	GetByID(id uint) (*model.AuditLog, error)
	GetByPublicID(publicID string) (*model.AuditLog, error)
	GetByUserID(userID uint, limit int, offset int) ([]*model.AuditLog, int64, error)
	GetByResource(resource string, resourceID string) ([]*model.AuditLog, error)
	GetSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error)
	DeleteOldLogs(olderThan time.Time) error
}

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new instance of AuditLogRepository
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create saves a new audit log entry
func (r *auditLogRepository) Create(auditLog *model.AuditLog) error {
	return r.db.Create(auditLog).Error
}

// List returns paginated audit logs with filters
func (r *auditLogRepository) List(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{}).Preload("User").Preload("Actor")

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ActorID != nil {
		query = query.Where("actor_id = ?", *filter.ActorID)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.Resource != nil {
		query = query.Where("resource = ?", *filter.Resource)
	}
	if filter.ResourceID != nil {
		query = query.Where("resource_id = ?", *filter.ResourceID)
	}
	if filter.StartDate != nil {
		query = query.Where("timestamp >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("timestamp <= ?", *filter.EndDate)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	limit := 50 // default
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	if err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// GetByID retrieves an audit log by ID
func (r *auditLogRepository) GetByID(id uint) (*model.AuditLog, error) {
	var auditLog model.AuditLog
	if err := r.db.Preload("User").Preload("Actor").First(&auditLog, id).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// GetByPublicID retrieves an audit log by public ID
func (r *auditLogRepository) GetByPublicID(publicID string) (*model.AuditLog, error) {
	var auditLog model.AuditLog
	if err := r.db.Preload("User").Preload("Actor").Where("public_id = ?", publicID).First(&auditLog).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// GetByUserID retrieves audit logs for a specific user
func (r *auditLogRepository) GetByUserID(userID uint, limit int, offset int) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if limit <= 0 {
		limit = 50
	}

	if err := query.Preload("User").Preload("Actor").
		Order("timestamp DESC").
		Limit(limit).Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// GetByResource retrieves audit logs for a specific resource
func (r *auditLogRepository) GetByResource(resource string, resourceID string) ([]*model.AuditLog, error) {
	var auditLogs []*model.AuditLog
	if err := r.db.Preload("User").Preload("Actor").
		Where("resource = ? AND resource_id = ?", resource, resourceID).
		Order("timestamp DESC").
		Find(&auditLogs).Error; err != nil {
		return nil, err
	}
	return auditLogs, nil
}

// GetSummary generates summary statistics for audit logs
func (r *auditLogRepository) GetSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	summary := &model.AuditLogSummary{
		ActionsByType: make(map[string]int64),
		UserActivity:  make(map[string]int64),
	}

	query := r.db.Model(&model.AuditLog{})
	if startDate != nil {
		query = query.Where("timestamp >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("timestamp <= ?", *endDate)
	}

	// Total actions
	if err := query.Count(&summary.TotalActions).Error; err != nil {
		return nil, err
	}

	// Actions by type
	rows, err := query.Select("action, COUNT(*) as count").
		Group("action").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var action string
		var count int64
		if err := rows.Scan(&action, &count); err != nil {
			return nil, err
		}
		summary.ActionsByType[action] = count
	}

	// Recent actions (last 10)
	var recentActions []model.AuditLog
	if err := r.db.Order("timestamp DESC").Limit(10).Find(&recentActions).Error; err != nil {
		return nil, err
	}
	summary.RecentActions = recentActions

	// User activity (top 10 most active users)
	userActivityRows, err := r.db.Model(&model.AuditLog{}).
		Select("COALESCE(CAST(actor_id AS TEXT), 'system') as actor, COUNT(*) as count").
		Group("actor_id").
		Order("count DESC").
		Limit(10).Rows()
	if err != nil {
		return nil, err
	}
	defer userActivityRows.Close()

	for userActivityRows.Next() {
		var actor string
		var count int64
		if err := userActivityRows.Scan(&actor, &count); err != nil {
			return nil, err
		}
		summary.UserActivity[actor] = count
	}
	return summary, nil
}

// DeleteOldLogs removes audit logs older than the specified date (for data retention)
func (r *auditLogRepository) DeleteOldLogs(olderThan time.Time) error {
	return r.db.Unscoped().Where("timestamp < ?", olderThan).Delete(&model.AuditLog{}).Error
}
