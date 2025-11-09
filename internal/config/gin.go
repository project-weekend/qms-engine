package config

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/project-weekend/qms-engine/server/config"
)

// NewGinEngine initializes and configures a new Gin engine with middleware
func NewGinEngine(config *config.Config, log *slog.Logger) *gin.Engine {
	// Set Gin mode based on environment
	if config.Env == "production" || config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create new Gin engine
	engine := gin.New()

	// Add custom recovery middleware
	engine.Use(RecoveryMiddleware(log))

	// Add custom logging middleware
	engine.Use(LoggingMiddleware(log))

	// Add CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add health check endpoints
	engine.GET("/health", HealthCheckHandler())
	engine.GET("/ping", PingHandler())

	log.Info("Gin engine initialized successfully")
	return engine
}

// LoggingMiddleware creates a custom logging middleware using slog
func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log request details
		log.Info("HTTP request",
			"status", statusCode,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"latency", latency.Milliseconds(),
			"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}

// RecoveryMiddleware creates a custom recovery middleware using slog
func RecoveryMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					"error", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"ip", c.ClientIP(),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// HealthCheckHandler returns a handler for health check endpoint
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"usecase": "qms-engine",
			"time":    time.Now().Format(time.RFC3339),
		})
	}
}

// PingHandler returns a handler for ping endpoint
func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
