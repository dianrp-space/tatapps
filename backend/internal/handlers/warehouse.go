package handlers

import (
	"net/http"
	"strconv"
	"tatapps/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseHandler struct {
	db *gorm.DB
}

func NewWarehouseHandler(db *gorm.DB) *WarehouseHandler {
	return &WarehouseHandler{db: db}
}

func (h *WarehouseHandler) GetAll(c *gin.Context) {
	var warehouses []models.Warehouse
	if err := h.db.Preload("Manager").Find(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch warehouses",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": warehouses,
	})
}

func (h *WarehouseHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var warehouse models.Warehouse
	if err := h.db.Preload("Manager").First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Warehouse not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": warehouse,
	})
}

func (h *WarehouseHandler) Create(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preload manager to return complete data
	h.db.Preload("Manager").First(&warehouse, warehouse.ID)

	c.JSON(http.StatusCreated, gin.H{"data": warehouse})
}

func (h *WarehouseHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var warehouse models.Warehouse
	if err := h.db.First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preload manager to return complete data
	h.db.Preload("Manager").First(&warehouse, warehouse.ID)

	c.JSON(http.StatusOK, gin.H{"data": warehouse})
}

func (h *WarehouseHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.db.Delete(&models.Warehouse{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
}
