package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisClient "shorten_url/pkg/core/redis"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage() (*RedisStorage, error) {
	client := redisClient.GetRedis()

	// testing redis connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

func (rs *RedisStorage) Get(key string) (interface{}, error) {
	val, err := rs.client.Get(rs.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("redis Get error: %w", err)
	}

	var intVal int64
	_, err = fmt.Sscanf(val, "%d", &intVal)
	if err == nil {
		return intVal, nil
	}

	var jsonVal interface{}
	err = json.Unmarshal([]byte(val), &jsonVal)
	if err == nil {
		return jsonVal, nil
	}

	return val, nil

}

func (rs *RedisStorage) Set(key string, value interface{}, ttl time.Duration) error {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case int, int64, float64:
		strValue = fmt.Sprintf("%v", v)
	default:
		// 複雜類型使用 JSON 序列化
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("json marshal error: %w", err)
		}
		strValue = string(jsonBytes)
	}
	err := rs.client.Set(rs.ctx, key, strValue, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis Set error: %w", err)
	}

	return nil
}

func (rs *RedisStorage) Increment(key string, ttl time.Duration) (int64, error) {
	newValue, err := rs.client.Incr(rs.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis Incr error: %w", err)
	}

	if newValue == 1 && ttl > 0 {
		err = rs.client.Expire(rs.ctx, key, ttl).Err()
		if err != nil {
			return newValue, fmt.Errorf("redis Expire error: %w", err)
		}
	}

	return newValue, nil

}

func (rs *RedisStorage) Delete(key string) error {
	err := rs.client.Del(rs.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis Del error: %w", err)
	}

	return nil
}
