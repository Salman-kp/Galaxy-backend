package routes

import (
	adminHandlers "event-management-backend/internal/handlers/admin"
	"event-management-backend/internal/middleware"
	"event-management-backend/internal/repository"
	"event-management-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	userRepo := repository.NewUserRepository()
	wageRepo := repository.NewRoleWageRepository()
	refreshRepo := repository.NewRefreshTokenRepository()
	jwtService := services.NewJWTService()

	userService := services.NewAdminUserService(userRepo, wageRepo)
	userHandler := adminHandlers.NewAdminUserHandler(userService)

	admin := r.Group("/admin")
	admin.Use(
		middleware.JWTAuthMiddleware(jwtService, refreshRepo),
		middleware.AdminMiddleware(),
	)

	admin.POST("/users", userHandler.CreateUser)
	admin.GET("/users", userHandler.ListUsers)
	admin.GET("/users/:id", userHandler.GetUser)
	admin.PUT("/users/:id", userHandler.UpdateUser)
	admin.PUT("/users/block/:id", userHandler.BlockUser)
	admin.PUT("/users/unblock/:id", userHandler.UnblockUser)
	admin.DELETE("/users/:id", userHandler.DeleteUser)
	admin.PUT("/users/reset-password/:id", userHandler.ResetPassword)
}