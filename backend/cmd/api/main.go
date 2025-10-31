package main

import (
	"context"
	"log"
	"tatapps/internal/config"
	"tatapps/internal/database"
	"tatapps/internal/routes"
	"tatapps/internal/services/notification"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from backend directory
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg)

	// Auto migrate database
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data (only for development)
	if cfg.AppEnv == "development" {
		if err := database.SeedData(db); err != nil {
			log.Println("Warning: Failed to seed database:", err)
		}
	}

	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Shared services
	notifService := notification.NewNotificationService(cfg, db)
	scheduler := notification.NewLowStockScheduler(db, notifService)
	scheduler.Start()
	defer scheduler.Stop(context.Background())

	// Create Gin router
	router := gin.Default()

	// Serve static files from uploads directory
	router.Static("/uploads", "./uploads")

	// Setup routes
	routes.SetupRoutes(router, db, cfg, notifService)

	// Start server
	log.Printf("Server starting on port %s...", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
