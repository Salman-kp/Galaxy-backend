package worker

import (
	"event-management-backend/internal/services/worker"
	"event-management-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkerBookingHandler struct {
	service *worker.WorkerBookingService
}

func NewWorkerBookingHandler(service *worker.WorkerBookingService) *WorkerBookingHandler {
	return &WorkerBookingHandler{service: service}
}

// ---------------- BOOK EVENT ----------------
// Worker books an event for their role
func (h *WorkerBookingHandler) BookEvent(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.BookEvent(userID, eventID, role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event booked successfully"})
}

// ---------------- LIST MY BOOKINGS ----------------
func (h *WorkerBookingHandler) ListMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookings, err := h.service.ListMyBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *WorkerBookingHandler) ListCompletedBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookings, err := h.service.ListCompletedBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch completed bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
