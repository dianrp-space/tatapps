package database

import (
	"log"
	"tatapps/internal/models"
	"tatapps/internal/utils"
	"time"

	"gorm.io/gorm"
)

// SeedData seeds initial data for development
func SeedData(db *gorm.DB) error {
	// Check if already seeded
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database...")

	// Create Roles
	roles := []models.Role{
		{Name: "admin", Description: "Administrator with full access"},
		{Name: "manager", Description: "Manager with limited access"},
		{Name: "employee", Description: "Employee with basic access"},
	}

	for i := range roles {
		if err := db.Create(&roles[i]).Error; err != nil {
			return err
		}
	}

	// Create default site setting
	siteSetting := models.SiteSetting{
		AppName:      "TatApps",
		SMTPPort:     587,
		SMTPFromName: "TatApps",
	}
	if err := db.Create(&siteSetting).Error; err != nil {
		return err
	}

	// Create Permissions
	permissions := []models.Permission{
		// Inventory permissions
		{Name: "inventory.view", Description: "View inventory items", Module: "inventory", Action: "view"},
		{Name: "inventory.create", Description: "Create inventory items", Module: "inventory", Action: "create"},
		{Name: "inventory.update", Description: "Update inventory items", Module: "inventory", Action: "update"},
		{Name: "inventory.delete", Description: "Delete inventory items", Module: "inventory", Action: "delete"},

		// Category permissions
		{Name: "category.view", Description: "View categories", Module: "category", Action: "view"},
		{Name: "category.create", Description: "Create categories", Module: "category", Action: "create"},
		{Name: "category.update", Description: "Update categories", Module: "category", Action: "update"},
		{Name: "category.delete", Description: "Delete categories", Module: "category", Action: "delete"},

		// Warehouse permissions
		{Name: "warehouse.view", Description: "View warehouses", Module: "warehouse", Action: "view"},
		{Name: "warehouse.create", Description: "Create warehouse", Module: "warehouse", Action: "create"},
		{Name: "warehouse.update", Description: "Update warehouse", Module: "warehouse", Action: "update"},
		{Name: "warehouse.delete", Description: "Delete warehouse", Module: "warehouse", Action: "delete"},

		// Employee permissions
		{Name: "employee.view", Description: "View employees", Module: "employee", Action: "view"},
		{Name: "employee.create", Description: "Create employee", Module: "employee", Action: "create"},
		{Name: "employee.update", Description: "Update employee", Module: "employee", Action: "update"},
		{Name: "employee.delete", Description: "Delete employee", Module: "employee", Action: "delete"},

		// Lead permissions
		{Name: "lead.view", Description: "View leads", Module: "lead", Action: "view"},
		{Name: "lead.create", Description: "Create lead", Module: "lead", Action: "create"},
		{Name: "lead.update", Description: "Update lead", Module: "lead", Action: "update"},
		{Name: "lead.delete", Description: "Delete lead", Module: "lead", Action: "delete"},

		// Project permissions
		{Name: "project.view", Description: "View projects", Module: "project", Action: "view"},
		{Name: "project.create", Description: "Create project", Module: "project", Action: "create"},
		{Name: "project.update", Description: "Update project", Module: "project", Action: "update"},
		{Name: "project.delete", Description: "Delete project", Module: "project", Action: "delete"},

		// PO permissions
		{Name: "po.view", Description: "View purchase orders", Module: "po", Action: "view"},
		{Name: "po.create", Description: "Create purchase order", Module: "po", Action: "create"},
		{Name: "po.update", Description: "Update purchase order", Module: "po", Action: "update"},
		{Name: "po.approve", Description: "Approve purchase order", Module: "po", Action: "approve"},
		{Name: "po.delete", Description: "Delete purchase order", Module: "po", Action: "delete"},
	}

	for _, permission := range permissions {
		if err := db.Create(&permission).Error; err != nil {
			return err
		}
	}

	// Assign all permissions to admin role
	var adminRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)
	var allPermissions []models.Permission
	db.Find(&allPermissions)
	db.Model(&adminRole).Association("Permissions").Append(&allPermissions)

	// Assign default menu visibility
	menuAssignments := map[string][]string{
		"admin": {
			"dashboard",
			"warehouses",
			"inventory",
			"inventory.items",
			"inventory.categories",
			"inventory.transactions",
			"employees",
			"leads",
			"projects",
			"purchase_orders",
			"support",
			"settings.profile",
			"settings.company",
			"settings.notifications",
			"settings.users",
			"settings.sites",
		},
		"manager": {
			"dashboard",
			"warehouses",
			"inventory",
			"inventory.items",
			"inventory.transactions",
			"employees",
			"leads",
			"projects",
			"purchase_orders",
			"support",
			"settings.profile",
			"settings.notifications",
		},
		"employee": {
			"dashboard",
			"inventory",
			"inventory.items",
			"support",
			"settings.profile",
		},
	}

	for roleName, keys := range menuAssignments {
		var role models.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			return err
		}
		for _, key := range keys {
			if err := db.Create(&models.RoleMenu{
				RoleID:  role.ID,
				MenuKey: key,
			}).Error; err != nil {
				return err
			}
		}
	}

	// Create admin user
	hashedPassword, _ := utils.HashPassword("admin123")
	adminUser := models.User{
		Email:    "admin@tatapps.com",
		Password: hashedPassword,
		FullName: "System Administrator",
		Phone:    "081234567890",
		RoleID:   adminRole.ID,
		IsActive: true,
	}
	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	// Create sample warehouse
	warehouse := models.Warehouse{
		Code:       "WH001",
		Name:       "Main Warehouse",
		Address:    "Jl. Contoh No. 123",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "12345",
		Phone:      "021-12345678",
		Email:      "warehouse@tatapps.com",
		ManagerID:  &adminUser.ID,
		IsActive:   true,
	}
	if err := db.Create(&warehouse).Error; err != nil {
		return err
	}

	// Create sample employee
	employee := models.Employee{
		EmployeeCode:   "EMP001",
		NIK:            "3173000000000001",
		FullName:       "John Doe",
		Email:          "john.doe@tatapps.com",
		Phone:          "081234567891",
		Address:        "Jl. Contoh No. 456",
		City:           "Jakarta",
		Province:       "DKI Jakarta",
		Department:     "Operations",
		JobTitle:       "Warehouse Manager",
		JoinDate:       time.Now(),
		EmploymentType: "full-time",
		Status:         "active",
		BasicSalary:    8000000,
		WarehouseID:    &warehouse.ID,
	}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}

	// Create sample lead
	lead := models.Lead{
		CompanyName:    "PT Contoh Jaya",
		ContactPerson:  "Jane Smith",
		Email:          "jane@contohjaya.com",
		Phone:          "081234567892",
		WhatsApp:       "081234567892",
		Address:        "Jl. Bisnis No. 789",
		City:           "Jakarta",
		Province:       "DKI Jakarta",
		Source:         "website",
		Industry:       "Manufacturing",
		Status:         "new",
		Priority:       "high",
		EstimatedValue: 50000000,
		AssignedToID:   adminUser.ID,
		Notes:          "Interested in warehouse management solution",
	}
	if err := db.Create(&lead).Error; err != nil {
		return err
	}

	log.Println("Database seeded successfully!")
	return nil
}
