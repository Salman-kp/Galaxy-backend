package routes

import (
	workerHandlers "event-management-backend/internal/handlers/worker"
	"event-management-backend/internal/middleware"
	"event-management-backend/internal/repository"
	"event-management-backend/internal/services/auth"
	"event-management-backend/internal/services/worker"

	"github.com/gin-gonic/gin"
)

func WorkerRoutes(r *gin.Engine) {
	// ---------------- Repositories ----------------
	refreshRepo := repository.NewRefreshTokenRepository()
	eventRepo := repository.NewEventRepository()
	bookingRepo := repository.NewBookingRepository()
	userRepo := repository.NewUserRepository()

	// ---------------- Services ----------------
	jwtService := auth.NewJWTService()
	eventService := worker.NewWorkerEventService(eventRepo)
	bookingService := worker.NewWorkerBookingService(
		bookingRepo,
		eventRepo,
		userRepo,
	)

	// ---------------- Handlers ----------------
	eventHandler := workerHandlers.NewWorkerEventHandler(eventService)
	bookingHandler := workerHandlers.NewWorkerBookingHandler(bookingService)

	// ---------------- Routes ----------------
	workerGroup := r.Group("/worker")
	workerGroup.Use(
		middleware.JWTAuthMiddleware(jwtService, refreshRepo),
		middleware.WorkerMiddleware(), // sub_captain, main_boy, junior_boy
	)

	// HOME
	workerGroup.GET("/events", eventHandler.ListEvents)
	workerGroup.GET("/events/:id", eventHandler.GetEvent)

	// BOOK EVENT
	workerGroup.POST("/events/:event_id/book", bookingHandler.BookEvent)

	// BOOKINGS
	workerGroup.GET("/bookings", bookingHandler.ListMyBookings)
	workerGroup.GET("/bookings/:booking_id", bookingHandler.GetBookingDetails)

	// COMPLETED
	workerGroup.GET("/bookings/completed", bookingHandler.ListCompletedBookings)
}