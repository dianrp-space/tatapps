package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequirePermission ensures the authenticated user owns all requested permissions.
func RequirePermission(db *gorm.DB, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requirePermissions(c, db, permissions, false)
	}
}

// RequireAnyPermission ensures the authenticated user owns at least one permission from the list.
func RequireAnyPermission(db *gorm.DB, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requirePermissions(c, db, permissions, true)
	}
}

func requirePermissions(c *gin.Context, db *gorm.DB, permissions []string, allowAny bool) {
	if len(permissions) == 0 {
		c.Next()
		return
	}

	roleName := ""
	if value, ok := c.Get("role_name"); ok {
		if name, isString := value.(string); isString {
			roleName = name
			if strings.EqualFold(name, "admin") {
				c.Next()
				return
			}
		}
	}

	roleID, ok := extractRoleID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
		c.Abort()
		return
	}

	var count int64
	if err := db.
		Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Where("permissions.name IN ?", permissions).
		Distinct("permissions.name").
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions"})
		c.Abort()
		return
	}

	if allowAny {
		if count == 0 && !strings.EqualFold(roleName, "employee") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
	} else {
		if count != int64(len(permissions)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
	}

	c.Next()
}

func extractRoleID(c *gin.Context) (uint, bool) {
	roleIDValue, exists := c.Get("role_id")
	if !exists {
		return 0, false
	}

	switch v := roleIDValue.(type) {
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
	default:
		return 0, false
	}
}
