package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"tatapps/internal/config"
	"tatapps/internal/models"
	"tatapps/internal/services/notification"
	"tatapps/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db    *gorm.DB
	notif *notification.NotificationService
	cfg   *config.Config
}

func NewUserHandler(db *gorm.DB, notif *notification.NotificationService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		db:    db,
		notif: notif,
		cfg:   cfg,
	}
}

type RoleRequest struct {
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	Color          string   `json:"color"`
	MenuKeys       []string `json:"menu_keys"`
	PermissionKeys []string `json:"permission_keys"`
}

type MenuDefinition struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Parent   string `json:"parent,omitempty"`
	Category string `json:"category"`
}

var menuDefinitions = []MenuDefinition{
	{Key: "dashboard", Label: "Dashboard", Category: "main"},
	{Key: "warehouses", Label: "Warehouses", Category: "operations"},
	{Key: "inventory", Label: "Inventory", Category: "operations"},
	{Key: "inventory.items", Label: "Inventory • Items", Parent: "inventory", Category: "operations"},
	{Key: "inventory.categories", Label: "Inventory • Categories", Parent: "inventory", Category: "operations"},
	{Key: "inventory.transactions", Label: "Inventory • Transactions", Parent: "inventory", Category: "operations"},
	{Key: "employees", Label: "Employees", Category: "operations"},
	{Key: "employees.data", Label: "Employees • Data Karyawan", Parent: "employees", Category: "operations"},
	{Key: "employees.divisions", Label: "Employees • Divisi", Parent: "employees", Category: "operations"},
	{Key: "employees.positions", Label: "Employees • Jabatan", Parent: "employees", Category: "operations"},
	{Key: "leads", Label: "Leads", Category: "operations"},
	{Key: "projects", Label: "Projects", Category: "operations"},
	{Key: "purchase_orders", Label: "Purchase Orders", Category: "operations"},
	{Key: "support", Label: "Support", Category: "support"},
	{Key: "settings", Label: "Settings", Category: "settings"},
	{Key: "settings.profile", Label: "Settings • Profile", Parent: "settings", Category: "settings"},
	{Key: "settings.company", Label: "Settings • Profil Perusahaan", Parent: "settings", Category: "settings"},
	{Key: "settings.notifications", Label: "Settings • Notifications", Parent: "settings", Category: "settings"},
	{Key: "settings.users", Label: "Settings • Users", Parent: "settings", Category: "settings"},
	{Key: "settings.sites", Label: "Settings • Sites", Parent: "settings", Category: "settings"},
}

var menuDefinitionLookup map[string]MenuDefinition
var permissionDefinitionLookup map[string]PermissionDefinition

func init() {
	menuDefinitionLookup = make(map[string]MenuDefinition, len(menuDefinitions))
	for _, def := range menuDefinitions {
		menuDefinitionLookup[def.Key] = def
	}
	permissionDefinitionLookup = make(map[string]PermissionDefinition, len(permissionDefinitions))
	for _, def := range permissionDefinitions {
		permissionDefinitionLookup[def.Name] = def
	}
}

type PermissionDefinition struct {
	Name        string
	Description string
	Module      string
	Action      string
}

var permissionDefinitions = []PermissionDefinition{
	// Inventory
	{Name: "inventory.view", Description: "View inventory items", Module: "inventory", Action: "view"},
	{Name: "inventory.create", Description: "Create inventory items", Module: "inventory", Action: "create"},
	{Name: "inventory.update", Description: "Update inventory items", Module: "inventory", Action: "update"},
	{Name: "inventory.delete", Description: "Delete inventory items", Module: "inventory", Action: "delete"},

	// Category
	{Name: "category.view", Description: "View categories", Module: "category", Action: "view"},
	{Name: "category.create", Description: "Create categories", Module: "category", Action: "create"},
	{Name: "category.update", Description: "Update categories", Module: "category", Action: "update"},
	{Name: "category.delete", Description: "Delete categories", Module: "category", Action: "delete"},

	// Warehouse
	{Name: "warehouse.view", Description: "View warehouses", Module: "warehouse", Action: "view"},
	{Name: "warehouse.create", Description: "Create warehouse", Module: "warehouse", Action: "create"},
	{Name: "warehouse.update", Description: "Update warehouse", Module: "warehouse", Action: "update"},
	{Name: "warehouse.delete", Description: "Delete warehouse", Module: "warehouse", Action: "delete"},

	// Employee
	{Name: "employee.view", Description: "View employees", Module: "employee", Action: "view"},
	{Name: "employee.create", Description: "Create employee", Module: "employee", Action: "create"},
	{Name: "employee.update", Description: "Update employee", Module: "employee", Action: "update"},
	{Name: "employee.delete", Description: "Delete employee", Module: "employee", Action: "delete"},

	// Lead
	{Name: "lead.view", Description: "View leads", Module: "lead", Action: "view"},
	{Name: "lead.create", Description: "Create lead", Module: "lead", Action: "create"},
	{Name: "lead.update", Description: "Update lead", Module: "lead", Action: "update"},
	{Name: "lead.delete", Description: "Delete lead", Module: "lead", Action: "delete"},

	// Project
	{Name: "project.view", Description: "View projects", Module: "project", Action: "view"},
	{Name: "project.create", Description: "Create project", Module: "project", Action: "create"},
	{Name: "project.update", Description: "Update project", Module: "project", Action: "update"},
	{Name: "project.delete", Description: "Delete project", Module: "project", Action: "delete"},

	// Purchase Order
	{Name: "po.view", Description: "View purchase orders", Module: "po", Action: "view"},
	{Name: "po.create", Description: "Create purchase order", Module: "po", Action: "create"},
	{Name: "po.update", Description: "Update purchase order", Module: "po", Action: "update"},
	{Name: "po.approve", Description: "Approve purchase order", Module: "po", Action: "approve"},
	{Name: "po.delete", Description: "Delete purchase order", Module: "po", Action: "delete"},
}

// GetAll returns all users with their roles
func (h *UserHandler) GetAll(c *gin.Context) {
	var users []models.User

	// Get query parameters for filtering
	roleID := c.Query("role_id")

	query := h.db.Preload("Role").Preload("Role.Menus").Preload("Role.Permissions").
		Preload("Warehouses").Preload("Warehouses.Warehouse")

	// Filter by role if provided
	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Remove password from response
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetByID returns a single user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.Preload("Role").Preload("Role.Menus").Preload("Role.Permissions").
		Preload("Warehouses").Preload("Warehouses.Warehouse").
		First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		FullName     string `json:"full_name" binding:"required"`
		Email        string `json:"email" binding:"required,email"`
		Phone        string `json:"phone"`
		Password     string `json:"password" binding:"required,min=6"`
		RoleID       uint   `json:"role_id" binding:"required"`
		SendWelcome  bool   `json:"send_welcome"`
		WarehouseIDs []uint `json:"warehouse_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists (including soft-deleted records)
	var existingUser models.User
	if err := h.db.Unscoped().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Verify role exists
	var role models.Role
	if err := h.db.Preload("Menus").First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		RoleID:   req.RoleID,
		IsActive: true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		// Check for duplicate email error (including soft-deleted records causing constraint violation)
		errStr := err.Error()
		if errStr != "" && (errStr == "UNIQUE constraint failed: users.email" ||
			errStr == "Error 1062: Duplicate entry" ||
			errStr == "pq: duplicate key value violates unique constraint \"users_email_key\"" ||
			errStr == "ERROR: duplicate key value violates unique constraint \"idx_users_email\" (SQLSTATE 23505)") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	if err := h.syncUserWarehouses(user.ID, req.WarehouseIDs); err != nil {
		h.db.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign warehouse access"})
		return
	}

	// Load relationships
	h.db.Preload("Role").Preload("Role.Menus").Preload("Role.Permissions").
		Preload("Warehouses").Preload("Warehouses.Warehouse").
		First(&user, user.ID)
	user.Password = ""

	// Send welcome message via WhatsApp if requested and phone is provided
	if req.SendWelcome && req.Phone != "" {
		appURL := h.cfg.FrontendURL
		welcomeMessage := fmt.Sprintf(
			"🎉 *Welcome to TatApps!*\n\n"+
				"Your account has been created successfully.\n\n"+
				"*Account Details:*\n"+
				"📧 Email: %s\n"+
				"🔑 Temporary Password: %s\n"+
				"👤 Role: %s\n\n"+
				"Akses aplikasi: %s\n\n"+
				"⚠️ Please change your password after first login for security.",
			user.Email,
			req.Password,
			role.Name,
			appURL,
		)

		if err := h.notif.WhatsApp.SendMessage(notification.WhatsAppMessage{
			Phone:   req.Phone,
			Message: welcomeMessage,
		}); err != nil {
			// Log error but don't fail the user creation
			fmt.Printf("Failed to send welcome message: %v\n", err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser updates core user information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")

	var req struct {
		FullName     string  `json:"full_name" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		Phone        string  `json:"phone"`
		RoleID       uint    `json:"role_id" binding:"required"`
		Password     *string `json:"password"`
		WarehouseIDs []uint  `json:"warehouse_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, idParam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Ensure email unique for other users
	var existing models.User
	if err := h.db.Unscoped().Where("email = ? AND id != ?", req.Email, user.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Validate role
	var role models.Role
	if err := h.db.Preload("Menus").First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	user.FullName = strings.TrimSpace(req.FullName)
	user.Email = strings.TrimSpace(req.Email)
	user.Phone = strings.TrimSpace(req.Phone)
	user.RoleID = req.RoleID

	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		if len(*req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
			return
		}

		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if err := h.syncUserWarehouses(user.ID, req.WarehouseIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update warehouse access"})
		return
	}

	h.db.Preload("Role").Preload("Role.Menus").Preload("Role.Permissions").
		Preload("Warehouses").Preload("Warehouses.Warehouse").
		First(&user, user.ID)
	user.Password = ""

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}

// UpdateUserStatus updates user active status
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	user.IsActive = req.IsActive
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Get current user ID from context (set by auth middleware)
	currentUserID, exists := c.Get("user_id")
	if exists && fmt.Sprintf("%v", currentUserID) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot delete your own account"})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) syncUserWarehouses(userID uint, warehouseIDs []uint) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("user_id = ?", userID).Delete(&models.UserWarehouse{}).Error; err != nil {
			return err
		}

		if len(warehouseIDs) == 0 {
			return nil
		}

		unique := make(map[uint]struct{})
		for _, id := range warehouseIDs {
			if id == 0 {
				continue
			}
			unique[id] = struct{}{}
		}

		if len(unique) == 0 {
			return nil
		}

		records := make([]models.UserWarehouse, 0, len(unique))
		for id := range unique {
			records = append(records, models.UserWarehouse{
				UserID:      userID,
				WarehouseID: id,
			})
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}

		return nil
	})
}

const defaultRoleColor = "#2563EB"

var roleColorPattern = regexp.MustCompile(`^(?i)#?([0-9a-f]{3}|[0-9a-f]{6})$`)

func sanitizeRoleColor(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return defaultRoleColor, nil
	}

	if !strings.HasPrefix(trimmed, "#") {
		trimmed = "#" + trimmed
	}

	if !roleColorPattern.MatchString(trimmed) {
		return "", fmt.Errorf("invalid color format: %s", value)
	}

	hex := strings.TrimPrefix(trimmed, "#")
	if len(hex) == 3 {
		expanded := make([]byte, 0, 6)
		for i := 0; i < 3; i++ {
			expanded = append(expanded, hex[i], hex[i])
		}
		hex = string(expanded)
	}

	return "#" + strings.ToUpper(hex), nil
}

func sanitizeMenuKeys(keys []string) ([]string, error) {
	seen := make(map[string]struct{})
	cleaned := make([]string, 0, len(keys))

	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}
		if _, exists := menuDefinitionLookup[trimmed]; !exists {
			return nil, fmt.Errorf("invalid menu key: %s", trimmed)
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			cleaned = append(cleaned, trimmed)
		}
	}

	if len(cleaned) == 0 {
		return nil, errors.New("menu_keys cannot be empty")
	}

	return cleaned, nil
}

func (h *UserHandler) assignRoleMenus(tx *gorm.DB, roleID uint, menuKeys []string) error {
	if err := tx.Unscoped().Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
		return err
	}

	if len(menuKeys) == 0 {
		return nil
	}

	menus := make([]models.RoleMenu, 0, len(menuKeys))
	for _, key := range menuKeys {
		menus = append(menus, models.RoleMenu{RoleID: roleID, MenuKey: key})
	}

	return tx.Create(&menus).Error
}

func sanitizePermissionKeys(keys []string) ([]string, error) {
	cleaned := make([]string, 0, len(keys))
	seen := make(map[string]struct{})

	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}
		if _, exists := permissionDefinitionLookup[trimmed]; !exists {
			return nil, fmt.Errorf("invalid permission key: %s", trimmed)
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		cleaned = append(cleaned, trimmed)
	}

	return cleaned, nil
}

func (h *UserHandler) assignRolePermissions(tx *gorm.DB, role *models.Role, permissionKeys []string) error {
	if err := tx.Model(role).Association("Permissions").Clear(); err != nil {
		return err
	}

	if len(permissionKeys) == 0 {
		return nil
	}

	for _, key := range permissionKeys {
		def, exists := permissionDefinitionLookup[key]
		if !exists {
			return fmt.Errorf("invalid permission key: %s", key)
		}

		var permission models.Permission
		if err := tx.
			Where("name = ?", key).
			Attrs(models.Permission{
				Description: def.Description,
				Module:      def.Module,
				Action:      def.Action,
			}).
			FirstOrCreate(&permission).Error; err != nil {
			return err
		}
	}

	var permissions []models.Permission
	if err := tx.Where("name IN ?", permissionKeys).Find(&permissions).Error; err != nil {
		return err
	}

	if len(permissions) != len(permissionKeys) {
		validMap := make(map[string]struct{}, len(permissions))
		for _, p := range permissions {
			validMap[p.Name] = struct{}{}
		}
		for _, key := range permissionKeys {
			if _, ok := validMap[key]; !ok {
				return fmt.Errorf("invalid permission key: %s", key)
			}
		}
	}

	if err := tx.Model(role).Association("Permissions").Replace(&permissions); err != nil {
		return err
	}

	return nil
}

// GetRoles returns all available roles
func (h *UserHandler) GetRoles(c *gin.Context) {
	var roles []models.Role

	if err := h.db.Preload("Menus").Preload("Permissions").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *UserHandler) CreateRole(c *gin.Context) {
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	color, err := sanitizeRoleColor(req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuKeys, err := sanitizeMenuKeys(req.MenuKeys)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionKeys, err := sanitizePermissionKeys(req.PermissionKeys)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.Role{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Color:       color,
	}

	tx := h.db.Begin()
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	if err := h.assignRoleMenus(tx, role.ID, menuKeys); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save role menus"})
		return
	}

	if err := h.assignRolePermissions(tx, &role, permissionKeys); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save role permissions"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	if err := h.db.Preload("Menus").Preload("Permissions").First(&role, role.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"data":    role,
	})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")
	roleID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var role models.Role
	if err := h.db.Preload("Menus").First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role"})
		return
	}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	color, err := sanitizeRoleColor(req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuKeys, err := sanitizeMenuKeys(req.MenuKeys)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionKeys, err := sanitizePermissionKeys(req.PermissionKeys)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role name is required"})
		return
	}

	role.Name = name
	role.Description = strings.TrimSpace(req.Description)
	role.Color = color

	tx := h.db.Begin()
	if err := tx.Model(&role).Updates(map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
		"color":       role.Color,
	}).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	if err := h.assignRoleMenus(tx, role.ID, menuKeys); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role menus"})
		return
	}

	if err := h.assignRolePermissions(tx, &role, permissionKeys); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role permissions"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	if err := h.db.Preload("Menus").Preload("Permissions").First(&role, role.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"data":    role,
	})
}

func (h *UserHandler) DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	roleID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var role models.Role
	if err := h.db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role"})
		return
	}

	if strings.EqualFold(role.Name, "admin") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Default admin role cannot be deleted"})
		return
	}

	var userCount int64
	if err := h.db.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role usage"})
		return
	}

	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is still assigned to users"})
		return
	}

	tx := h.db.Begin()
	if err := tx.Where("role_id = ?", role.ID).Delete(&models.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role menus"})
		return
	}

	if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role permissions"})
		return
	}

	if err := tx.Delete(&role).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

func (h *UserHandler) GetRoleMenuOptions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": menuDefinitions})
}

func (h *UserHandler) GetRolePermissionOptions(c *gin.Context) {
	for _, def := range permissionDefinitions {
		var perm models.Permission
		err := h.db.Unscoped().Where("name = ?", def.Name).First(&perm).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			perm = models.Permission{
				Name:        def.Name,
				Description: def.Description,
				Module:      def.Module,
				Action:      def.Action,
			}
			if err := h.db.Create(&perm).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure permission definitions"})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure permission definitions"})
			return
		} else {
			updates := make(map[string]interface{})
			if perm.Description != def.Description {
				updates["description"] = def.Description
			}
			if perm.Module != def.Module {
				updates["module"] = def.Module
			}
			if perm.Action != def.Action {
				updates["action"] = def.Action
			}
			if perm.DeletedAt.Valid {
				updates["deleted_at"] = nil
			}

			if len(updates) > 0 {
				if err := h.db.Unscoped().Model(&perm).Updates(updates).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure permission definitions"})
					return
				}
			}
		}
	}

	var permissions []models.Permission
	if err := h.db.Order("module ASC, action ASC").Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}
