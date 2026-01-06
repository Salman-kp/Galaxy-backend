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

// ======================= BOOK EVENT =======================
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

// ======================= TODAY =======================
func (h *CaptainBookingHandler) ListTodayBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListTodayBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch today's bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ======================= UPCOMING =======================
func (h *CaptainBookingHandler) ListUpcomingBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListUpcomingBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ======================= COMPLETED =======================
func (h *CaptainBookingHandler) ListCompletedBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := h.service.ListCompletedBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch completed bookings"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ======================= EVENT ATTENDANCE =======================
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

// ======================= UPDATE ATTENDANCE =======================
func (h *CaptainBookingHandler) UpdateAttendance(c *gin.Context) {
	captainID := c.GetUint("user_id")

	var req validations.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateAttendance(
		captainID,
		req.BookingID,
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
// ======================= FILTER BY STATUS =======================
func (h *CaptainBookingHandler) ListEventBookingsByStatus(c *gin.Context) {
	captainID := c.GetUint("user_id")

	eventID := utils.ParseUintParam(c.Param("event_id"))
	status := c.Param("status")
	if eventID == 0 || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	data, err := h.service.ListEventBookingsByStatus(captainID, eventID, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ======================= SEARCH BY NAME =======================
func (h *CaptainBookingHandler) SearchEventBookingsByName(c *gin.Context) {
	captainID := c.GetUint("user_id")

	eventID := utils.ParseUintParam(c.Param("event_id"))
	name := c.Query("name")
	if eventID == 0 || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	data, err := h.service.SearchEventBookingsByName(captainID, eventID, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}