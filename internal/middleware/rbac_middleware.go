package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func HasPermission(requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, _ := c.Get("role")
        if role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
            c.Abort()
            return
        }

        perms, exists := c.Get("permissions")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "no permissions assigned"})
            c.Abort()
            return
        }

        permissions, ok := perms.([]string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "invalid permissions format"})
            c.Abort()
            return
        }

        for _, p := range permissions {
            if p == requiredPermission {
                c.Next()
                return
            }
        }
        
        c.JSON(http.StatusForbidden, gin.H{"error": "permission denied: " + requiredPermission})
        c.Abort()
    }
}