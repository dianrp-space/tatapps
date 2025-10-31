package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Project Information
	ProjectCode     string    `gorm:"uniqueIndex;not null" json:"project_code"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	
	// Client Information
	ClientName      string    `json:"client_name"`
	ClientEmail     string    `json:"client_email"`
	ClientPhone     string    `json:"client_phone"`
	ClientAddress   string    `json:"client_address"`
	
	// Project Details
	Status          string    `gorm:"default:'planning'" json:"status"` // planning, in-progress, on-hold, completed, cancelled
	Priority        string    `gorm:"default:'medium'" json:"priority"` // low, medium, high
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	
	// Financial
	Budget          float64   `gorm:"default:0" json:"budget"`
	ActualCost      float64   `gorm:"default:0" json:"actual_cost"`
	
	// Assignment
	ManagerID       uint      `gorm:"not null" json:"manager_id"`
	Manager         User      `gorm:"foreignKey:ManagerID" json:"manager"`
	
	// Relations
	PurchaseOrders  []PurchaseOrder `gorm:"foreignKey:ProjectID" json:"purchase_orders,omitempty"`
	
	// Progress
	ProgressPercent int       `gorm:"default:0" json:"progress_percent"`
	Notes           string    `json:"notes"`
}
