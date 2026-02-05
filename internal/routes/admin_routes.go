package routes

import (
	"event-management-backend/internal/config"
	adminHandlers "event-management-backend/internal/handlers/admin"
	"event-management-backend/internal/middleware"
	"event-management-backend/internal/repository"
	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	// ---------------- Repositories ----------------
	userRepo := repository.NewUserRepository()
	wageRepo := repository.NewRoleWageRepository()
	refreshRepo := repository.NewRefreshTokenRepository()
	eventRepo := repository.NewEventRepository()
	bookingRepo := repository.NewBookingRepository()
	roleRepo := repository.NewRoleRepository()
	permRepo := repository.NewPermissionRepository()

	// ---------------- Services ----------------
	jwtService := auth.NewJWTService()

	userService := admin.NewAdminUserService(userRepo, wageRepo)
	eventService := admin.NewAdminEventService(eventRepo)
	bookingService := admin.NewAdminBookingService(bookingRepo, eventRepo)
	wageService := admin.NewWageService(bookingRepo, eventRepo)
	dashboardService := admin.NewDashboardService()
	roleWageService := admin.NewRoleWageService(wageRepo, userRepo)
	roleService := admin.NewRoleService(roleRepo, permRepo)

	// ---------------- Handlers ----------------
	userHandler := adminHandlers.NewAdminUserHandler(userService)
	eventHandler := adminHandlers.NewAdminEventHandler(eventService)
	bookingHandler := adminHandlers.NewAdminBookingHandler(bookingService)
	wageHandler := adminHandlers.NewAdminWageHandler(wageService)
	dashboardHandler := adminHandlers.NewAdminDashboardHandler(dashboardService)
	profileHandler := adminHandlers.NewAdminProfileHandler(userService)
	roleWageHandler := adminHandlers.NewRoleWageHandler(roleWageService)
	roleHandler := adminHandlers.NewAdminRoleHandler(roleService)
	settingHandler := adminHandlers.NewSettingHandler(config.DB)

	// ---------------- Routes ----------------
	adminGroup := r.Group("/admin")
	adminGroup.Use(
		middleware.JWTAuthMiddleware(jwtService, refreshRepo),
		middleware.AdminMiddleware(),
	)

	adminGroup.GET("/settings", middleware.HasPermission("system:manage"), settingHandler.GetSettings)
	adminGroup.PUT("/settings", middleware.HasPermission("system:manage"), settingHandler.UpdateSetting)

    // --- USER MANAGEMENT ---
	users := adminGroup.Group("/users")
	{
		users.POST("/",middleware.HasPermission("user:create"), userHandler.CreateUser)
		users.GET("/",middleware.HasPermission("user:view"), userHandler.ListUsers)
		users.GET("/role/:role",middleware.HasPermission("user:view"), userHandler.ListUsersByRole)
		users.GET("/search",middleware.HasPermission("user:view"), userHandler.SearchUsersByPhone)
		users.GET("/:id",middleware.HasPermission("user:view"), userHandler.GetUser)
		users.PUT("/:id",middleware.HasPermission("user:edit"), userHandler.UpdateUser)
		users.PUT("/block/:id", middleware.HasPermission("user:status"),userHandler.BlockUser)
		users.PUT("/unblock/:id",middleware.HasPermission("user:status"), userHandler.UnblockUser)
		users.DELETE("/:id/photo",middleware.HasPermission("user:edit"), userHandler.RemoveUserPhoto)
		users.DELETE("/:id",middleware.HasPermission("user:delete"), userHandler.DeleteUser)
		users.PUT("/reset-password/:id",middleware.HasPermission("user:password"), userHandler.ResetPassword)
	}

    // --- EVENT MANAGEMENT ---
	events := adminGroup.Group("/events")
	{
        events.GET("/", middleware.HasPermission("event:view"), eventHandler.ListEvents)
    	events.GET("/:id", middleware.HasPermission("event:view"), eventHandler.GetEvent)
		events.POST("/", middleware.HasPermission("event:create"), eventHandler.CreateEvent)
        events.PUT("/:id", middleware.HasPermission("event:edit"), eventHandler.UpdateEvent)
        events.DELETE("/:id", middleware.HasPermission("event:delete"), eventHandler.DeleteEvent)
        
        // Operational access LIFE CYCLE
        events.PUT("/start/:id", middleware.HasPermission("event:operate"), eventHandler.StartEvent)
        events.PUT("/complete/:id", middleware.HasPermission("event:operate"), eventHandler.CompleteEvent)
		events.PUT("/cancel/:id", middleware.HasPermission("event:operate"), eventHandler.CancelEvent)
	}

   // --- BOOKINGS & WAGES ---
	adminGroup.GET("/events/bookings/:event_id", middleware.HasPermission("event:view"), bookingHandler.ListEventBookings)
	adminGroup.DELETE("/events/bookings/:event_id/:booking_id", middleware.HasPermission("event:operate"), bookingHandler.RemoveUserFromEvent)
	adminGroup.PUT("/bookings/:booking_id/attendance", middleware.HasPermission("event:operate"), bookingHandler.UpdateAttendance)
	adminGroup.PUT("/bookings/:booking_id/wage", middleware.HasPermission("wage:edit"), wageHandler.OverrideWage)
	adminGroup.GET("/events/bookings/:event_id/status/:status",middleware.HasPermission("event:view"), bookingHandler.ListEventBookingsByStatus)
	adminGroup.GET("/events/bookings/:event_id/search",middleware.HasPermission("event:view"), bookingHandler.SearchEventBookingsByName)
	adminGroup.GET("/reports/events/:event_id/wages/summary", middleware.HasPermission("wage:view"), bookingHandler.GetEventWageSummary)
	
	
    // --- DASHBOARD & PROFILE ---
	adminGroup.GET("/dashboard/summary", middleware.HasPermission("dashboard:view"), dashboardHandler.GetSummary)
	adminGroup.GET("/dashboard/charts/monthly", middleware.HasPermission("dashboard:view"), dashboardHandler.GetMonthlyChart)
	adminGroup.GET("/dashboard/charts/daily", middleware.HasPermission("dashboard:view"), dashboardHandler.GetDailyChart)
	adminGroup.PUT("/profile", middleware.HasPermission("profile:edit"), profileHandler.UpdateProfile)

    // --- ROLE WAGES ---
	adminGroup.GET("/wages", middleware.HasPermission("managewages:view"), roleWageHandler.List)
	adminGroup.PUT("/wages/:role", middleware.HasPermission("managewages:view"), roleWageHandler.Update)

    // --- RBAC MANAGEMENT ---
    rbac := adminGroup.Group("/rbac")
    rbac.Use(middleware.HasPermission("rbac:view"))
    {
    rbac.POST("/users/invite", userHandler.CreateUser)

	rbac.GET("/permissions", roleHandler.ListPermissions)
    rbac.GET("/roles", roleHandler.ListRoles)
    rbac.GET("/roles/:id", roleHandler.GetRoleDetails)
    rbac.POST("/roles", roleHandler.CreateRole)
    rbac.PUT("/roles/:id", roleHandler.UpdateRole)
    rbac.DELETE("/roles/:id", roleHandler.DeleteRole)

    rbac.PUT("/update-role/:id", userHandler.UpdateUserRole) 
    rbac.DELETE("/admins/:id", userHandler.DeleteUser)
  }
}