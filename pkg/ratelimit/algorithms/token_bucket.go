package algorithms

import (
	"fmt"
	"shorten_url/pkg/ratelimit"
	"time"
)

type TokenBucketRateLimiter struct {
	storage    ratelimit.Storage
	config     ratelimit.Config
	refillRate float64
}

func NewTokenBucketRateLimiter(storage ratelimit.Storage, config ratelimit.Config) *TokenBucketRateLimiter {
	refillRate := float64(config.Limit) / config.Window.Seconds()

	if config.BurstSize == 0 {
		config.BurstSize = config.Limit
	}

	return &TokenBucketRateLimiter{
		storage:    storage,
		config:     config,
		refillRate: refillRate,
	}
}

type bucketState struct {
	Tokens       float64
	LastRefillAt int64
}

func (tb *TokenBucketRateLimiter) Allow(key string) (bool, ratelimit.RateLimitInfo, error) {

	now := time.Now()
	bucketKey := fmt.Sprintf("token_bucket:%s", key)

	value, err := tb.storage.Get(bucketKey)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	var state bucketState
	if value == nil {
		state = bucketState{
			Tokens:       float64(tb.config.BurstSize),
			LastRefillAt: now.Unix(),
		}
	} else {
		state = value.(bucketState)

		elapsed := now.Unix() - state.LastRefillAt
		tokenToAdd := float64(elapsed) * tb.refillRate
		state.Tokens += tokenToAdd

		if state.Tokens > float64(tb.config.BurstSize) {
			state.Tokens = float64(tb.config.BurstSize)
		}

		state.LastRefillAt = now.Unix()

	}
	allowed := state.Tokens >= 1.0
	if allowed {
		state.Tokens -= 1.0
	}

	err = tb.storage.Set(bucketKey, state, 0)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	remaining := max(int(state.Tokens), 0)

	retryAfter := 0
	if !allowed {
		retryAfter = max(int(1.0/tb.refillRate), 1)
	}

	tokensNeeded := float64(tb.config.BurstSize) - state.Tokens
	secondToFull := int(tokensNeeded / tb.refillRate)
	resetTime := now.Add(time.Duration(secondToFull) * time.Second).Unix()

	info := ratelimit.RateLimitInfo{
		Limit:      tb.config.Limit,
		Remaining:  remaining,
		ResetTime:  resetTime,
		RetryAfter: retryAfter,
	}

	return allowed, info, nil
}

func (tb *TokenBucketRateLimiter) Reset(key string) error {
	bucketKey := fmt.Sprintf("token_bucket:%s", key)
	return tb.storage.Delete(bucketKey)
}
