package admin

import (
	"event-management-backend/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingHandler struct {
	DB *gorm.DB
}

func NewSettingHandler(db *gorm.DB) *SettingHandler {
	return &SettingHandler{DB: db}
}

// GetSettings fetches all system configurations
func (h *SettingHandler) GetSettings(c *gin.Context) {
	var settings []models.SystemSetting
	h.DB.Find(&settings)
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSetting toggles a boolean setting
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"` // "true" or "false"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.SystemSetting
	if err := h.DB.Where("key = ?", req.Key).First(&setting).Error; err != nil {
		// Create if doesn't exist
		setting = models.SystemSetting{Key: req.Key, Value: req.Value}
		h.DB.Create(&setting)
	} else {
		// Update existing
		setting.Value = req.Value
		h.DB.Save(&setting)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully", "setting": setting})
}