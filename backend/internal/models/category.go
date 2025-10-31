package models

import (
	"time"
)

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;unique" json:"name"`
	Code        string    `gorm:"size:20;not null;unique" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"size:20;not null;default:#3B82F6" json:"color"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	ItemCount   int       `gorm:"-" json:"item_count"` // Virtual field for count
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
