package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionUserLogin         AuditAction = "user_login"
	AuditActionUserLogout        AuditAction = "user_logout"
	AuditActionUserCreated       AuditAction = "user_created"
	AuditActionUserUpdated       AuditAction = "user_updated"
	AuditActionUserDeactivated   AuditAction = "user_deactivated"
	AuditActionUserReactivated   AuditAction = "user_reactivated"
	AuditActionUserDeleted       AuditAction = "user_deleted"
	AuditActionUserExported      AuditAction = "user_exported"
	AuditActionRoleChanged       AuditAction = "role_changed"
	AuditActionDataAccessed      AuditAction = "data_accessed"
	AuditActionDeletionRequested AuditAction = "deletion_requested"
	AuditActionDeletionConfirmed AuditAction = "deletion_confirmed"
	AuditActionDeletionCancelled AuditAction = "deletion_cancelled"
	AuditActionDataAnonymized    AuditAction = "data_anonymized"
)

// AuditLog represents an audit log entry for LGPD compliance
type AuditLog struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	PublicID   uuid.UUID   `json:"public_id" gorm:"type:uuid;uniqueIndex"`
	UserID     *uint       `json:"user_id,omitempty" gorm:"index"`           // User being acted upon (nullable for system actions)
	ActorID    *uint       `json:"actor_id,omitempty" gorm:"index"`          // User performing the action (nullable for system actions)
	ActorType  string      `json:"actor_type" gorm:"size:50;default:'user'"` // 'user', 'admin', 'system'
	Action     AuditAction `json:"action" gorm:"size:50;not null"`
	Resource   string      `json:"resource" gorm:"size:100;not null"`           // 'user', 'order', 'product', etc.
	ResourceID *string     `json:"resource_id,omitempty" gorm:"size:100;index"` // ID of the resource
	Details    string      `json:"details" gorm:"type:text"`                    // JSON string with additional details
	OldValues  string      `json:"old_values,omitempty" gorm:"type:text"`       // JSON string of old values (for updates)
	NewValues  string      `json:"new_values,omitempty" gorm:"type:text"`       // JSON string of new values (for updates)
	Timestamp  time.Time   `json:"timestamp" gorm:"not null;index"`
	Compliance string      `json:"compliance,omitempty" gorm:"size:100"`   // LGPD, GDPR, etc.
	Severity   string      `json:"severity" gorm:"size:20;default:'info'"` // 'info', 'warning', 'error', 'critical'
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`

	// Relations
	User  *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Actor *User `json:"actor,omitempty" gorm:"foreignKey:ActorID"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate hook to generate UUID
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.PublicID == uuid.Nil {
		a.PublicID = uuid.New()
	}
	return nil
}

// AuditLogFilter represents filters for audit log queries
type AuditLogFilter struct {
	UserID     *uint        `json:"user_id,omitempty"`
	ActorID    *uint        `json:"actor_id,omitempty"`
	Action     *AuditAction `json:"action,omitempty"`
	Resource   *string      `json:"resource,omitempty"`
	ResourceID *string      `json:"resource_id,omitempty"`
	StartDate  *time.Time   `json:"start_date,omitempty"`
	EndDate    *time.Time   `json:"end_date,omitempty"`
	Severity   *string      `json:"severity,omitempty"`
	Limit      int          `json:"limit,omitempty" validate:"min=1,max=1000"`
	Offset     int          `json:"offset,omitempty" validate:"min=0"`
}

// AuditLogSummary represents audit log summary statistics
type AuditLogSummary struct {
	TotalActions  int64            `json:"total_actions"`
	ActionsByType map[string]int64 `json:"actions_by_type"`
	RecentActions []AuditLog       `json:"recent_actions"`
	UserActivity  map[string]int64 `json:"user_activity"`
}
