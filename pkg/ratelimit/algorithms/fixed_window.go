package algorithms

import (
	"fmt"
	"shorten_url/pkg/ratelimit"
	"time"
)

type FixedWindowRateLimiter struct {
	storage ratelimit.Storage
	config  ratelimit.Config
}

// FixedWindowLimiter constructor
func NewFixedWindowRateLimiter(storage ratelimit.Storage, config ratelimit.Config) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		storage: storage,
		config:  config,
	}
}

// check if current request is allow
func (fw *FixedWindowRateLimiter) Allow(key string) (bool, ratelimit.RateLimitInfo, error) {
	now := time.Now()

	// get time param according to current window
	windowStart := now.Truncate(fw.config.Window)
	windowEnd := windowStart.Add(fw.config.Window)
	resetTime := windowEnd.Unix()

	windowKey := fmt.Sprintf("fixed window: %s:%d", key, windowStart.Unix())

	count, err := fw.storage.Increment(windowKey, fw.config.Window)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	remaining := fw.config.Limit - int(count)
	if remaining < 0 {
		remaining = 0

	}

	allowed := count <= int64(fw.config.Limit)

	retryAfter := 0
	if !allowed {
		retryAfter = max(int(windowEnd.Sub(now).Seconds()), 0)
	}

	info := ratelimit.RateLimitInfo{
		Limit:      fw.config.Limit,
		Remaining:  remaining,
		ResetTime:  resetTime,
		RetryAfter: retryAfter,
	}

	return allowed, info, nil

}

func (fw *FixedWindowRateLimiter) Reset(key string) error {
	now := time.Now()
	windowStart := now.Truncate(time.Duration(fw.config.Window))
	windowKey := fmt.Sprintf("fixed window:%s:%d", key, windowStart.Unix())

	return fw.storage.Delete(windowKey)
}
