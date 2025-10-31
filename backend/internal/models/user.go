package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	FullName string `gorm:"not null" json:"full_name"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role"`

	Warehouses []UserWarehouse `gorm:"foreignKey:UserID" json:"warehouses,omitempty"`
}

type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string       `gorm:"uniqueIndex;not null" json:"name"` // admin, employee, manager, etc
	Description string       `json:"description"`
	Color       string       `gorm:"size:20;default:#2563EB" json:"color"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	Menus       []RoleMenu   `gorm:"foreignKey:RoleID" json:"menus"`
}

type Permission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"uniqueIndex;not null" json:"name"` // warehouse.view, warehouse.create, etc
	Description string `json:"description"`
	Module      string `json:"module"` // warehouse, employee, lead, project, po
	Action      string `json:"action"` // view, create, update, delete
}

type RoleMenu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RoleID  uint   `gorm:"index:idx_role_menu,unique;not null" json:"role_id"`
	MenuKey string `gorm:"size:100;index:idx_role_menu,unique;not null" json:"menu_key"`
}
