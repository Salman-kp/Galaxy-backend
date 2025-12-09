package middleware
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CaptainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "captain" {
			c.JSON(http.StatusForbidden, gin.H{"error": "captain only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
