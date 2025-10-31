package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
	"tatapps/internal/models"
	"tatapps/internal/services/notification"
)

type InventoryHandler struct {
	db *gorm.DB
}

func NewInventoryHandler(db *gorm.DB) *InventoryHandler {
	return &InventoryHandler{db: db}
}

func toUint(value any) (uint, bool) {
	switch v := value.(type) {
	case uint:
		return v, true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case uint64:
		return uint(v), true
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	}
	return 0, false
}

func buildUintSet(ids []uint) map[uint]struct{} {
	set := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		set[id] = struct{}{}
	}
	return set
}

func (h *InventoryHandler) contextUserID(c *gin.Context) (uint, bool) {
	if id := c.GetUint("user_id"); id != 0 {
		return id, true
	}
	value, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := toUint(value)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}

func (h *InventoryHandler) getUserWarehouseIDs(userID uint) ([]uint, error) {
	var ids []uint
	if err := h.db.Model(&models.UserWarehouse{}).
		Where("user_id = ?", userID).
		Pluck("warehouse_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func filterTransactionsByWarehouses(transactions []models.InventoryTransaction, allowed map[uint]struct{}) []models.InventoryTransaction {
	if len(allowed) == 0 {
		return []models.InventoryTransaction{}
	}

	filtered := make([]models.InventoryTransaction, 0, len(transactions))
	for _, tx := range transactions {
		if _, ok := allowed[tx.Item.WarehouseID]; ok {
			filtered = append(filtered, tx)
			continue
		}
		if tx.FromWarehouseID != nil {
			if _, ok := allowed[*tx.FromWarehouseID]; ok {
				filtered = append(filtered, tx)
				continue
			}
		}
		if tx.ToWarehouseID != nil {
			if _, ok := allowed[*tx.ToWarehouseID]; ok {
				filtered = append(filtered, tx)
				continue
			}
		}
	}
	return filtered
}

// GetAllItems godoc
// @Summary Get all inventory items
// @Tags Inventory
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /inventory [get]
func (h *InventoryHandler) GetAllItems(c *gin.Context) {
	var items []models.InventoryItem

	query := applyInventoryFilters(h.db.Preload("Warehouse"), c)

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch inventory items",
			"message": err.Error(),
		})
		return
	}

	if roleName, ok := c.Get("role_name"); ok {
		if name, isString := roleName.(string); isString && strings.EqualFold(name, "employee") {
			userID, ok := h.contextUserID(c)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unable to resolve user context"})
				return
			}

			allowedIDs, err := h.getUserWarehouseIDs(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve warehouse permissions"})
				return
			}
			if len(allowedIDs) == 0 {
				items = []models.InventoryItem{}
			} else {
				allowedSet := buildUintSet(allowedIDs)
				filtered := make([]models.InventoryItem, 0, len(items))
				for _, item := range items {
					if _, ok := allowedSet[item.WarehouseID]; ok {
						filtered = append(filtered, item)
					}
				}
				items = filtered
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

// GetItemByID godoc
// @Summary Get inventory item by ID
// @Tags Inventory
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Router /inventory/{id} [get]
func (h *InventoryHandler) GetItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.InventoryItem
	if err := h.db.Preload("Warehouse").First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch item",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// CreateItem godoc
// @Summary Create new inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Param item body models.InventoryItem true "Item data"
// @Success 201 {object} map[string]interface{}
// @Router /inventory [post]
func (h *InventoryHandler) CreateItem(c *gin.Context) {
	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	item.SN = strings.TrimSpace(item.SN)

	// Check if SN already exists
	if item.SN != "" {
		var existing models.InventoryItem
		if err := h.db.Where("sku = ?", item.SN).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "SN already exists",
			})
			return
		}
	}

	// Set IsActive based on initial quantity
	if item.Quantity > 0 {
		item.IsActive = true
	} else {
		item.IsActive = false
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create item",
			"message": err.Error(),
		})
		return
	}

	// Reload with warehouse
	h.db.Preload("Warehouse").First(&item, item.ID)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "Item created successfully",
	})
}

// UpdateItem godoc
// @Summary Update inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body models.InventoryItem true "Item data"
// @Success 200 {object} map[string]interface{}
// @Router /inventory/{id} [put]

type updateInventoryItemRequest struct {
	SN          *string  `json:"sn"`
	Name        *string  `json:"name"`
	Category    *string  `json:"category"`
	WarehouseID *uint    `json:"warehouse_id"`
	Quantity    *float64 `json:"quantity"`
	MinStock    *float64 `json:"min_stock"`
	MaxStock    *float64 `json:"max_stock"`
	UnitPrice   *float64 `json:"unit_price"`
	Unit        *string  `json:"unit"`
	Description *string  `json:"description"`
	IsActive    *bool    `json:"is_active"`
}

func (h *InventoryHandler) UpdateItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.InventoryItem
	if err := h.db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch item",
			"message": err.Error(),
		})
		return
	}

	var req updateInventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	if req.SN != nil {
		trimmed := strings.TrimSpace(*req.SN)
		if trimmed != "" && trimmed != item.SN {
			var existing models.InventoryItem
			if err := h.db.Where("sku = ? AND id != ?", trimmed, id).First(&existing).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "SN already exists",
				})
				return
			}
		}
		item.SN = trimmed
	}

	if req.Name != nil {
		item.Name = strings.TrimSpace(*req.Name)
	}
	if req.Category != nil {
		item.Category = strings.TrimSpace(*req.Category)
	}
	if req.WarehouseID != nil {
		item.WarehouseID = *req.WarehouseID
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.MinStock != nil {
		item.MinStock = *req.MinStock
	}
	if req.MaxStock != nil {
		item.MaxStock = *req.MaxStock
	}
	if req.UnitPrice != nil {
		item.UnitPrice = *req.UnitPrice
	}
	if req.Unit != nil {
		item.Unit = strings.TrimSpace(*req.Unit)
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	statusProvided := req.IsActive != nil
	if statusProvided && req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if !statusProvided {
		if item.Quantity == 0 {
			item.IsActive = false
		} else if item.Quantity > 0 {
			item.IsActive = true
		}
	}

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update item",
			"message": err.Error(),
		})
		return
	}

	if err := h.db.Preload("Warehouse").First(&item, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reload item",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Item updated successfully",
	})
}

// DeleteItem godoc
// @Summary Delete inventory item
// @Tags Inventory
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Router /inventory/{id} [delete]
func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.InventoryItem
	if err := h.db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch item",
			"message": err.Error(),
		})
		return
	}

	if err := h.db.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete item",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
	})
}

// DeleteItemsBatch deletes multiple inventory items by IDs
func (h *InventoryHandler) DeleteItemsBatch(c *gin.Context) {
	var payload struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	if len(payload.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No item IDs provided",
		})
		return
	}

	if err := h.db.Where("id IN ?", payload.IDs).Delete(&models.InventoryItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete items",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d items deleted successfully", len(payload.IDs)),
		"deleted": payload.IDs,
	})
}

// GetItemTransactions godoc
// @Summary Get transactions for an item
// @Tags Inventory
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Router /inventory/{id}/transactions [get]
func (h *InventoryHandler) GetItemTransactions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if roleName, ok := c.Get("role_name"); ok {
		if name, isString := roleName.(string); isString && strings.EqualFold(name, "employee") {
			userID, ok := h.contextUserID(c)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unable to resolve user context"})
				return
			}

			allowedIDs, err := h.getUserWarehouseIDs(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve warehouse permissions"})
				return
			}
			if len(allowedIDs) == 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "No warehouse access configured for this account"})
				return
			}

			allowedSet := buildUintSet(allowedIDs)
			var item models.InventoryItem
			if err := h.db.Select("id", "warehouse_id").First(&item, id).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
				return
			}
			if _, ok := allowedSet[item.WarehouseID]; !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to access transactions for this warehouse"})
				return
			}
		}
	}

	var transactions []models.InventoryTransaction
	if err := h.db.
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("CreatedBy").
		Where("item_id = ?", id).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch transactions",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// GetAllTransactions godoc
// @Summary Get all inventory transactions
// @Tags Inventory
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /inventory/transactions [get]
func (h *InventoryHandler) GetAllTransactions(c *gin.Context) {
	var transactions []models.InventoryTransaction

	query := applyTransactionFilters(h.db.Model(&models.InventoryTransaction{}), c)

	if err := query.
		Preload("Item").
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch transactions",
			"message": err.Error(),
		})
		return
	}

	if roleName, ok := c.Get("role_name"); ok {
		if name, isString := roleName.(string); isString && strings.EqualFold(name, "employee") {
			userID, ok := h.contextUserID(c)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unable to resolve user context"})
				return
			}
			allowedIDs, err := h.getUserWarehouseIDs(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve warehouse permissions"})
				return
			}
			if len(allowedIDs) == 0 {
				transactions = []models.InventoryTransaction{}
			} else {
				allowedSet := buildUintSet(allowedIDs)
				transactions = filterTransactionsByWarehouses(transactions, allowedSet)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// RecordTransaction godoc
// @Summary Record inventory transaction
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param transaction body models.InventoryTransaction true "Transaction data"
// @Success 201 {object} map[string]interface{}
// @Router /inventory/items/{id}/transactions [post]
func (h *InventoryHandler) RecordTransaction(c *gin.Context) {
	// Get item ID from URL param
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid item ID",
			"message": err.Error(),
		})
		return
	}

	// Debug logging
	println("=== RECORD TRANSACTION DEBUG ===")
	println("Item ID from URL:", itemID)

	var transaction models.InventoryTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
		return
	}

	// Debug logging
	println("Transaction type:", transaction.Type)
	println("Quantity:", transaction.Quantity)

	// Set item ID from URL param
	transaction.ItemID = uint(itemID)

	// Get user context
	userID, ok := h.contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	transaction.CreatedByID = userID

	roleName := ""
	if roleValue, exists := c.Get("role_name"); exists {
		if name, isString := roleValue.(string); isString {
			roleName = name
		}
	}

	isEmployee := strings.EqualFold(roleName, "employee")
	var allowedSet map[uint]struct{}
	if isEmployee {
		allowedIDs, err := h.getUserWarehouseIDs(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve warehouse permissions"})
			return
		}
		if len(allowedIDs) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "No warehouse access configured for this account"})
			return
		}
		allowedSet = buildUintSet(allowedIDs)
	}

	// First, check if item exists (before starting DB transaction)
	var item models.InventoryItem
	println("Checking if item exists with ID:", transaction.ItemID)
	if err := h.db.Preload("Warehouse").First(&item, transaction.ItemID).Error; err != nil {
		println("Error finding item:", err.Error())
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Item not found",
				"message": "Inventory item with ID " + strconv.Itoa(itemID) + " does not exist",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch item",
				"message": err.Error(),
			})
		}
		return
	}

	println("Item found:", item.Name, "SN:", item.SN, "Current quantity:", item.Quantity)

	if isEmployee {
		if _, exists := allowedSet[item.WarehouseID]; !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to record transactions for this warehouse"})
			return
		}
		if transaction.Type == "transfer" {
			if transaction.ToWarehouseID == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Destination warehouse is required for transfer"})
				return
			}
			if _, exists := allowedSet[*transaction.ToWarehouseID]; !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Destination warehouse is not permitted"})
				return
			}
		}
	}

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the item again within transaction
	if err := tx.Preload("Warehouse").First(&item, transaction.ItemID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to lock item",
			"message": err.Error(),
		})
		return
	}

	// Validate and update quantity based on transaction type
	switch transaction.Type {
	case "in":
		item.Quantity += transaction.Quantity

	case "out":
		if item.Quantity < transaction.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock",
			})
			return
		}
		item.Quantity -= transaction.Quantity

	case "transfer":
		if transaction.ToWarehouseID == nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Destination warehouse is required for transfer",
			})
			return
		}

		if *transaction.ToWarehouseID == item.WarehouseID {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot transfer to the same warehouse",
			})
			return
		}

		if item.Quantity < transaction.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock for transfer",
			})
			return
		}

		// Reduce from source warehouse
		item.Quantity -= transaction.Quantity
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update source item",
				"message": err.Error(),
			})
			return
		}

		// Find or create item in destination warehouse
		var destItem models.InventoryItem
		err := tx.Where("sku = ? AND warehouse_id = ?", item.SN, *transaction.ToWarehouseID).
			Preload("Warehouse").
			First(&destItem).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			destItem = models.InventoryItem{
				WarehouseID: *transaction.ToWarehouseID,
				SN:          item.SN,
				Name:        item.Name,
				Description: item.Description,
				Category:    item.Category,
				Unit:        item.Unit,
				Quantity:    transaction.Quantity,
				MinStock:    item.MinStock,
				MaxStock:    item.MaxStock,
				UnitPrice:   item.UnitPrice,
				IsActive:    transaction.Quantity > 0,
			}
			if err := tx.Create(&destItem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to create destination item",
					"message": err.Error(),
				})
				return
			}
		case err == nil:
			destItem.Name = item.Name
			destItem.Description = item.Description
			destItem.Category = item.Category
			destItem.Unit = item.Unit
			destItem.MinStock = item.MinStock
			destItem.MaxStock = item.MaxStock
			destItem.UnitPrice = item.UnitPrice

			destItem.Quantity += transaction.Quantity
			if destItem.Quantity > 0 {
				destItem.IsActive = true
			}

			if err := tx.Save(&destItem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to update destination item",
					"message": err.Error(),
				})
				return
			}
		default:
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check destination item",
				"message": err.Error(),
			})
			return
		}

		transaction.FromWarehouseID = &item.WarehouseID

	case "adjustment":
		// For adjustment, quantity can be positive (add) or negative (reduce)
		newQuantity := item.Quantity + transaction.Quantity
		if newQuantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Adjustment would result in negative stock",
			})
			return
		}
		item.Quantity = newQuantity

	default:
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction type",
		})
		return
	}

	// Update item quantity
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update item quantity",
			"message": err.Error(),
		})
		return
	}

	// Auto-update IsActive based on quantity
	if item.Quantity == 0 {
		item.IsActive = false
	} else if item.Quantity > 0 {
		item.IsActive = true
	}
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update item status",
			"message": err.Error(),
		})
		return
	}

	// Create transaction record
	transaction.CreatedAt = time.Now()
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to record transaction",
			"message": err.Error(),
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to commit transaction",
			"message": err.Error(),
		})
		return
	}

	// Reload with relations
	h.db.
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("CreatedBy").
		First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, gin.H{
		"data":    transaction,
		"message": "Transaction recorded successfully",
	})
}

// DeleteTransaction removes a transaction and reverts the stock changes it introduced.
// @Summary Delete inventory transaction
// @Tags Inventory
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]interface{}
// @Router /inventory/transactions/{id} [delete]
func (h *InventoryHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction models.InventoryTransaction
	if err := h.db.Preload("Item").First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch transaction",
			"message": err.Error(),
		})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var item models.InventoryItem
	itemResult := tx.Unscoped().First(&item, transaction.ItemID)
	if itemResult.Error != nil {
		if errors.Is(itemResult.Error, gorm.ErrRecordNotFound) {
			// Item no longer exists; delete transaction without stock rollback.
			if err := tx.Delete(&models.InventoryTransaction{}, transaction.ID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete transaction",
					"message": err.Error(),
				})
				return
			}
			if err := tx.Commit().Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to commit delete transaction",
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Transaction deleted (related inventory item already removed)",
			})
			return
		}
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to lock inventory item",
			"message": itemResult.Error.Error(),
		})
		return
	}

	switch transaction.Type {
	case "in":
		if item.Quantity < transaction.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot delete transaction because stock would become negative",
			})
			return
		}
		item.Quantity -= transaction.Quantity

	case "out":
		item.Quantity += transaction.Quantity

	case "transfer":
		item.Quantity += transaction.Quantity

		if transaction.ToWarehouseID != nil {
			var destItem models.InventoryItem
			destResult := tx.Unscoped().Where("sku = ? AND warehouse_id = ?", item.SN, *transaction.ToWarehouseID).First(&destItem)
			if destResult.Error != nil {
				if !errors.Is(destResult.Error, gorm.ErrRecordNotFound) {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to fetch destination inventory item",
						"message": destResult.Error.Error(),
					})
					return
				}
				// Destination item missing; continue without rollback.
			} else {
				if destItem.Quantity < transaction.Quantity {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Destination stock is lower than transferred quantity",
					})
					return
				}

				destItem.Quantity -= transaction.Quantity
				if err := tx.Save(&destItem).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to update destination inventory item",
						"message": err.Error(),
					})
					return
				}
			}
		}

	case "adjustment":
		item.Quantity -= transaction.Quantity
		if item.Quantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Deleting this adjustment would result in negative stock",
			})
			return
		}

	default:
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported transaction type",
		})
		return
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update inventory item",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Delete(&models.InventoryTransaction{}, transaction.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete transaction",
			"message": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to commit delete transaction",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})
}

// GetLowStockItems godoc
// @Summary Get items with low stock
// @Tags Inventory
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /inventory/low-stock [get]
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	var allItems []models.InventoryItem
	if err := h.db.
		Preload("Warehouse").
		Where("is_active = ?", true).
		Find(&allItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch low stock items",
			"message": err.Error(),
		})
		return
	}

	lowStock := notification.ComputeLowStockEntries(allItems)

	c.JSON(http.StatusOK, gin.H{
		"data": lowStock,
	})
}

// ImportItemsFromCSV handles bulk inventory import from CSV file
func (h *InventoryHandler) ImportItemsFromCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "CSV file is required",
			"message": err.Error(),
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Unable to open uploaded file",
			"message": err.Error(),
		})
		return
	}
	defer src.Close()

	reader := csv.NewReader(src)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err == io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "CSV file is empty",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to read CSV header",
			"message": err.Error(),
		})
		return
	}

	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	headerMap := make(map[string]int, len(header))
	warehouseCodeKey := ""
	warehouseNameKey := ""
	for idx, column := range header {
		normalized := strings.ToLower(strings.TrimSpace(column))
		switch normalized {
		case "sn", "serial number", "serial", "sku", "kode sn":
			headerMap["sn"] = idx
		case "name", "item name", "nama item", "nama barang", "item":
			headerMap["name"] = idx
		case "category", "kategori":
			headerMap["category"] = idx
		case "warehouse code", "warehouse_code", "kode_gudang", "kode warehouse", "kode gudang":
			warehouseCodeKey = normalized
			headerMap["warehouse_code"] = idx
		case "warehouse name", "warehouse", "nama_gudang", "nama warehouse", "gudang":
			warehouseNameKey = normalized
			headerMap["warehouse_name"] = idx
		case "quantity", "qty", "stok", "stock", "kuantitas", "kuantiti":
			headerMap["quantity"] = idx
		case "min stock", "minimum stock", "minstok", "minimum stok", "min stok":
			headerMap["min_stock"] = idx
		case "max stock", "maximum stock", "maxstok", "maks stok":
			headerMap["max_stock"] = idx
		case "unit", "satuan":
			headerMap["unit"] = idx
		case "unit price", "price", "harga":
			headerMap["unit_price"] = idx
		case "description", "deskripsi", "keterangan":
			headerMap["description"] = idx
		case "status", "aktif", "status item", "status aktif":
			headerMap["status"] = idx
		default:
			if normalized != "" {
				headerMap[normalized] = idx
			}
		}
	}

	requiredColumns := []string{"sn", "name"}
	for _, column := range requiredColumns {
		if _, ok := headerMap[column]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Missing required column '%s' in CSV", column),
			})
			return
		}
	}

	// Check if either warehouse_code or warehouse_name is present
	if warehouseCodeKey == "" && warehouseNameKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required column for warehouse (e.g. 'warehouse_code', 'Warehouse Code', 'warehouse_name', 'Warehouse Name') in CSV",
		})
		return
	}

	type importSummary struct {
		Inserted int      `json:"inserted"`
		Updated  int      `json:"updated"`
		Errors   []string `json:"errors"`
	}

	summary := importSummary{}

	tx := h.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to start database transaction",
			"message": tx.Error.Error(),
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	lineNumber := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNumber++

		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: %v", lineNumber, err))
			continue
		}

		if isCSVRecordEmpty(record) {
			continue
		}

		var warehouse models.Warehouse

		if idx, ok := headerMap["warehouse_code"]; ok && idx < len(record) {
			if code := strings.TrimSpace(record[idx]); code != "" {
				if err := tx.Where("code = ?", code).First(&warehouse).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: warehouse with code '%s' not found", lineNumber, code))
						continue
					}
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to fetch warehouse by code",
						"message": err.Error(),
					})
					return
				}
			}
		}

		if warehouse.ID == 0 {
			warehouseName := csvValue(record, headerMap, "warehouse_name")
			if warehouseName == "" {
				summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: warehouse name is required", lineNumber))
				continue
			}

			if err := tx.Where("name = ?", warehouseName).First(&warehouse).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: warehouse with name '%s' not found", lineNumber, warehouseName))
					continue
				}
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to fetch warehouse by name",
					"message": err.Error(),
				})
				return
			}
		}

		sn := csvValue(record, headerMap, "sn")

		name := csvValue(record, headerMap, "name")
		if name == "" {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: name is required", lineNumber))
			continue
		}

		category := csvValue(record, headerMap, "category")
		description := csvValue(record, headerMap, "description")
		unit := csvValue(record, headerMap, "unit")
		if unit == "" {
			unit = "pcs"
		}

		quantity, quantityProvided, err := parseOptionalFloat(csvValue(record, headerMap, "quantity"))
		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: invalid quantity value", lineNumber))
			continue
		}

		minStock, minStockProvided, err := parseOptionalFloat(csvValue(record, headerMap, "min_stock"))
		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: invalid min_stock value", lineNumber))
			continue
		}

		maxStock, maxStockProvided, err := parseOptionalFloat(csvValue(record, headerMap, "max_stock"))
		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: invalid max_stock value", lineNumber))
			continue
		}

		unitPrice, unitPriceProvided, err := parseOptionalFloat(csvValue(record, headerMap, "unit_price"))
		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: invalid unit_price value", lineNumber))
			continue
		}

		statusProvided := false
		isActive := true
		if _, ok := headerMap["status"]; ok {
			statusValue := strings.ToLower(csvValue(record, headerMap, "status"))
			if statusValue != "" {
				statusProvided = true
				switch statusValue {
				case "active", "aktif", "1", "true", "ya", "yes":
					isActive = true
				case "inactive", "nonaktif", "tidak aktif", "0", "false", "tidak", "no":
					isActive = false
				default:
					summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: invalid status value '%s' (use Active/Inactive)", lineNumber, statusValue))
					continue
				}
			}
		}

		var existing models.InventoryItem
		itemQuery := tx.Where("warehouse_id = ?", warehouse.ID)
		if sn != "" {
			itemQuery = itemQuery.Where("sku = ?", sn)
		} else {
			itemQuery = itemQuery.Where("name = ?", name)
		}

		err = itemQuery.First(&existing).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			newItem := models.InventoryItem{
				WarehouseID: warehouse.ID,
				SN:          sn,
				Name:        name,
				Description: description,
				Category:    category,
				Unit:        unit,
				Quantity:    quantity,
				MinStock:    minStock,
				MaxStock:    maxStock,
				UnitPrice:   unitPrice,
				IsActive:    isActive,
			}

			if !quantityProvided {
				newItem.Quantity = 0
			}
			if !minStockProvided {
				newItem.MinStock = 0
			}
			if !maxStockProvided {
				newItem.MaxStock = 0
			}
			if !unitPriceProvided {
				newItem.UnitPrice = 0
			}
			if !statusProvided {
				newItem.IsActive = quantity > 0
			}

			if err := tx.Create(&newItem).Error; err != nil {
				summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: failed to create item (%v)", lineNumber, err))
				continue
			}
			summary.Inserted++

		case err != nil:
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check existing inventory item",
				"message": err.Error(),
			})
			return

		default:
			if category != "" {
				existing.Category = category
			}
			if description != "" {
				existing.Description = description
			}
			if unit != "" {
				existing.Unit = unit
			}

			existing.Name = name

			if quantityProvided {
				existing.Quantity = quantity
			}
			if minStockProvided {
				existing.MinStock = minStock
			}
			if maxStockProvided {
				existing.MaxStock = maxStock
			}
			if unitPriceProvided {
				existing.UnitPrice = unitPrice
			}
			if statusProvided {
				existing.IsActive = isActive
			}

			if err := tx.Save(&existing).Error; err != nil {
				summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: failed to update item (%v)", lineNumber, err))
				continue
			}
			summary.Updated++
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to commit import transaction",
			"message": err.Error(),
		})
		return
	}

	status := http.StatusOK
	message := fmt.Sprintf("Import completed. Inserted: %d, Updated: %d", summary.Inserted, summary.Updated)
	if len(summary.Errors) > 0 {
		status = http.StatusMultiStatus
	}

	c.JSON(status, gin.H{
		"message": message,
		"summary": summary,
	})
}

// ExportItemsToCSV streams inventory data as CSV
func (h *InventoryHandler) ExportItemsToCSV(c *gin.Context) {
	var items []models.InventoryItem

	query := applyInventoryFilters(h.db.Preload("Warehouse"), c)
	if err := query.Order("name ASC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch inventory items for export",
			"message": err.Error(),
		})
		return
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	headers := []string{
		"SN",
		"Item Name",
		"Category",
		"Warehouse",
		"Quantity",
		"Unit",
		"Min Stock",
		"Status",
	}

	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to write CSV header",
			"message": err.Error(),
		})
		return
	}

	for _, item := range items {
		status := "Inactive"
		if item.IsActive {
			status = "Active"
		}
		record := []string{
			item.SN,
			item.Name,
			item.Category,
			item.Warehouse.Name,
			formatFloat(item.Quantity),
			item.Unit,
			formatFloat(item.MinStock),
			status,
		}

		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to write CSV record",
				"message": err.Error(),
			})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to finalize CSV export",
			"message": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("inventory-%s.csv", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}

// ExportItemsToPDF streams inventory data as PDF document
func (h *InventoryHandler) ExportItemsToPDF(c *gin.Context) {
	var items []models.InventoryItem

	query := applyInventoryFilters(h.db.Preload("Warehouse"), c)
	if err := query.Order("name ASC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch inventory items for export",
			"message": err.Error(),
		})
		return
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "L",
		UnitStr:        "mm",
		SizeStr:        "A4",
	})
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Inventory Report")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Generated at: %s", time.Now().Format("02 Jan 2006 15:04")))
	pdf.Ln(10)

	headers := []string{"No", "SN", "Name", "Warehouse", "Category", "Qty", "Unit", "Unit Price", "Min"}
	widths := []float64{10, 28, 58, 60, 36, 18, 20, 26, 18}
	lineHeight := 6.0

	renderHeader := func() {
		pdf.SetFillColor(240, 240, 240)
		pdf.SetFont("Arial", "B", 9)
		for idx, header := range headers {
			pdf.CellFormat(widths[idx], 8, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont("Arial", "", 9)
	}

	renderHeader()

	_, _, _, bottomMargin := pdf.GetMargins()
	_, pageHeight := pdf.GetPageSize()
	usableHeight := pageHeight - bottomMargin

	for idx, item := range items {
		row := []string{
			strconv.Itoa(idx + 1),
			item.SN,
			item.Name,
			fmt.Sprintf("%s - %s", item.Warehouse.Code, item.Warehouse.Name),
			item.Category,
			formatFloat(item.Quantity),
			item.Unit,
			formatFloat(item.UnitPrice),
			formatFloat(item.MinStock),
		}

		rowHeight := lineHeight
		cellLines := make([][][]byte, len(row))
		for idx, value := range row {
			text := strings.TrimSpace(value)
			if text == "" {
				text = "-"
			}
			wrapWidth := widths[idx] - 2
			if wrapWidth <= 0 {
				wrapWidth = widths[idx]
			}
			lines := pdf.SplitLines([]byte(text), wrapWidth)
			if len(lines) == 0 {
				lines = [][]byte{[]byte(" ")}
			}
			cellLines[idx] = lines
			cellHeight := float64(len(lines)) * lineHeight
			if cellHeight > rowHeight {
				rowHeight = cellHeight
			}
		}

		if pdf.GetY()+rowHeight > usableHeight {
			pdf.AddPage()
			renderHeader()
			_, pageHeight = pdf.GetPageSize()
			usableHeight = pageHeight - bottomMargin
		}

		xLeft := pdf.GetX()
		yTop := pdf.GetY()

		for idx, lines := range cellLines {
			width := widths[idx]
			align := "L"
			if idx == 4 || idx == 6 || idx == 7 {
				align = "R"
			}

			cellX := pdf.GetX()
			cellY := pdf.GetY()
			pdf.Rect(cellX, cellY, width, rowHeight, "")

			var builder strings.Builder
			for i, line := range lines {
				builder.Write(line)
				if i < len(lines)-1 {
					builder.WriteByte('\n')
				}
			}
			content := builder.String()

			pdf.SetXY(cellX, cellY)
			pdf.MultiCell(width, lineHeight, content, "", align, false)
			pdf.SetXY(cellX+width, cellY)
		}

		pdf.SetXY(xLeft, yTop+rowHeight)
	}

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate PDF file",
			"message": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("inventory-%s.pdf", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}

// DownloadImportTemplate returns a ready-to-use CSV template for inventory imports
func (h *InventoryHandler) DownloadImportTemplate(c *gin.Context) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	headers := []string{
		"SN",
		"Item Name",
		"Category",
		"Warehouse Name",
		"Quantity",
		"Unit",
		"Min Stock",
		"Status",
	}

	sample := []string{
		"SN-0001",
		"Laptop Dell XPS 13",
		"Laptop",
		"Main Warehouse",
		"10",
		"unit",
		"2",
		"Active",
	}

	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to write template header",
			"message": err.Error(),
		})
		return
	}

	if err := writer.Write(sample); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to write template sample row",
			"message": err.Error(),
		})
		return
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to prepare template file",
			"message": err.Error(),
		})
		return
	}

	filename := "inventory-import-template.csv"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}

func applyInventoryFilters(query *gorm.DB, c *gin.Context) *gorm.DB {
	// Filter by warehouse
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	// Filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Search by name or SN
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Low stock filter
	if c.Query("low_stock") == "true" {
		query = query.Where("quantity <= min_stock")
	}

	return query
}

func applyTransactionFilters(query *gorm.DB, c *gin.Context) *gorm.DB {
	if transType := c.Query("type"); transType != "" {
		query = query.Where("type = ?", transType)
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("(from_warehouse_id = ? OR to_warehouse_id = ?)", warehouseID, warehouseID)
	}

	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		like := "%" + search + "%"
		query = query.Joins("LEFT JOIN inventory_items AS export_items ON export_items.id = inventory_transactions.item_id").
			Where("export_items.name ILIKE ? OR export_items.sku ILIKE ? OR inventory_transactions.reference ILIKE ? OR COALESCE(inventory_transactions.notes, '') ILIKE ?", like, like, like, like).
			Distinct()
	}

	return query
}

func csvValue(record []string, header map[string]int, key string) string {
	idx, ok := header[key]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

func parseOptionalFloat(value string) (float64, bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false, err
	}

	return parsed, true, nil
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func isCSVRecordEmpty(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

// ExportTransactionsToCSV streams inventory transactions as CSV
func (h *InventoryHandler) ExportTransactionsToCSV(c *gin.Context) {
	var transactions []models.InventoryTransaction

	query := applyTransactionFilters(h.db.Model(&models.InventoryTransaction{}), c)

	if err := query.
		Preload("Item").
		Preload("Item.Warehouse").
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch transactions for export",
			"message": err.Error(),
		})
		return
	}

	if roleName, ok := c.Get("role_name"); ok {
		if name, isString := roleName.(string); isString && strings.EqualFold(name, "employee") {
			userID, ok := h.contextUserID(c)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unable to resolve user context"})
				return
			}

			allowedIDs, err := h.getUserWarehouseIDs(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve warehouse permissions"})
				return
			}

			if len(allowedIDs) == 0 {
				transactions = []models.InventoryTransaction{}
			} else {
				allowedSet := buildUintSet(allowedIDs)
				transactions = filterTransactionsByWarehouses(transactions, allowedSet)
			}
		}
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	headers := []string{
		"No",
		"Date",
		"Type",
		"Item SN",
		"Item Name",
		"Quantity",
		"Unit",
		"From Warehouse",
		"To Warehouse",
		"Reference",
		"Notes",
		"Created By",
	}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to write CSV header",
			"message": err.Error(),
		})
		return
	}

	for idx, t := range transactions {
		itemSN := "-"
		itemName := "-"
		unit := "-"
		if t.Item.ID != 0 {
			if strings.TrimSpace(t.Item.SN) != "" {
				itemSN = t.Item.SN
			}
			if strings.TrimSpace(t.Item.Name) != "" {
				itemName = t.Item.Name
			}
			if strings.TrimSpace(t.Item.Unit) != "" {
				unit = t.Item.Unit
			}
		}

		fromWarehouse := "-"
		if t.FromWarehouse != nil && strings.TrimSpace(t.FromWarehouse.Name) != "" {
			fromWarehouse = t.FromWarehouse.Name
		}

		toWarehouse := "-"
		if t.ToWarehouse != nil && strings.TrimSpace(t.ToWarehouse.Name) != "" {
			toWarehouse = t.ToWarehouse.Name
		}

		reference := "-"
		if strings.TrimSpace(t.Reference) != "" {
			reference = t.Reference
		}

		notes := "-"
		if strings.TrimSpace(t.Notes) != "" {
			notes = strings.ReplaceAll(t.Notes, "\n", " ")
		}

		createdBy := "-"
		if t.CreatedBy.ID != 0 && strings.TrimSpace(t.CreatedBy.FullName) != "" {
			createdBy = t.CreatedBy.FullName
		}

		record := []string{
			strconv.Itoa(idx + 1),
			t.CreatedAt.Format("2006-01-02 15:04:05"),
			strings.ToUpper(t.Type),
			itemSN,
			itemName,
			fmt.Sprintf("%.2f", t.Quantity),
			unit,
			fromWarehouse,
			toWarehouse,
			reference,
			notes,
			createdBy,
		}

		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to write CSV record",
				"message": err.Error(),
			})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to finalize CSV",
			"message": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("inventory-transactions-%s.csv", time.Now().Format("20060102-150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}

// ExportTransactionsToPDF streams inventory transactions as PDF
func (h *InventoryHandler) ExportTransactionsToPDF(c *gin.Context) {
	var transactions []models.InventoryTransaction

	query := applyTransactionFilters(h.db.Model(&models.InventoryTransaction{}), c)

	if err := query.
		Preload("Item").
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch transactions for export",
			"message": err.Error(),
		})
		return
	}

	if roleName, ok := c.Get("role_name"); ok {
		if name, isString := roleName.(string); isString && strings.EqualFold(name, "employee") {
			userID, ok := h.contextUserID(c)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unable to resolve user context"})
				return
			}

			allowedIDs, err := h.getUserWarehouseIDs(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve warehouse permissions"})
				return
			}

			if len(allowedIDs) == 0 {
				transactions = []models.InventoryTransaction{}
			} else {
				allowedSet := buildUintSet(allowedIDs)
				transactions = filterTransactionsByWarehouses(transactions, allowedSet)
			}
		}
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetTitle("Inventory Transactions", false)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Inventory Transactions")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)

	headers := []string{
		"No",
		"Date",
		"Type",
		"Item",
		"Quantity",
		"From",
		"To",
		"Reference",
		"Notes",
		"Created By",
	}
	colWidths := []float64{10, 35, 18, 75, 24, 30, 30, 30, 55, 30}

	lineHeight := 6.0

	renderHeader := func() {
		pdf.SetFillColor(230, 230, 230)
		pdf.SetFont("Arial", "B", 10)
		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont("Arial", "", 9)
	}

	renderHeader()

	_, _, _, bottomMargin := pdf.GetMargins()
	_, pageHeight := pdf.GetPageSize()
	usableHeight := pageHeight - bottomMargin

	for idx, t := range transactions {
		itemName := "-"
		itemSN := ""
		unit := "-"
		if t.Item.ID != 0 {
			if strings.TrimSpace(t.Item.Name) != "" {
				itemName = t.Item.Name
			}
			if strings.TrimSpace(t.Item.SN) != "" {
				itemSN = t.Item.SN
			}
			if strings.TrimSpace(t.Item.Unit) != "" {
				unit = t.Item.Unit
			}
		}

		itemDisplay := itemName
		if itemSN != "" {
			itemDisplay = fmt.Sprintf("%s (%s)", itemDisplay, itemSN)
		}

		fromWarehouse := "-"
		if t.FromWarehouse != nil && strings.TrimSpace(t.FromWarehouse.Name) != "" {
			fromWarehouse = t.FromWarehouse.Name
		}

		toWarehouse := "-"
		if t.ToWarehouse != nil && strings.TrimSpace(t.ToWarehouse.Name) != "" {
			toWarehouse = t.ToWarehouse.Name
		}

		reference := "-"
		if strings.TrimSpace(t.Reference) != "" {
			reference = t.Reference
		}

		notes := "-"
		if strings.TrimSpace(t.Notes) != "" {
			notes = strings.ReplaceAll(strings.TrimSpace(t.Notes), "\r\n", "\n")
		}

		createdBy := "-"
		if t.CreatedBy.ID != 0 && strings.TrimSpace(t.CreatedBy.FullName) != "" {
			createdBy = t.CreatedBy.FullName
		}

		row := []string{
			strconv.Itoa(idx + 1),
			t.CreatedAt.Format("2006-01-02 15:04"),
			strings.ToUpper(t.Type),
			itemDisplay,
			fmt.Sprintf("%.2f %s", t.Quantity, unit),
			fromWarehouse,
			toWarehouse,
			reference,
			notes,
			createdBy,
		}

		rowHeight := lineHeight
		cellLines := make([][][]byte, len(row))
		for i, value := range row {
			text := strings.TrimSpace(value)
			if text == "" {
				text = "-"
			}
			wrapWidth := colWidths[i] - 2
			if wrapWidth <= 0 {
				wrapWidth = colWidths[i]
			}
			lines := pdf.SplitLines([]byte(text), wrapWidth)
			if len(lines) == 0 {
				lines = [][]byte{[]byte(" ")}
			}
			cellLines[i] = lines
			height := float64(len(lines)) * lineHeight
			if height > rowHeight {
				rowHeight = height
			}
		}

		if pdf.GetY()+rowHeight > usableHeight {
			pdf.AddPage()
			renderHeader()
			_, pageHeight = pdf.GetPageSize()
			usableHeight = pageHeight - bottomMargin
		}

		xLeft := pdf.GetX()
		yTop := pdf.GetY()

		for i, lines := range cellLines {
			width := colWidths[i]
			align := "L"
			if i == 0 {
				align = "C"
			}
			if i == 4 {
				align = "R"
			}

			cellX := pdf.GetX()
			cellY := pdf.GetY()
			pdf.Rect(cellX, cellY, width, rowHeight, "")

			for lineIdx, line := range lines {
				pdf.SetXY(cellX, cellY+float64(lineIdx)*lineHeight)
				pdf.CellFormat(width, lineHeight, string(line), "", 0, align, false, 0, "")
			}
			pdf.SetXY(cellX+width, cellY)
		}

		pdf.SetXY(xLeft, yTop+rowHeight)
	}

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate PDF",
			"message": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("inventory-transactions-%s.pdf", time.Now().Format("20060102-150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}
