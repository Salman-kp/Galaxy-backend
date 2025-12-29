package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/validations"

	"github.com/gin-gonic/gin"
)

type AdminProfileHandler struct {
	service *admin.AdminUserService
}

func NewAdminProfileHandler(service *admin.AdminUserService) *AdminProfileHandler {
	return &AdminProfileHandler{service: service}
}

func (h *AdminProfileHandler) UpdateProfile(c *gin.Context) {
	adminID := c.GetUint("user_id")

	jsonData := c.PostForm("json")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is required"})
		return
	}

	var req validations.UpdateAdminSelfProfileRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var photoName string
	file, err := c.FormFile("photo")
	if err == nil && file != nil && file.Filename != "" {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "photo must be jpg/jpeg/png"})
			return
		}

		os.MkdirAll("uploads/users", os.ModePerm)
		photoName = fmt.Sprintf("user_%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads/users", photoName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo"})
			return
		}
	}

	user := &models.User{
		ID:            adminID,
		Name:          req.Name,
		Phone:         req.Phone,
		Branch:        req.Branch,
		StartingPoint: req.StartingPoint,
		BloodGroup:    req.BloodGroup,
	}

	if req.DOB != "" {
		if parsed, err := time.Parse("2006-01-02", req.DOB); err == nil {
			user.DOB = &parsed
		}
	}

	if photoName != "" {
		user.Photo = photoName
	}

	if err := h.service.UpdateUser(user); err != nil {
		if err.Error() == "no changes detected" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no changes detected"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}