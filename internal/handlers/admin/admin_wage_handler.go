package admin

import (
	"net/http"

	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/utils"
	"event-management-backend/internal/validations"

	"github.com/gin-gonic/gin"
)

type AdminWageHandler struct {
	service *admin.WageService
}

func NewAdminWageHandler(service *admin.WageService) *AdminWageHandler {
	return &AdminWageHandler{service: service}
}

// ---------------- OVERRIDE WAGE (ADMIN) ----------------
func (h *AdminWageHandler) OverrideWage(c *gin.Context) {
	bookingID := utils.ParseUintParam(c.Param("booking_id"))
	if bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var req validations.UpdateWageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.OverrideWage(
		bookingID,
		req.TAAmount,
		req.BonusAmount,
		req.FineAmount,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wage updated successfully"})
}
