package admin

import (
	"net/http"
	"strconv"
	"time"

	"event-management-backend/internal/services/admin"

	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct {
	service *admin.DashboardService
}

func NewAdminDashboardHandler(service *admin.DashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{service: service}
}

// DASHBOARD SUMMARY
// Top 5 boxes
func (h *AdminDashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// MONTH-BASED CHART
func (h *AdminDashboardHandler) GetMonthlyChart(c *gin.Context) {
	yearStr := c.Query("year")
	year := time.Now().Year()

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	data, err := h.service.GetMonthlyEventChart(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch monthly chart"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// DATE-BASED CHART
// (Admin clicks a month)
func (h *AdminDashboardHandler) GetDailyChart(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}

	year, err1 := strconv.Atoi(yearStr)
	month, err2 := strconv.Atoi(monthStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year or month"})
		return
	}

	data, err := h.service.GetDailyEventChart(year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch daily chart"})
		return
	}

	c.JSON(http.StatusOK, data)
}

