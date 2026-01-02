package admin

import (
	"net/http"

	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/utils"
	"event-management-backend/internal/validations"

	"github.com/gin-gonic/gin"
)

type AdminBookingHandler struct {
	service *admin.AdminBookingService
}

func NewAdminBookingHandler(service *admin.AdminBookingService) *AdminBookingHandler {
	return &AdminBookingHandler{service: service}
}

// ---------------- LIST EVENT BOOKINGS ----------------

func (h *AdminBookingHandler) ListEventBookings(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	bookings, err := h.service.ListEventBookings(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// ---------------- REMOVE USER FROM EVENT ----------------

func (h *AdminBookingHandler) RemoveUserFromEvent(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	bookingID := utils.ParseUintParam(c.Param("booking_id"))
	if bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := h.service.RemoveUserFromEvent(eventID, bookingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user removed from event"})
}

// ---------------- UPDATE ATTENDANCE (ADMIN) ----------------
// RETURNS UPDATED BOOKING WITH TOTAL AMOUNT
//
func (h *AdminBookingHandler) UpdateAttendance(c *gin.Context) {
	bookingID := utils.ParseUintParam(c.Param("booking_id"))
	if bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var req validations.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Validate() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	booking, err := h.service.UpdateAttendance(
		bookingID,
		req.Status,
		req.TAAmount,
		req.BonusAmount,
		req.FineAmount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"booking_id":   booking.ID,
		"status":       booking.Status,
		"base_amount":  booking.BaseAmount,
		"extra_amount": booking.ExtraAmount,
		"ta_amount":    booking.TAAmount,
		"bonus_amount": booking.BonusAmount,
		"fine_amount":  booking.FineAmount,
		"total_amount": booking.TotalAmount,
	})
}