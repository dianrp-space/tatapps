package routes

import (
	"tatapps/internal/config"
	"tatapps/internal/handlers"
	"tatapps/internal/middleware"
	"tatapps/internal/services/notification"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config, notifService *notification.NotificationService) {
	// CORS middleware
	router.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg, notifService)
	warehouseHandler := handlers.NewWarehouseHandler(db)
	poHandler := handlers.NewPOHandler(db)
	inventoryHandler := handlers.NewInventoryHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	userHandler := handlers.NewUserHandler(db, notifService, cfg)
	employeeHandler := handlers.NewEmployeeHandler(db)
	settingsHandler := handlers.NewSettingsHandler(db, notifService, cfg)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
		public.GET("/settings/site", settingsHandler.GetSiteSettings)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// Profile
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.PUT("/users/profile", authHandler.UpdateProfile)
		protected.PUT("/users/change-password", authHandler.ChangePassword)

		// Warehouses
		warehouses := protected.Group("/warehouses")
		{
			warehouses.GET("", warehouseHandler.GetAll)
			warehouses.GET("/:id", warehouseHandler.GetByID)
			warehouses.POST("", middleware.AdminOrManager(), warehouseHandler.Create)
			warehouses.PUT("/:id", middleware.AdminOrManager(), warehouseHandler.Update)
			warehouses.DELETE("/:id", middleware.AdminOnly(), warehouseHandler.Delete)
		}

		// Purchase Orders
		purchaseOrders := protected.Group("/purchase-orders")
		{
			purchaseOrders.GET("", poHandler.GetAll)
			purchaseOrders.GET("/:id", poHandler.GetByID)
			purchaseOrders.POST("", poHandler.Create)
			purchaseOrders.PUT("/:id", poHandler.Update)
			purchaseOrders.POST("/:id/approve", middleware.AdminOrManager(), poHandler.Approve)
			purchaseOrders.POST("/:id/reject", middleware.AdminOrManager(), poHandler.Reject)
		}

		// Inventory
		inventory := protected.Group("/inventory")
		{
			inventory.GET("", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.GetAllItems)
			inventory.GET("/low-stock", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.GetLowStockItems)
			inventory.GET("/transactions", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.GetAllTransactions)
			inventory.DELETE("/transactions/:id", middleware.RequirePermission(db, "inventory.delete"), inventoryHandler.DeleteTransaction)
			inventory.GET("/import/template", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.DownloadImportTemplate)
			inventory.POST("/import/csv", middleware.RequirePermission(db, "inventory.create"), inventoryHandler.ImportItemsFromCSV)
			inventory.GET("/export/csv", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.ExportItemsToCSV)
			inventory.GET("/export/pdf", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.ExportItemsToPDF)
			inventory.GET("/transactions/export/csv", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.ExportTransactionsToCSV)
			inventory.GET("/transactions/export/pdf", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.ExportTransactionsToPDF)
			inventory.GET("/items/:id", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.GetItemByID)
			inventory.GET("/items/:id/transactions", middleware.RequirePermission(db, "inventory.view"), inventoryHandler.GetItemTransactions)
			inventory.POST("/items", middleware.RequirePermission(db, "inventory.create"), inventoryHandler.CreateItem)
			inventory.DELETE("/items", middleware.RequirePermission(db, "inventory.delete"), inventoryHandler.DeleteItemsBatch)
			inventory.PUT("/items/:id", middleware.RequirePermission(db, "inventory.update"), inventoryHandler.UpdateItem)
			inventory.DELETE("/items/:id", middleware.RequirePermission(db, "inventory.delete"), inventoryHandler.DeleteItem)
			inventory.POST("/items/:id/transactions", middleware.RequireAnyPermission(db, "inventory.update", "inventory.create"), inventoryHandler.RecordTransaction)
		}

		// Categories
		categories := protected.Group("/categories")
		{
			categories.GET("", middleware.RequirePermission(db, "category.view"), categoryHandler.GetAllCategories)
			categories.GET("/:id", middleware.RequirePermission(db, "category.view"), categoryHandler.GetCategoryByID)
			categories.POST("", middleware.RequirePermission(db, "category.create"), categoryHandler.CreateCategory)
			categories.PUT("/:id", middleware.RequirePermission(db, "category.update"), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.RequirePermission(db, "category.delete"), categoryHandler.DeleteCategory)
		}

		// Settings
		settings := protected.Group("/settings")
		{
			settings.GET("/site/admin", middleware.AdminOnly(), settingsHandler.GetSiteSettingsAdmin)
			settings.GET("/notifications", settingsHandler.GetNotificationSettings)
			settings.PUT("/notifications", settingsHandler.UpdateNotificationSettings)
			settings.PUT("/site", middleware.AdminOnly(), settingsHandler.UpdateSiteSettings)
			settings.GET("/database/backup", middleware.AdminOnly(), settingsHandler.BackupDatabase)
			settings.POST("/database/restore", middleware.AdminOnly(), settingsHandler.RestoreDatabase)

			// User management
			users := settings.Group("/users")
			{
				users.GET("", middleware.RequirePermission(db, "employee.view"), userHandler.GetAll)
				users.GET("/:id", middleware.RequirePermission(db, "employee.view"), userHandler.GetByID)
				users.POST("", middleware.RequirePermission(db, "employee.create"), userHandler.CreateUser)
				users.PUT("/:id", middleware.RequirePermission(db, "employee.update"), userHandler.UpdateUser)
				users.PUT("/:id/status", middleware.RequirePermission(db, "employee.update"), userHandler.UpdateUserStatus)
				users.DELETE("/:id", middleware.RequirePermission(db, "employee.delete"), userHandler.DeleteUser)
			}

			roles := settings.Group("/roles")
			{
				roles.GET("", middleware.AdminOnly(), userHandler.GetRoles)
				roles.GET("/menu-options", middleware.AdminOnly(), userHandler.GetRoleMenuOptions)
				roles.GET("/permission-options", middleware.AdminOnly(), userHandler.GetRolePermissionOptions)
				roles.POST("", middleware.AdminOnly(), userHandler.CreateRole)
				roles.PUT("/:id", middleware.AdminOnly(), userHandler.UpdateRole)
				roles.DELETE("/:id", middleware.AdminOnly(), userHandler.DeleteRole)
			}
		}

		// Roles
		protected.GET("/roles", userHandler.GetRoles)

		// Employees
		employees := protected.Group("/employees")
		{
			employees.GET("", middleware.RequirePermission(db, "employee.view"), employeeHandler.ListEmployees)
			employees.GET("/:id", middleware.RequirePermission(db, "employee.view"), employeeHandler.GetEmployee)
			employees.POST("", middleware.RequirePermission(db, "employee.create"), employeeHandler.CreateEmployee)
			employees.PUT("/:id", middleware.RequirePermission(db, "employee.update"), employeeHandler.UpdateEmployee)
			employees.DELETE("", middleware.RequirePermission(db, "employee.delete"), employeeHandler.DeleteEmployeesBatch)
			employees.DELETE("/:id", middleware.RequirePermission(db, "employee.delete"), employeeHandler.DeleteEmployee)

			divisions := employees.Group("/divisions")
			{
				divisions.GET("", middleware.RequirePermission(db, "employee.view"), employeeHandler.ListDivisions)
				divisions.POST("", middleware.RequirePermission(db, "employee.create"), employeeHandler.CreateDivision)
				divisions.PUT(":id", middleware.RequirePermission(db, "employee.update"), employeeHandler.UpdateDivision)
				divisions.DELETE(":id", middleware.RequirePermission(db, "employee.delete"), employeeHandler.DeleteDivision)
			}

			positions := employees.Group("/positions")
			{
				positions.GET("", middleware.RequirePermission(db, "employee.view"), employeeHandler.ListPositions)
				positions.POST("", middleware.RequirePermission(db, "employee.create"), employeeHandler.CreatePosition)
				positions.PUT(":id", middleware.RequirePermission(db, "employee.update"), employeeHandler.UpdatePosition)
				positions.DELETE(":id", middleware.RequirePermission(db, "employee.delete"), employeeHandler.DeletePosition)
			}
		}

		// Notifications
		notifications := protected.Group("/notifications")
		{
			notifications.POST("/test", settingsHandler.SendTestNotification)
			notifications.POST("/check-low-stock", settingsHandler.CheckLowStock)
			notifications.GET("/history", settingsHandler.GetNotificationHistory)
		}

		// Additional routes for other modules can be added here
		// - Employees
		// - Leads
		// - Projects
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
