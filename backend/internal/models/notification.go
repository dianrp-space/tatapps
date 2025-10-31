package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Notification Details
	Type    string `gorm:"not null" json:"type"` // email, whatsapp, system
	Channel string `json:"channel"`              // email, whatsapp, both
	Module  string `json:"module"`               // warehouse, employee, lead, project, po

	// Recipient
	UserID *uint  `json:"user_id,omitempty"`
	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`

	// Content
	Subject string `json:"subject"`
	Message string `json:"message"`
	Data    string `gorm:"type:jsonb" json:"data"` // JSON data for additional info

	// Status
	Status   string     `gorm:"default:'pending'" json:"status"` // pending, sent, failed, read
	SentAt   *time.Time `json:"sent_at,omitempty"`
	ReadAt   *time.Time `json:"read_at,omitempty"`
	ErrorMsg string     `json:"error_msg"`

	// Retry
	RetryCount int `gorm:"default:0" json:"retry_count"`
	MaxRetries int `gorm:"default:3" json:"max_retries"`
}

// NotificationSetting stores user's notification preferences
type NotificationSetting struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	// Low Stock Alert Settings
	Enabled        bool   `gorm:"default:false" json:"enabled"`
	Threshold      int    `gorm:"default:10" json:"threshold"`
	CheckFrequency string `gorm:"default:'daily'" json:"check_frequency"` // hourly, daily, weekly
	ScheduleMode   string `gorm:"default:'preset'" json:"schedule_mode"`
	CronExpression string `json:"cron_expression"`
	TimeZone       string `gorm:"default:'Asia/Jakarta'" json:"timezone"`

	// Notification Channels
	WhatsAppEnabled bool       `gorm:"default:false" json:"whatsapp_enabled"`
	WhatsAppNumber  string     `json:"whatsapp_number"`
	EmailEnabled    bool       `gorm:"default:false" json:"email_enabled"`
	EmailAddress    string     `json:"email_address"`
	LastRunAt       *time.Time `json:"last_run_at"`
}

// NotificationHistory stores sent notifications for auditing
type NotificationHistory struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	Type    string `json:"type"` // low_stock, test, critical
	Title   string `json:"title"`
	Message string `json:"message"`

	WhatsAppSent bool `gorm:"default:false" json:"whatsapp_sent"`
	EmailSent    bool `gorm:"default:false" json:"email_sent"`
}
