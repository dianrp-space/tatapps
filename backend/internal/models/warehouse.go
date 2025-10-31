package models

import (
	"time"

	"gorm.io/gorm"
)

type Warehouse struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Code       string `gorm:"uniqueIndex;not null" json:"code"`
	Name       string `gorm:"not null" json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	ManagerID  *uint  `json:"manager_id,omitempty"`
	Manager    *User  `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Color      string `gorm:"size:20;default:#6366F1" json:"color"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`

	// Relations
	InventoryItems []InventoryItem `gorm:"foreignKey:WarehouseID" json:"inventory_items,omitempty"`
}

type InventoryItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	WarehouseID uint      `gorm:"not null" json:"warehouse_id"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`

	SN          string  `gorm:"column:sku;not null" json:"sn"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Unit        string  `json:"unit"` // pcs, kg, box, etc
	Quantity    float64 `gorm:"default:0" json:"quantity"`
	MinStock    float64 `gorm:"default:0" json:"min_stock"`
	MaxStock    float64 `json:"max_stock"`
	UnitPrice   float64 `gorm:"default:0" json:"unit_price"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	// Relations
	Transactions []InventoryTransaction `gorm:"foreignKey:ItemID" json:"transactions,omitempty"`
}

type InventoryTransaction struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ItemID uint          `gorm:"not null" json:"item_id"`
	Item   InventoryItem `gorm:"foreignKey:ItemID" json:"item"`

	Type            string     `gorm:"not null" json:"type"` // in, out, transfer, adjustment
	Quantity        float64    `gorm:"not null" json:"quantity"`
	FromWarehouseID *uint      `json:"from_warehouse_id,omitempty"`
	FromWarehouse   *Warehouse `gorm:"foreignKey:FromWarehouseID" json:"from_warehouse,omitempty"`
	ToWarehouseID   *uint      `json:"to_warehouse_id,omitempty"`
	ToWarehouse     *Warehouse `gorm:"foreignKey:ToWarehouseID" json:"to_warehouse,omitempty"`

	Reference   string `json:"reference"` // PO number, transfer note, etc
	Notes       string `json:"notes"`
	CreatedByID uint   `gorm:"not null" json:"created_by_id"`
	CreatedBy   User   `gorm:"foreignKey:CreatedByID" json:"created_by"`
}
