package redis

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func GetRedis() *redis.Client {
	addr, password := getRedisConnection()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}) //&redis.Options{} to get redis.Client pointer

	// test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}
	return RedisClient
}

func getRedisConnection() (string, string) {
	ENV := os.Getenv("ENV")
	if ENV == "local" {
		currentDir, _ := os.Getwd()
		log.Printf("Current working directory: %s", currentDir)
		dbenvPath := filepath.Join(currentDir, "../.dbenv")
		log.Printf("Current env file Path: %s", dbenvPath)

		_ = godotenv.Load(dbenvPath)
	}
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	addr := REDIS_HOST + ":" + REDIS_PORT
	password := REDIS_PASSWORD
	return addr, password
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
