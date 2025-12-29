package admin

import (
	"net/http"
	"strings"

	"event-management-backend/internal/services/admin"

	"github.com/gin-gonic/gin"
)

type RoleWageHandler struct {
	service *admin.RoleWageService
}

func NewRoleWageHandler(service *admin.RoleWageService) *RoleWageHandler {
	return &RoleWageHandler{service: service}
}

// GET /admin/wages
func (h *RoleWageHandler) List(c *gin.Context) {
	wages, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wages"})
		return
	}

	c.JSON(http.StatusOK, wages)
}

// PUT /admin/wages/:role
func (h *RoleWageHandler) Update(c *gin.Context) {
	role := strings.ToLower(c.Param("role"))

	var req struct {
		Wage int64 `json:"wage"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.Update(role, req.Wage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role wage updated successfully"})
}
