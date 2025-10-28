package ratelimit

import (
	"fmt"
	"time"
)

// Config defines the configuration parameters for the rate limiter
type Config struct {
	// Algorithm specifies the rate limiting algorithm to use
	// Valid values: "fixed_window", "token_bucket", "sliding_window", "leaky_bucket"
	Algorithm string

	// Storage specifies the storage backend to use
	// Valid values: "memory", "redis"
	Storage string

	// Limit is the rate limit threshold (e.g., allow 5 requests per minute)
	Limit int

	// Window is the time window (e.g., 1 minute)
	Window time.Duration

	// BurstSize is the burst capacity for Token Bucket algorithm (optional)
	// Not used by other algorithms, can be set to 0
	BurstSize int
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Limit <= 0 {
		return fmt.Errorf("limit must be positive, got: %d", c.Limit)
	}

	if c.Window <= 0 {
		return fmt.Errorf("window must be positive, got: %v", c.Window)
	}

	validAlgorithms := map[string]bool{
		"fixed_window":   true,
		"token_bucket":   true,
		"sliding_window": true,
		"leaky_bucket":   true,
	}

	if !validAlgorithms[c.Algorithm] {
		return fmt.Errorf("invalid algorithm: %s", c.Algorithm)
	}

	validStorages := map[string]bool{
		"memory": true,
		"redis":  true,
	}

	if !validStorages[c.Storage] {
		return fmt.Errorf("invalid storage: %s", c.Storage)
	}

	return nil
}
