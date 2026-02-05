package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/services/admin"
	"event-management-backend/internal/validations"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	service *admin.AdminUserService
}

func NewAdminUserHandler(service *admin.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{service: service}
}

func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	jsonData := c.PostForm("json")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is required"})
		return
	}

	var req validations.CreateUserRequest
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
	if err == nil {
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
		Name:          req.Name,
		Phone:         req.Phone,
		Password:      req.Password,
		Role:          req.Role,
		Branch:        req.Branch,
		StartingPoint: req.StartingPoint,
		BloodGroup:    req.BloodGroup,
		Photo:         photoName,
	}
    if req.Role == models.RoleAdmin {
        user.AdminRoleID = req.AdminRoleID
    } else {
        user.AdminRoleID = nil 
    }

	if req.DOB != "" {
		if parsed, err := time.Parse("2006-01-02", req.DOB); err == nil {
			user.DOB = &parsed
		}
	}

	if err := h.service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")

	users, err := h.service.ListUsers(role, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id := parseID(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	id := parseID(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

    existingUser, err := h.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	jsonData := c.PostForm("json")
	if jsonData == "" {
		fmt.Println("UpdateUser Error: 'json' field is empty in form-data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is required"})
		return
	}

	var req validations.UpdateUserRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		fmt.Printf("UpdateUser Error: JSON Unmarshal failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if req.Role == models.RoleAdmin && req.AdminRoleID == nil {
		if existingUser.AdminRoleID != nil {
			req.AdminRoleID = existingUser.AdminRoleID
		}
	}

	if err := req.Validate(); err != nil {
		fmt.Printf("UpdateUser Error: Validation failed: %v\n", err)
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
		ID:            id,
		Name:          req.Name,
		Phone:         req.Phone,
		Role:          req.Role,
		Branch:        req.Branch,
		StartingPoint: req.StartingPoint,
		BloodGroup:    req.BloodGroup,
		Status:        req.Status,
	}

	if req.Role == models.RoleAdmin {
        user.AdminRoleID = req.AdminRoleID
    } else {
        user.AdminRoleID = nil 
    }

	if req.DOB != "" {
		if parsed, err := time.Parse("2006-01-02", req.DOB); err == nil {
			user.DOB = &parsed
		}
	}

	if photoName != "" {
		user.Photo = photoName
	}

	err = h.service.UpdateUser(user)
	if err != nil {
		if photoName != "" {
			os.Remove(filepath.Join("uploads/users", photoName))
		}
		
		if err.Error() == "no changes detected" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no changes detected"})
			return
		}

		fmt.Printf("UpdateUser Error: Service Update failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if photoName != "" && existingUser.Photo != "" {
		if photoName != existingUser.Photo {
			oldPath := filepath.Join("uploads/users", existingUser.Photo)
			_ = os.Remove(oldPath)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *AdminUserHandler) RemoveUserPhoto(c *gin.Context) {
    id := parseID(c.Param("id"))
    if id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    photoName, err := h.service.RemoveUserPhoto(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if photoName != "" {
        oldPath := filepath.Join("uploads/users", photoName)
        _ = os.Remove(oldPath)
    }
    c.JSON(http.StatusOK, gin.H{"message": "photo removed successfully"})
}

func (h *AdminUserHandler) BlockUser(c *gin.Context) {
	id := parseID(c.Param("id"))
	if err := h.service.BlockUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user blocked"})
}

func (h *AdminUserHandler) UnblockUser(c *gin.Context) {
	id := parseID(c.Param("id"))
	if err := h.service.UnblockUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user unblocked"})
}

func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	id := parseID(c.Param("id"))
	if err := h.service.SoftDeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *AdminUserHandler) ResetPassword(c *gin.Context) {
	id := parseID(c.Param("id"))

	var body struct {
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || len(body.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	if err := h.service.ResetPassword(id, body.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

func (h *AdminUserHandler) ListUsersByRole(c *gin.Context) {
	role := c.Param("role")

	if !models.ValidateRole(role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	data, err := h.service.ListUsersByRole(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AdminUserHandler) SearchUsersByPhone(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone is required"})
		return
	}

	data, err := h.service.SearchUsersByPhone(phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AdminUserHandler) UpdateUserRole(c *gin.Context) {
	id := parseID(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req validations.UpdateUserClearanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateUserRole(id, req.Role, req.AdminRoleID)
	if err != nil {
		if err.Error() == "no changes detected in clearance" {
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user role and clearance updated successfully"})
}

func parseID(s string) uint {
	id64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id64)
}
