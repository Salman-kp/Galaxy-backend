package captain

import (
	"net/http"

	"event-management-backend/internal/services/captain"
	"event-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CaptainEventHandler struct {
	service *captain.CaptainEventService
}

func NewCaptainEventHandler(service *captain.CaptainEventService) *CaptainEventHandler {
	return &CaptainEventHandler{service: service}
}

// ---------------- LIST AVAILABLE EVENTS ----------------
func (h *CaptainEventHandler) ListEvents(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	events, err := h.service.ListAvailableEvents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ---------------- GET EVENT ----------------
func (h *CaptainEventHandler) GetEvent(c *gin.Context) {
	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := h.service.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// ---------------- START EVENT ----------------
func (h *CaptainEventHandler) StartEvent(c *gin.Context) {
	captainID := c.GetUint("user_id")
	if captainID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.StartEvent(captainID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event started"})
}

// ---------------- COMPLETE EVENT ----------------
func (h *CaptainEventHandler) CompleteEvent(c *gin.Context) {
	captainID := c.GetUint("user_id")
	if captainID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.CompleteEvent(captainID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event completed"})
}