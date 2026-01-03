package worker

import (
	"net/http"

	"event-management-backend/internal/services/worker"
	"event-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type WorkerEventHandler struct {
	service *worker.WorkerEventService
}

func NewWorkerEventHandler(service *worker.WorkerEventService) *WorkerEventHandler {
	return &WorkerEventHandler{service: service}
}

// ---------------- LIST AVAILABLE EVENTS ----------------
func (h *WorkerEventHandler) ListEvents(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return
	}

	events, err := h.service.ListAvailableEvents(userID, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ---------------- GET EVENT ----------------
func (h *WorkerEventHandler) GetEvent(c *gin.Context) {
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