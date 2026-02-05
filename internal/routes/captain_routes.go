package routes

import (
	captainHandlers "event-management-backend/internal/handlers/captain"
	"event-management-backend/internal/middleware"
	"event-management-backend/internal/repository"
	"event-management-backend/internal/services/auth"
	"event-management-backend/internal/services/captain"

	"github.com/gin-gonic/gin"
)

func CaptainRoutes(r *gin.RouterGroup) {
	// ---------------- Repositories ----------------
	refreshRepo := repository.NewRefreshTokenRepository()
	eventRepo := repository.NewEventRepository()
	bookingRepo := repository.NewBookingRepository()
	userRepo := repository.NewUserRepository()

	// ---------------- Services ----------------
	jwtService := auth.NewJWTService()
	eventService := captain.NewCaptainEventService(eventRepo)
	bookingService := captain.NewCaptainBookingService(bookingRepo, eventRepo, userRepo)

	// ---------------- Handlers ----------------
	eventHandler := captainHandlers.NewCaptainEventHandler(eventService)
	bookingHandler := captainHandlers.NewCaptainBookingHandler(bookingService)

	// ---------------- Routes ----------------
	captainGroup := r.Group("/captain")
	captainGroup.Use(
		middleware.JWTAuthMiddleware(jwtService, refreshRepo),
		middleware.CaptainMiddleware(),
	)

	// EVENT ROUTES
	captainGroup.GET("/events",middleware.SystemGuard() ,eventHandler.ListEvents)
	captainGroup.GET("/events/:id", eventHandler.GetEvent)
	captainGroup.PUT("/events/start/:id", eventHandler.StartEvent)
	captainGroup.PUT("/events/complete/:id", eventHandler.CompleteEvent)

	// BOOK EVENT
	captainGroup.POST("/events/:event_id/book", bookingHandler.BookEvent)

	// BOOKING LISTS
	captainGroup.GET("/bookings/today", bookingHandler.ListTodayBookings)
	captainGroup.GET("/bookings/upcoming", bookingHandler.ListUpcomingBookings)
	captainGroup.GET("/bookings/completed", bookingHandler.ListCompletedBookings)

	// ATTENDANCE
	captainGroup.GET("/event-attendance/:event_id", bookingHandler.ListEventBookings)
	captainGroup.PUT("/event-attendance/:event_id", bookingHandler.UpdateAttendance)
	captainGroup.GET("/event-attendance/:event_id/status/:status", bookingHandler.ListEventBookingsByStatus)
	captainGroup.GET("/event-attendance/:event_id/search", bookingHandler.SearchEventBookingsByName)
	captainGroup.GET("/reports/events/:event_id/wages/summary",bookingHandler.GetEventWageSummary)
}