package admin

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"event-management-backend/internal/services/admin"
)

type AdminRoleHandler struct {
	service *admin.RoleService
}

func NewAdminRoleHandler(s *admin.RoleService) *AdminRoleHandler {
	return &AdminRoleHandler{service: s}
}

func (h *AdminRoleHandler) CreatePermission(c *gin.Context) {
	var body struct {
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePermission(body.Slug, body.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "permission created"})
}

func (h *AdminRoleHandler) ListPermissions(c *gin.Context) {
	perms, err := h.service.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}

func (h *AdminRoleHandler) CreateRole(c *gin.Context) {
	var body struct {
		Name          string `json:"name" binding:"required"`
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateRole(body.Name, body.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "role created"})
}

func (h *AdminRoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *AdminRoleHandler) GetRoleDetails(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := h.service.GetRoleDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (h *AdminRoleHandler) UpdateRole(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    
    var body struct {
        Name          string `json:"name" binding:"required"`
        PermissionIDs []uint `json:"permission_ids"`
    }

	if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    if err := h.service.UpdateRole(uint(id), body.Name, body.PermissionIDs); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

func (h *AdminRoleHandler) DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role deleted"})
}