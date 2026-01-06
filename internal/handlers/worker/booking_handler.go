package worker

import (
	"net/http"

	"event-management-backend/internal/services/worker"
	"event-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type WorkerBookingHandler struct {
	service *worker.WorkerBookingService
}

func NewWorkerBookingHandler(service *worker.WorkerBookingService) *WorkerBookingHandler {
	return &WorkerBookingHandler{service: service}
}

//
// ---------------- BOOK EVENT ----------------
//
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

//
// ---------------- LIST MY BOOKINGS ----------------
// Booked page (upcoming + ongoing)
//
func (h *WorkerBookingHandler) ListMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListMyBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ---------------- LIST COMPLETED BOOKINGS ----------------
// Completed page
//
func (h *WorkerBookingHandler) ListCompletedBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListCompletedBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch completed bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ---------------- GET BOOKING DETAILS ----------------
// Used by:
// - Booked details page
// - Completed details page
//
func (h *WorkerBookingHandler) GetBookingDetails(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookingID := utils.ParseUintParam(c.Param("booking_id"))
	if bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	data, err := h.service.GetBookingDetails(userID, bookingID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}