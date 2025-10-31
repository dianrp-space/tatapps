package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	EmployeeCode    string     `gorm:"uniqueIndex;not null" json:"employee_code"`
	NIK             string     `gorm:"column:nik;uniqueIndex" json:"nik"`
	FullName        string     `gorm:"not null" json:"full_name"`
	Email           string     `gorm:"uniqueIndex" json:"email"`
	Phone           string     `json:"phone"`
	BirthDate       *time.Time `gorm:"column:birth_date" json:"date_of_birth,omitempty"`
	BirthPlace      string     `gorm:"column:birth_place" json:"birth_place"`
	Gender          string     `json:"gender"`
	BloodType       string     `gorm:"column:blood_type" json:"blood_type"`
	MaritalStatus   string     `gorm:"column:marital_status" json:"marital_status"`
	Religion        string     `json:"religion"`
	IdentityType    string     `gorm:"column:identity_type" json:"identity_type"`
	IdentityNumber  string     `gorm:"column:identity_number" json:"identity_number"`
	AddressKTP      string     `gorm:"column:address_ktp" json:"address_ktp"`
	AddressDomicile string     `gorm:"column:address_domicile" json:"address_domicile"`
	Address         string     `json:"address"`
	City            string     `json:"city"`
	Province        string     `json:"province"`
	PostalCode      string     `gorm:"column:postal_code" json:"postal_code"`
	Timezone        string     `json:"timezone"`
	Photo           string     `gorm:"type:text" json:"photo"`

	// Employment Details
	Department     string            `json:"department"`
	JobTitle       string            `gorm:"column:job_title" json:"job_title"`
	JoinDate       time.Time         `gorm:"not null;column:join_date" json:"join_date"`
	EmploymentType string            `gorm:"column:employment_type" json:"employment_type"` // full-time, part-time, contract
	Status         string            `gorm:"default:'active'" json:"status"`                // active, resigned, terminated
	DivisionID     *uint             `gorm:"column:division_id" json:"division_id"`
	Division       *EmployeeDivision `gorm:"foreignKey:DivisionID" json:"division,omitempty"`
	PositionID     *uint             `gorm:"column:position_id" json:"position_id"`
	Position       *EmployeePosition `gorm:"foreignKey:PositionID" json:"position,omitempty"`

	// Salary
	BasicSalary float64 `gorm:"default:0;column:basic_salary" json:"basic_salary"`

	// Emergency Contact
	EmergencyContactName  string `gorm:"column:emergency_contact_name" json:"emergency_contact_name"`
	EmergencyContactPhone string `gorm:"column:emergency_contact_phone" json:"emergency_contact_phone"`

	// Relations
	ManagerID   *uint      `json:"manager_id,omitempty"`
	Manager     *Employee  `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	WarehouseID *uint      `json:"warehouse_id,omitempty"`
	Warehouse   *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
}
