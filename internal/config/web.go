package config

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupWebConfig(r *gin.Engine) {

	originsEnv := os.Getenv("CORS_ORIGINS")

	var origins []string
	if originsEnv != "" {
		origins = strings.Split(originsEnv, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.SetTrustedProxies(nil)
}