package main

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/repository"
	"event-management-backend/internal/routes"
	"event-management-backend/internal/seeders"
	"event-management-backend/migrations"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect DB
	config.ConnectDatabase()

	// Migrate DB
	if err := migrations.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Seeders
	seeders.SeedRBAC(config.DB)
	seeders.SeedRoleWages(config.DB)

	// Router
	router := gin.Default()
	router.RedirectTrailingSlash = false

	config.SetupWebConfig(router)

	router.Static("/uploads", "./uploads/users")
	api := router.Group("/api")

	// Auth routes
	routes.AuthRoutes(
		api,
		repository.NewUserRepository(),
		repository.NewRefreshTokenRepository(),
	)

	// Protected routes
	routes.AdminRoutes(api)
	routes.CaptainRoutes(api)
	routes.WorkerRoutes(api)

	// Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running at http://localhost%s", addr)
	router.Run(addr)
}
