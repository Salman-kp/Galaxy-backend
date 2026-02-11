package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/services/auth"
	"event-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(jwtService *auth.JWTService, refreshRepo interfaces.RefreshTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		header := c.GetHeader("Authorization")
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			accessToken = strings.TrimPrefix(header, "Bearer ")
		}
		if accessToken == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				accessToken = cookie
			}
		}
		claims, err := jwtService.ValidateAccessToken(accessToken)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
			c.Set("permissions", claims.Permissions)
			c.Next()
			return
		}
		if accessToken == "" || !errors.Is(err, jwt.ErrTokenExpired) {
			utils.ClearAccessToken(c)
			utils.ClearRefreshToken(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		rawRefresh, err := c.Cookie("refresh_token")
		if err != nil || rawRefresh == "" {
			utils.ClearAccessToken(c)
			utils.ClearRefreshToken(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "login required"})
			c.Abort()
			return
		}
		expiredClaims, err := jwtService.ParseExpiredAccessToken(accessToken)
		if err != nil || expiredClaims == nil {
			utils.ClearAccessToken(c)
			utils.ClearRefreshToken(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}
		hashed := utils.HashToken(rawRefresh)
		rt, err := refreshRepo.FindByHashedToken(hashed)
		if err != nil || rt.ExpiresAt.Before(time.Now().UTC()) {
			utils.ClearAccessToken(c)
			utils.ClearRefreshToken(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
			c.Abort()
			return
		}
		if rt.UserID != expiredClaims.UserID {
			utils.ClearAccessToken(c)
			utils.ClearRefreshToken(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token user mismatch"})
			c.Abort()
			return
		}
		perms := expiredClaims.Permissions
		if perms == nil {
			perms = []string{}
		}
		newAccess, err := jwtService.GenerateAccessToken(expiredClaims.UserID, expiredClaims.Role, perms)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
			c.Abort()
			return
		}
		utils.SetAccessToken(c, newAccess)
		c.Set("user_id", expiredClaims.UserID)
		c.Set("role", expiredClaims.Role)
		c.Set("permissions", perms)
		c.Next()
	}
}
