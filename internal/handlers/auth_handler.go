package handlers

import (
	"event-management-backend/internal/domain/interfaces"
	"event-management-backend/internal/domain/models"
	"event-management-backend/internal/services/auth"
	"event-management-backend/internal/utils"
	"event-management-backend/internal/validations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserRepo    interfaces.UserRepository
	RefreshRepo interfaces.RefreshTokenRepository
	JWTService  *auth.JWTService
}

func NewAuthHandler(u interfaces.UserRepository, r interfaces.RefreshTokenRepository, j *auth.JWTService) *AuthHandler {
	return &AuthHandler{UserRepo: u, RefreshRepo: r, JWTService: j}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req validations.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and password are required"})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.UserRepo.FindByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if user.Status == models.StatusBlocked || user.Status == models.StatusDeleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not allowed to login"})
		return
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	accessToken, err := h.JWTService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}
	utils.SetAccessToken(c, accessToken)
	rawRefresh, hashedRefresh, expiresAt, err := h.JWTService.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
		return
	}
	_ = h.RefreshRepo.DeleteByUserID(user.ID)
	err = h.RefreshRepo.Save(&models.RefreshToken{
		UserID:      user.ID,
		TokenHashed: hashedRefresh,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}
	utils.SetRefreshToken(c, rawRefresh)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
			"role": user.Role,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	access, _ := c.Cookie("access_token")
	refresh, _ := c.Cookie("refresh_token")
	if access == "" && refresh == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}
	userID, exists := c.Get("user_id")
	if exists {
		_ = h.RefreshRepo.DeleteByUserID(userID.(uint))
	}
	utils.ClearAccessToken(c)
	utils.ClearRefreshToken(c)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              user.ID,
		"name":            user.Name,
		"phone":           user.Phone,
		"role":            user.Role,
		"branch":          user.Branch,
		"starting_point":  user.StartingPoint,
		"blood_group":     user.BloodGroup,
		"dob":             user.DOB,
		"photo":           user.Photo,
		"joined_at":       user.JoinedAt,
		"completed_work":  user.CompletedWork,
		"current_wage":    user.CurrentWage,
		"status":          user.Status,
	})
}