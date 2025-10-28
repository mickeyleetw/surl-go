package algorithms

import (
	"fmt"
	"shorten_url/pkg/ratelimit"
	"time"
)

type LeakyBucketRateLimiter struct {
	storage   ratelimit.Storage
	config    ratelimit.Config
	leakyRate float64
}

func NewLeakyBucketRateLimiter(config ratelimit.Config, storage ratelimit.Storage) *LeakyBucketRateLimiter {
	leakyRate := float64(config.Limit) / config.Window.Seconds()

	capacity := config.Limit
	if config.BurstSize > 0 {
		capacity = config.BurstSize
	}
	config.BurstSize = capacity

	return &LeakyBucketRateLimiter{
		config:    config,
		storage:   storage,
		leakyRate: leakyRate,
	}

}

type bucketLevel struct {
	Count      float64
	LastLeakAt int64
}

func (lb *LeakyBucketRateLimiter) Allow(key string) (bool, ratelimit.RateLimitInfo, error) {
	now := time.Now()
	bucketKey := fmt.Sprintf("leaky_bucket:%s", key)

	value, err := lb.storage.Get(bucketKey)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	var level bucketLevel
	if value == nil {
		level = bucketLevel{
			Count:      0,
			LastLeakAt: now.Unix(),
		}
	} else {
		level = value.(bucketLevel)

		elapsed := now.Unix() - level.LastLeakAt
		leaked := float64(elapsed) - lb.leakyRate
		level.Count -= leaked

		if level.Count < 0 {
			level.Count = 0
		}

		level.LastLeakAt = now.Unix()
	}

	allowed := level.Count < float64(lb.config.BurstSize)
	if allowed {
		level.Count += 1.0
	}

	err = lb.storage.Set(bucketKey, level, 0)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	remaining := max(lb.config.BurstSize-int(level.Count), 0)

	retryAfter := 0
	if !allowed {
		overflow := level.Count - float64(lb.config.BurstSize) + 1.0
		retryAfter = max(int(overflow/lb.leakyRate), 1)
	}

	secondsToEmpty := int(level.Count / lb.leakyRate)
	resetTime := now.Add(time.Duration(secondsToEmpty) * time.Second).Unix()

	info := ratelimit.RateLimitInfo{
		Limit:      lb.config.Limit,
		Remaining:  remaining,
		RetryAfter: retryAfter,
		ResetTime:  resetTime,
	}

	return allowed, info, nil

}

func (lb *LeakyBucketRateLimiter) Reset(key string) error {
	bucketKey := fmt.Sprintf("leaky_bucket:%s", key)
	return lb.storage.Delete(bucketKey)
}
