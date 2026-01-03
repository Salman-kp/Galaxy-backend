package captain

import (
	"net/http"

	"event-management-backend/internal/services/captain"
	"event-management-backend/internal/utils"
	"event-management-backend/internal/validations"

	"github.com/gin-gonic/gin"
)

type CaptainBookingHandler struct {
	service *captain.CaptainBookingService
}

func NewCaptainBookingHandler(service *captain.CaptainBookingService) *CaptainBookingHandler {
	return &CaptainBookingHandler{service: service}
}

//
// ======================= BOOK EVENT =======================
//
func (h *CaptainBookingHandler) BookEvent(c *gin.Context) {
	userID := c.GetUint("user_id")

	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.BookEvent(userID, eventID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event booked successfully"})
}

//
// ======================= LIST MY BOOKINGS =======================
// Home / fallback list
//
func (h *CaptainBookingHandler) ListMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListMyBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ======================= LIST TODAY BOOKINGS =======================
// Today tab
//
func (h *CaptainBookingHandler) ListTodayBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListTodayBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch today's bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ======================= LIST UPCOMING BOOKINGS =======================
// Upcoming tab
//
func (h *CaptainBookingHandler) ListUpcomingBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListUpcomingBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ======================= LIST COMPLETED BOOKINGS =======================
// Completed tab
//
func (h *CaptainBookingHandler) ListCompletedBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListCompletedBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch completed bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ======================= LIST EVENT BOOKINGS =======================
// Attendance table (only captain of event)
//
func (h *CaptainBookingHandler) ListEventBookings(c *gin.Context) {
	captainID := c.GetUint("user_id")

	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	data, err := h.service.ListEventBookings(captainID, eventID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// ======================= UPDATE ATTENDANCE =======================
//
func (h *CaptainBookingHandler) UpdateAttendance(c *gin.Context) {
	bookingID := utils.ParseUintParam(c.Param("booking_id"))
	if bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var req validations.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	captainID := c.GetUint("user_id")
	if err := h.service.UpdateAttendance(
		captainID,
		bookingID,
		req.Status,
		req.TAAmount,
		req.BonusAmount,
		req.FineAmount,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance updated successfully"})
}