package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"tatapps/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetAllCategories retrieves all categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	var categories []models.Category

	if err := h.db.Order("name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(categories) > 0 {
		type countResult struct {
			Category string
			Count    int64
		}

		var counts []countResult
		if err := h.db.
			Model(&models.InventoryItem{}).
			Select("LOWER(category) AS category, COUNT(*) AS count").
			Group("LOWER(category)").
			Find(&counts).Error; err == nil {
			countMap := make(map[string]int64, len(counts))
			for _, entry := range counts {
				countMap[entry.Category] = entry.Count
			}
			for idx, cat := range categories {
				key := strings.ToLower(cat.Name)
				if val, ok := countMap[key]; ok {
					categories[idx].ItemCount = int(val)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetCategoryByID retrieves a single category by ID
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	result := h.db.First(&category, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.ItemCount = int(h.countItemsForCategory(category.Name))

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color" binding:"required"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		Color:       input.Color,
		IsActive:    input.IsActive,
	}

	if result := h.db.Create(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// UpdateCategory updates a category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	if result := h.db.First(&category, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Color       string `json:"color"`
		IsActive    *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		category.Name = input.Name
	}
	if input.Code != "" {
		category.Code = input.Code
	}
	if input.Description != "" {
		category.Description = input.Description
	}
	if input.Color != "" {
		category.Color = input.Color
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	if result := h.db.Save(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	if result := h.db.First(&category, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category is being used by any items
	if h.countItemsForCategory(category.Name) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete category. It is being used by inventory items.",
		})
		return
	}

	if result := h.db.Delete(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *CategoryHandler) countItemsForCategory(name string) int64 {
	if name == "" {
		return 0
	}

	var count int64
	h.db.Model(&models.InventoryItem{}).
		Where("LOWER(category) = ?", strings.ToLower(name)).
		Count(&count)
	return count
}
