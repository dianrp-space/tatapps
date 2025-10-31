package models

import (
	"time"

	"gorm.io/gorm"
)

type Lead struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Lead Information
	CompanyName     string  `json:"company_name"`
	ContactPerson   string  `gorm:"not null" json:"contact_person"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	WhatsApp        string  `json:"whatsapp"`
	
	// Address
	Address         string  `json:"address"`
	City            string  `json:"city"`
	Province        string  `json:"province"`
	
	// Lead Details
	Source          string  `json:"source"` // website, referral, cold-call, etc
	Industry        string  `json:"industry"`
	Status          string  `gorm:"default:'new'" json:"status"` // new, contacted, qualified, proposal, negotiation, won, lost
	Priority        string  `gorm:"default:'medium'" json:"priority"` // low, medium, high
	EstimatedValue  float64 `gorm:"default:0" json:"estimated_value"`
	EstimatedCloseDate *time.Time `json:"estimated_close_date,omitempty"`
	
	// Assignment
	AssignedToID    uint    `gorm:"not null" json:"assigned_to_id"`
	AssignedTo      User    `gorm:"foreignKey:AssignedToID" json:"assigned_to"`
	
	// Notes & Follow-up
	Notes           string  `json:"notes"`
	LastContactDate *time.Time `json:"last_contact_date,omitempty"`
	NextFollowUpDate *time.Time `json:"next_follow_up_date,omitempty"`
	
	// Conversion
	ConvertedToProjectID *uint    `json:"converted_to_project_id,omitempty"`
	ConvertedToProject   *Project `gorm:"foreignKey:ConvertedToProjectID" json:"converted_to_project,omitempty"`
	ConvertedAt          *time.Time `json:"converted_at,omitempty"`
}
