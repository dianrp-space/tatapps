package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"tatapps/internal/config"
	"tatapps/internal/routes"
	"tatapps/internal/services/notification"
)

func main() {
	router := gin.New()
	cfg := config.LoadConfig()
	notifService := notification.NewNotificationService(cfg, nil)

	routes.SetupRoutes(router, nil, cfg, notifService)

	for _, r := range router.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}
}
