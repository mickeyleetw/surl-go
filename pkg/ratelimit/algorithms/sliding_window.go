package algorithms

import (
	"fmt"
	"shorten_url/pkg/ratelimit"
	"time"
)

type SlidingWindowRateLimiter struct {
	config  ratelimit.Config
	storage ratelimit.Storage
}

func NewSlidingWindowRateLimiter(storage ratelimit.Storage, config ratelimit.Config) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		config:  config,
		storage: storage,
	}
}

type requestLog struct {
	TimeStamps []int64
}

func (sw *SlidingWindowRateLimiter) Allow(key string) (bool, ratelimit.RateLimitInfo, error) {
	now := time.Now()
	windowStart := now.Add(-sw.config.Window).Unix()
	logkey := fmt.Sprintf("sliding_window:%s", key)

	value, err := sw.storage.Get(logkey)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	var log requestLog
	if value != nil {
		log = value.(requestLog)
	} else {
		log = requestLog{TimeStamps: []int64{}}
	}

	validTimestamps := []int64{}
	for _, ts := range log.TimeStamps {
		if ts > windowStart {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	currentCount := len(validTimestamps)
	allowed := currentCount < sw.config.Limit
	if allowed {
		validTimestamps = append(validTimestamps, now.Unix())
	}

	log.TimeStamps = validTimestamps
	err = sw.storage.Set(logkey, log, sw.config.Window*2)
	if err != nil {
		return false, ratelimit.RateLimitInfo{}, err
	}

	remaining := max(sw.config.Limit-len(validTimestamps), 0)

	resetTime := now.Add(sw.config.Window).Unix()
	if len(validTimestamps) > 0 {
		oldTimeStamp := validTimestamps[0]
		resetTime = time.Unix(oldTimeStamp, 0).Add(sw.config.Window).Unix()
	}

	retryAfter := 0
	if !allowed && len(validTimestamps) > 0 {
		oldTimeStamp := validTimestamps[0]
		retryAfter = max(int(time.Unix(oldTimeStamp, 0).Add(sw.config.Window).Sub(now).Seconds()), 0)
	}

	info := ratelimit.RateLimitInfo{
		Limit:      sw.config.Limit,
		Remaining:  remaining,
		ResetTime:  resetTime,
		RetryAfter: retryAfter,
	}

	return allowed, info, nil
}

func (sw *SlidingWindowRateLimiter) Reset(key string) error {
	logKey := fmt.Sprintf("sliding_window:%s", key)
	return sw.storage.Delete(logKey)
}
