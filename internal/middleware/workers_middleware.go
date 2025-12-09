package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WorkerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == "sub_captain" || role == "main_boy" || role == "junior_boy" {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "worker only"})
		c.Abort()
	}
}
