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
// Captain books an event for himself
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
// All bookings (used by Home filtering)
//
func (h *CaptainBookingHandler) ListMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookings, err := h.service.ListMyBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

//
// ======================= LIST TODAY BOOKINGS =======================
// Today page
//
func (h *CaptainBookingHandler) ListTodayBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookings, err := h.service.ListTodayBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch today's bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

//
// ======================= LIST UPCOMING BOOKINGS =======================
// Booked (future) page
//
func (h *CaptainBookingHandler) ListUpcomingBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookings, err := h.service.ListUpcomingBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

//
// ======================= LIST COMPLETED BOOKINGS =======================
// Completed page
//
func (h *CaptainBookingHandler) ListCompletedBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	bookings, err := h.service.ListCompletedBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch completed bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

//
// ======================= LIST EVENT BOOKINGS =======================
// Attendance table (event context)
//
func (h *CaptainBookingHandler) ListEventBookings(c *gin.Context) {
	eventID := utils.ParseUintParam(c.Param("event_id"))
	if eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	bookings, err := h.service.ListEventBookings(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch event bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

//
// ======================= UPDATE ATTENDANCE =======================
// Captain updates attendance & wage for a booking
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