package admin

import (
	"net/http"

	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/utils"
	"event-management-backend/internal/validations"

	"github.com/gin-gonic/gin"
)

type AdminEventHandler struct {
	service *admin.AdminEventService
}

func NewAdminEventHandler(service *admin.AdminEventService) *AdminEventHandler {
	return &AdminEventHandler{service: service}
}

// ---------------- CREATE ----------------

func (h *AdminEventHandler) CreateEvent(c *gin.Context) {
	var req validations.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &models.Event{
		EventName:           req.Name,
		Date:                req.Date,
		TimeSlot:            req.TimeSlot,
		ReportingTime:       req.ReportingTime,
		WorkType:            req.WorkType,
		LocationLink:        req.LocationLink,
		RequiredCaptains:    req.RequiredCaptains,
		RequiredSubCaptains: req.RequiredSubCaptains,
		RequiredMainBoys:    req.RequiredMainBoys,
		RequiredJuniors:     req.RequiredJuniors,
		LongWork:            req.LongWork,
		TransportProvided:   req.TransportProvided,
		TransportType:       req.TransportType,
		ExtraWageAmount:     req.ExtraWageAmount,
	}

	if err := h.service.CreateEvent(event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "event created successfully",
		"id":      event.ID,
	})
}

// ---------------- LIST ----------------

func (h *AdminEventHandler) ListEvents(c *gin.Context) {
	status := c.Query("status")
	date := c.Query("date")

	events, err := h.service.ListEvents(status, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ---------------- GET ----------------

func (h *AdminEventHandler) GetEvent(c *gin.Context) {
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

// ---------------- UPDATE ----------------

func (h *AdminEventHandler) UpdateEvent(c *gin.Context) {
	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req validations.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &models.Event{
		ID:                  id,
		EventName:           req.Name,
		Date:                req.Date,
		TimeSlot:            req.TimeSlot,
		ReportingTime:       req.ReportingTime,
		WorkType:            req.WorkType,
		LocationLink:        req.LocationLink,
		RequiredCaptains:    req.RequiredCaptains,
		RequiredSubCaptains: req.RequiredSubCaptains,
		RequiredMainBoys:    req.RequiredMainBoys,
		RequiredJuniors:     req.RequiredJuniors,
		LongWork:            req.LongWork,
		TransportProvided:   req.TransportProvided,
		TransportType:       req.TransportType,
		ExtraWageAmount:     req.ExtraWageAmount,
	}

	if err := h.service.UpdateEvent(event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

// ---------------- START ----------------

func (h *AdminEventHandler) StartEvent(c *gin.Context) {
	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.StartEvent(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event started"})
}

// ---------------- COMPLETE ----------------

func (h *AdminEventHandler) CompleteEvent(c *gin.Context) {
	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.CompleteEvent(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event completed"})
}

// ---------------- CANCEL ----------------

func (h *AdminEventHandler) CancelEvent(c *gin.Context) {
	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.CancelEvent(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event cancelled"})
}

// ---------------- DELETE ----------------

func (h *AdminEventHandler) DeleteEvent(c *gin.Context) {
	id := utils.ParseUintParam(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.service.DeleteEvent(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}
