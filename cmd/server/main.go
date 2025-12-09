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
	// Connect to DB
	config.ConnectDatabase()
	//migrate db
	if err := migrations.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}
    seeders.SeedAdmin(config.DB)
	
	// Create router
	router := gin.Default()

	routes.AuthRoutes(
		router,
		repository.NewUserRepository(),
		repository.NewRefreshTokenRepository(),
	)

	routes.AdminRoutes(router)

	// Server port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running at http://localhost%s", addr)
	router.Run(addr)
}