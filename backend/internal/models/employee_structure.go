package models

import (
	"time"

	"gorm.io/gorm"
)

type EmployeeDivision struct {
	ID                uint               `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         gorm.DeletedAt     `gorm:"index" json:"-"`
	Name              string             `gorm:"not null" json:"name"`
	Description       string             `json:"description"`
	RecruitmentStatus string             `gorm:"default:'Stabil'" json:"recruitment_status"`
	HeadEmployeeID    *uint              `json:"head_employee_id"`
	HeadPositionID    *uint              `json:"head_position_id"`
	Head              string             `json:"head"`
	HeadTitle         string             `json:"head_title"`
	HeadEmployee      *Employee          `gorm:"-" json:"head_employee,omitempty"`
	HeadPosition      *EmployeePosition  `gorm:"-" json:"head_position,omitempty"`
	Employees         []Employee         `gorm:"-" json:"-"`
	Positions         []EmployeePosition `gorm:"-" json:"-"`
}

type EmployeePosition struct {
	ID          uint               `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
	Title       string             `gorm:"not null" json:"title"`
	Code        string             `json:"code"`
	DivisionID  *uint              `json:"division_id"`
	ParentID    *uint              `json:"parent_id"`
	Notes       string             `json:"notes"`
	Grade       string             `json:"grade"`
	SalaryRange string             `json:"salary_range"`
	Division    *EmployeeDivision  `gorm:"foreignKey:DivisionID" json:"division,omitempty"`
	Parent      *EmployeePosition  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []EmployeePosition `gorm:"foreignKey:ParentID" json:"-"`
	Employees   []Employee         `gorm:"foreignKey:PositionID" json:"-"`
}
