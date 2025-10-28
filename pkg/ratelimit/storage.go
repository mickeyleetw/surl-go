package ratelimit

import "time"

// Storage defines the storage backend interface for rate limit data
// Supports both in-memory storage and Redis storage implementations
type Storage interface {
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{}, ttl time.Duration) error
	Increment(key string, ttl time.Duration) (int64, error)
	Delete(key string) error
	Close() error
}
