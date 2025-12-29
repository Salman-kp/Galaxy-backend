package routes

import (
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/handlers"
	"event-management-backend/internal/middleware"
	"event-management-backend/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, userRepo interfaces.UserRepository, refreshRepo interfaces.RefreshTokenRepository) {
	jwtService := auth.NewJWTService()
	authHandler := handlers.NewAuthHandler(userRepo, refreshRepo, jwtService)

	r.POST("/auth/login", authHandler.Login)

	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuthMiddleware(jwtService, refreshRepo))

	auth.POST("/logout", authHandler.Logout)
	auth.GET("/profile", authHandler.Profile)
}