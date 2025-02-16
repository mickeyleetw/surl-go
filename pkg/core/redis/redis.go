package redis

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
)

// RedisClient is a pointer to a redis.Client
var RedisClient *redis.Client

// GetRedis returns a singleton instance of the redis client
func GetRedis() *redis.Client {
	addr, password := getRedisConnection()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}) // &redis.Options{} to get redis.Client pointer

	// test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}
	return RedisClient
}

func getRedisConnection() (string, string) {
	env := os.Getenv("ENV")
	if env == "local" {
		currentDir, _ := os.Getwd()
		log.Printf("Current working directory: %s", currentDir)
		dbenvPath := filepath.Join(currentDir, "../.dbenv")
		log.Printf("Current env file Path: %s", dbenvPath)

		_ = godotenv.Load(dbenvPath)
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	addr := redisHost + ":" + redisPort
	password := redisPassword
	return addr, password
}

// CloseRedis is a function that closes the Redis client connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
