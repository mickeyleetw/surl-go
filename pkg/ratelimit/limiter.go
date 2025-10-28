package ratelimit

// All rate limiting algorithms must implement this interface
type RateLimiter interface {
	Allow(key string) (allowed bool, info RateLimitInfo, err error)
	Reset(key string) error
}

// RateLimitInfo contains detailed information about rate limit status
// Used to return to middleware for setting HTTP headers
type RateLimitInfo struct {
	Limit      int
	Remaining  int
	ResetTime  int64
	RetryAfter int
}
