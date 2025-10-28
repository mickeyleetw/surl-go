package server

import (
	"fmt"
	"net/http"

	"shorten_url/pkg/ratelimit"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

// ValidateRequest is a middleware that validates the request payload
func ValidateRequest[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		requestValidator := validator.New()
		if err := requestValidator.Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"validation_error": validationErrors.Error()})
			c.Abort()
			return
		}

		c.Set("validated_data", req)
		c.Next()
	}
}

// Ratelimiter middleware
func RequestRateLimiter(rateLimiter ratelimit.RateLimiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//ratelimiter key
		key := keyFunc(c)
		allowed, info, err := rateLimiter.Allow(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			c.Abort()
			return
		}
		// 設置 rate limit HTTP headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", info.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", info.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", info.ResetTime))

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%d", info.RetryAfter))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()

	}
}
