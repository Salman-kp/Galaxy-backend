package routes

import (
	adminHandlers "event-management-backend/internal/handlers/admin"
	"event-management-backend/internal/middleware"
	"event-management-backend/internal/repository"
	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	// ---------------- Repositories ----------------
	userRepo := repository.NewUserRepository()
	wageRepo := repository.NewRoleWageRepository()
	refreshRepo := repository.NewRefreshTokenRepository()
	eventRepo := repository.NewEventRepository()
	bookingRepo := repository.NewBookingRepository()

	// ---------------- Services ----------------
	jwtService := auth.NewJWTService()

	userService := admin.NewAdminUserService(userRepo, wageRepo)
	eventService := admin.NewAdminEventService(eventRepo)
	bookingService := admin.NewAdminBookingService(bookingRepo, eventRepo)
	wageService := admin.NewWageService(bookingRepo)
	dashboardService := admin.NewDashboardService()
	roleWageService := admin.NewRoleWageService(wageRepo)

	// ---------------- Handlers ----------------
	userHandler := adminHandlers.NewAdminUserHandler(userService)
	eventHandler := adminHandlers.NewAdminEventHandler(eventService)
	bookingHandler := adminHandlers.NewAdminBookingHandler(bookingService)
	wageHandler := adminHandlers.NewAdminWageHandler(wageService)
	dashboardHandler := adminHandlers.NewAdminDashboardHandler(dashboardService)
	profileHandler := adminHandlers.NewAdminProfileHandler(userService)
	roleWageHandler := adminHandlers.NewRoleWageHandler(roleWageService)

	// ---------------- Routes ----------------
	adminGroup := r.Group("/admin")
	adminGroup.Use(
		middleware.JWTAuthMiddleware(jwtService, refreshRepo),
		middleware.AdminMiddleware(),
	)

	// USER ROUTES
	adminGroup.POST("/users", userHandler.CreateUser)
	adminGroup.GET("/users", userHandler.ListUsers)
	adminGroup.GET("/users/:id", userHandler.GetUser)
	adminGroup.PUT("/users/:id", userHandler.UpdateUser)
	adminGroup.PUT("/users/block/:id", userHandler.BlockUser)
	adminGroup.PUT("/users/unblock/:id", userHandler.UnblockUser)
	adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
	adminGroup.PUT("/users/reset-password/:id", userHandler.ResetPassword)

	// EVENT ROUTES
	adminGroup.POST("/events", eventHandler.CreateEvent)
	adminGroup.GET("/events", eventHandler.ListEvents)
	adminGroup.GET("/events/:id", eventHandler.GetEvent)
	adminGroup.PUT("/events/:id", eventHandler.UpdateEvent)
	adminGroup.DELETE("/events/:id", eventHandler.DeleteEvent)

	// Event lifecycle control
	adminGroup.PUT("/events/start/:id", eventHandler.StartEvent)
	adminGroup.PUT("/events/complete/:id", eventHandler.CompleteEvent)
	adminGroup.PUT("/events/cancel/:id", eventHandler.CancelEvent)

	// BOOKING ROUTES (GIN-SAFE)
	adminGroup.GET("/events/bookings/:event_id", bookingHandler.ListEventBookings)
	adminGroup.DELETE("/events/bookings/:event_id/:booking_id", bookingHandler.RemoveUserFromEvent)
	adminGroup.PUT("/bookings/:booking_id/attendance", bookingHandler.UpdateAttendance)
	adminGroup.PUT("/bookings/:booking_id/wage", wageHandler.OverrideWage)

	// DASHBOARD ROUTES
	adminGroup.GET("/dashboard/summary", dashboardHandler.GetSummary)
	adminGroup.GET("/dashboard/charts/monthly", dashboardHandler.GetMonthlyChart)
	adminGroup.GET("/dashboard/charts/daily", dashboardHandler.GetDailyChart)
	adminGroup.GET("/dashboard/today", dashboardHandler.GetTodayEvents)

	// PROFILE
	adminGroup.PUT("/profile", profileHandler.UpdateProfile)

	// ROLE WAGE ROUTES
	adminGroup.GET("/wages", roleWageHandler.List)
	adminGroup.PUT("/wages/:role", roleWageHandler.Update)
}