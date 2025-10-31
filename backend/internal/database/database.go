package database

import (
	"fmt"
	"log"
	"tatapps/internal/config"
	"tatapps/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
	return db
}

func AutoMigrate(db *gorm.DB) error {
	// Migrate in correct order to handle foreign key dependencies
	// Step 1: Base tables with no dependencies
	log.Println("Migrating Role and Permission tables...")
	if err := db.AutoMigrate(&models.Role{}, &models.Permission{}); err != nil {
		log.Println("Error migrating Role/Permission:", err)
		return err
	}
	log.Println("Role and Permission tables migrated successfully")

	// Step 2: User table (depends on Role)
	log.Println("Migrating User table...")
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Println("Error migrating User:", err)
		return err
	}
	log.Println("User table migrated successfully")

	// Step 3: Employee structure tables
	log.Println("Migrating Employee division table...")
	if err := db.AutoMigrate(&models.EmployeeDivision{}); err != nil {
		log.Println("Error migrating Employee division table:", err)
		return err
	}
	log.Println("Employee division table migrated successfully")

	log.Println("Migrating Employee position table...")
	if err := db.AutoMigrate(&models.EmployeePosition{}); err != nil {
		log.Println("Error migrating Employee position table:", err)
		return err
	}
	log.Println("Employee position table migrated successfully")

	log.Println("Migrating Employee table...")
	if err := db.AutoMigrate(&models.Employee{}); err != nil {
		log.Println("Error migrating Employee:", err)
		return err
	}
	log.Println("Employee table migrated successfully")

	log.Println("Migrating Warehouse, Lead, Project tables...")
	if err := db.AutoMigrate(&models.Warehouse{}, &models.Lead{}, &models.Project{}); err != nil {
		log.Println("Error migrating Warehouse/Lead/Project:", err)
		return err
	}
	log.Println("Warehouse, Lead, Project tables migrated successfully")

	log.Println("Migrating UserWarehouse table...")
	if err := db.AutoMigrate(&models.UserWarehouse{}); err != nil {
		log.Println("Error migrating UserWarehouse:", err)
		return err
	}
	log.Println("UserWarehouse table migrated successfully")

	// Step 4: Category table
	log.Println("Migrating Category table...")
	if err := db.AutoMigrate(&models.Category{}); err != nil {
		log.Println("Error migrating Category:", err)
		return err
	}
	log.Println("Category table migrated successfully")

	// Step 5: Tables with other dependencies
	log.Println("Migrating InventoryItem table...")
	if err := db.AutoMigrate(&models.InventoryItem{}); err != nil {
		log.Println("Error migrating InventoryItem:", err)
		return err
	}

	log.Println("Migrating PurchaseOrder and POItem tables...")
	if err := db.AutoMigrate(&models.PurchaseOrder{}, &models.POItem{}); err != nil {
		log.Println("Error migrating PurchaseOrder/POItem:", err)
		return err
	}

	log.Println("Migrating InventoryTransaction and Notification tables...")
	if err := db.AutoMigrate(&models.InventoryTransaction{}, &models.Notification{}); err != nil {
		log.Println("Error migrating InventoryTransaction/Notification:", err)
		return err
	}

	log.Println("Migrating NotificationSetting and NotificationHistory tables...")
	if err := db.AutoMigrate(&models.NotificationSetting{}, &models.NotificationHistory{}); err != nil {
		log.Println("Error migrating NotificationSetting/NotificationHistory:", err)
		return err
	}

	log.Println("Migrating SiteSetting table...")
	if err := db.AutoMigrate(&models.SiteSetting{}); err != nil {
		log.Println("Error migrating SiteSetting:", err)
		return err
	}

	log.Println("Migrating RoleMenu table...")
	if err := db.AutoMigrate(&models.RoleMenu{}); err != nil {
		log.Println("Error migrating RoleMenu:", err)
		return err
	}

	log.Println("All tables migrated successfully")
	return nil
}
