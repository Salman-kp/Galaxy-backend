package middleware

import (
	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SystemGuard checks for global restrictions like Maintenance Mode
func SystemGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bypass check for Admin Login to allow them to fix things
		if c.Request.URL.Path == "/api/auth/login" {
			c.Next()
			return
		}

		roleValue, exists := c.Get("role")
		if !exists {
			c.Next()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.Next()
			return
		}

		//. Check "Boys Entering Disable" (Worker Access Disabled)
		if role != models.RoleAdmin {
			var workerSetting models.SystemSetting
			if err := config.DB.
				Where("key = ?", "worker_access_disabled").
				First(&workerSetting).Error; err == nil {

				if workerSetting.Value == "true" {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "Access is currently disabled by the administrator.",
					})
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
