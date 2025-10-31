package models

import (
	"time"

	"gorm.io/gorm"
)

// UserWarehouse defines the warehouses a user is allowed to access.
type UserWarehouse struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint      `gorm:"index:idx_user_warehouse,unique;not null" json:"user_id"`
	WarehouseID uint      `gorm:"index:idx_user_warehouse,unique;not null" json:"warehouse_id"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
}
