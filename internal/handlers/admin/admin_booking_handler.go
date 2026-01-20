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

	data, err := h.service.ListEventBookings(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ---------------- FILTER BY STATUS (ADMIN) ----------------
func (h *AdminBookingHandler) ListEventBookingsByStatus(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	status := c.Param("status")

	if eventID == 0 || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	data, err := h.service.ListEventBookingsByStatus(eventID, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ---------------- SEARCH BY NAME (ADMIN) ----------------
func (h *AdminBookingHandler) SearchEventBookingsByName(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	name := c.Query("name")

	if eventID == 0 || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	data, err := h.service.SearchEventBookingsByName(eventID, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ---------------- REMOVE USER FROM EVENT ----------------

func (h *AdminBookingHandler) RemoveUserFromEvent(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	bookingID := utils.ParseUintParam(c.Param("booking_id"))

	if eventID == 0 || bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	if err := h.service.RemoveUserFromEvent(eventID, bookingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user removed from event"})
}

// ---------------- UPDATE ATTENDANCE (ADMIN) ----------------
func (h *AdminBookingHandler) UpdateAttendance(c *gin.Context) {
	bookingID := utils.ParseUintParam(c.Param("booking_id"))
	if bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}
    
	var req validations.UpdateAttendanceRequest
		req.BookingID = bookingID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}


	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.UpdateAttendance(
		req.BookingID,
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

// ---------------- EVENT WAGE SUMMARY (ADMIN) ----------------
// RETURNS TOTALS OF ALL WAGE COLUMNS FOR AN EVENT
//
func (h *AdminBookingHandler) GetEventWageSummary(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	summary, err := h.service.GetEventWageSummary(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}