package apps

import (
	"context"
	"log"
	"net/http"
	"shorten_url/pkg/core/redis"
	"shorten_url/pkg/models"

	repositories "shorten_url/pkg/repositories"

	"github.com/gin-gonic/gin"
)

// CreateURLShorten is a function that creates a short URL
func CreateURLShorten(c *gin.Context) {
	// wrap json to validated payload model
	createURLShortenModel := c.MustGet("validated_data").(models.CreateURLShortenModel)
	shortURL, err := repositories.GetOrCreateShortURL(createURLShortenModel.LongURL)
	if err != nil {
		log.Println("Error getting or creating short url: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// set the short url to redis
	redisClient := redis.GetRedis()

	result, _ := redisClient.Get(context.Background(), shortURL).Result()
	if result == "" {
		redisErr := redisClient.Set(context.Background(), shortURL, createURLShortenModel.LongURL, 0).Err()
		if redisErr != nil {
			log.Println("Error setting short url to redis: ", redisErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": redisErr.Error()})
			return
		}
		log.Println("Short url set to redis: ", shortURL)
	}

	// if the original url is already shortened, return the short url
	if shortURL != "" {
		c.JSON(http.StatusCreated, gin.H{"short_url": shortURL})
	}
}

// GetURLByShorten is a function that gets a long URL from a short URL
func GetURLByShorten(c *gin.Context, shortURL string) {
	// check if the short url is already exist in redis
	result, redisErr := redis.GetRedis().Get(context.Background(), shortURL).Result()
	if redisErr != nil {
		log.Println("Error getting short url from redis: ", redisErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": redisErr.Error()})
		return
	}
	if result == "" {
		url, err := repositories.GetLongURLByShortURL(shortURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, url)
	} else {
		log.Println("Short url found in redis: ", result)
		c.Redirect(http.StatusTemporaryRedirect, result)
	}
}
