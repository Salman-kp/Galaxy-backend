package utils

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func isProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

func sameSiteMode() http.SameSite {
	if isProduction() {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

func SetAccessToken(c *gin.Context, token string) {
	maxAge := 60 * 60

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: sameSiteMode(),
	})
}

func ClearAccessToken(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: sameSiteMode(),
	})
}

func SetRefreshToken(c *gin.Context, token string) {
	maxAge := 7 * 24 * 60 * 60

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: sameSiteMode(),
	})
}

func ClearRefreshToken(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: sameSiteMode(),
	})
}
