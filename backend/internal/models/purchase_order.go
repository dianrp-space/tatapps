package models

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseOrder struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// PO Information
	PONumber        string    `gorm:"uniqueIndex;not null" json:"po_number"`
	PODate          time.Time `gorm:"not null" json:"po_date"`
	
	// Supplier Information
	SupplierName    string    `gorm:"not null" json:"supplier_name"`
	SupplierEmail   string    `json:"supplier_email"`
	SupplierPhone   string    `json:"supplier_phone"`
	SupplierAddress string    `json:"supplier_address"`
	
	// PO Details
	Status          string    `gorm:"default:'draft'" json:"status"` // draft, pending, approved, rejected, ordered, received, cancelled
	Priority        string    `gorm:"default:'medium'" json:"priority"` // low, medium, high
	DeliveryDate    *time.Time `json:"delivery_date,omitempty"`
	DeliveryAddress string    `json:"delivery_address"`
	
	// Financial
	Subtotal        float64   `gorm:"default:0" json:"subtotal"`
	TaxPercent      float64   `gorm:"default:0" json:"tax_percent"`
	TaxAmount       float64   `gorm:"default:0" json:"tax_amount"`
	DiscountAmount  float64   `gorm:"default:0" json:"discount_amount"`
	ShippingCost    float64   `gorm:"default:0" json:"shipping_cost"`
	TotalAmount     float64   `gorm:"default:0" json:"total_amount"`
	
	// Assignment & Approval
	RequestedByID   uint      `gorm:"not null" json:"requested_by_id"`
	RequestedBy     User      `gorm:"foreignKey:RequestedByID" json:"requested_by"`
	ApprovedByID    *uint     `json:"approved_by_id,omitempty"`
	ApprovedBy      *User     `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	
	// Relations
	ProjectID       *uint     `json:"project_id,omitempty"`
	Project         *Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	WarehouseID     *uint     `json:"warehouse_id,omitempty"`
	Warehouse       *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	
	Items           []POItem  `gorm:"foreignKey:POID" json:"items"`
	
	// Notes
	Notes           string    `json:"notes"`
	RejectionReason string    `json:"rejection_reason"`
}

type POItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	POID         uint    `gorm:"not null" json:"po_id"`
	PO           PurchaseOrder `gorm:"foreignKey:POID" json:"-"`
	
	ItemName     string  `gorm:"not null" json:"item_name"`
	ItemCode     string  `json:"item_code"`
	Description  string  `json:"description"`
	Unit         string  `json:"unit"` // pcs, kg, box, etc
	Quantity     float64 `gorm:"not null" json:"quantity"`
	UnitPrice    float64 `gorm:"not null" json:"unit_price"`
	TotalPrice   float64 `gorm:"not null" json:"total_price"`
	
	ReceivedQty  float64 `gorm:"default:0" json:"received_qty"`
	Notes        string  `json:"notes"`
}
