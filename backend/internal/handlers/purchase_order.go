package handlers

import (
	"net/http"
	"strconv"
	"tatapps/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type POHandler struct {
	db *gorm.DB
}

func NewPOHandler(db *gorm.DB) *POHandler {
	return &POHandler{db: db}
}

func (h *POHandler) GetAll(c *gin.Context) {
	var pos []models.PurchaseOrder
	query := h.db.Preload("RequestedBy").Preload("ApprovedBy").Preload("Project").Preload("Warehouse").Preload("Items")

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&pos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pos)
}

func (h *POHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var po models.PurchaseOrder
	if err := h.db.Preload("RequestedBy").Preload("ApprovedBy").Preload("Project").Preload("Warehouse").Preload("Items").First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	c.JSON(http.StatusOK, po)
}

func (h *POHandler) Create(c *gin.Context) {
	var po models.PurchaseOrder
	if err := c.ShouldBindJSON(&po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set requested by from token
	userID := c.GetUint("user_id")
	po.RequestedByID = userID
	po.Status = "draft"

	// Calculate totals
	var subtotal float64
	for i := range po.Items {
		po.Items[i].TotalPrice = po.Items[i].Quantity * po.Items[i].UnitPrice
		subtotal += po.Items[i].TotalPrice
	}

	po.Subtotal = subtotal
	po.TaxAmount = (po.Subtotal * po.TaxPercent) / 100
	po.TotalAmount = po.Subtotal + po.TaxAmount - po.DiscountAmount + po.ShippingCost

	if err := h.db.Create(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, po)
}

func (h *POHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var po models.PurchaseOrder
	if err := h.db.First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	if err := c.ShouldBindJSON(&po); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, po)
}

func (h *POHandler) Approve(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	var po models.PurchaseOrder
	if err := h.db.First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	if po.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending PO can be approved"})
		return
	}

	po.Status = "approved"
	po.ApprovedByID = &userID
	now := time.Now()
	po.ApprovedAt = &now

	if err := h.db.Save(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order approved successfully",
		"po":      po,
	})
}

func (h *POHandler) Reject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var po models.PurchaseOrder
	if err := h.db.First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	po.Status = "rejected"
	po.RejectionReason = req.Reason

	if err := h.db.Save(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order rejected",
		"po":      po,
	})
}
