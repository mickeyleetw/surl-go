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

func CreateUrlShorten(c *gin.Context) {
	// wrap json to validated payload model
	createUrlShortenModel := c.MustGet("validated_data").(models.CreateUrlShortenModel)
	shortUrl, err := repositories.GetOrCreateShortUrl(createUrlShortenModel.LongUrl)
	if err != nil {
		log.Println("Error getting or creating short url: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// set the short url to redis
	redis_client := redis.GetRedis()

	result, _ := redis_client.Get(context.Background(), shortUrl).Result()
	if result == "" {
		redis_err := redis_client.Set(context.Background(), shortUrl, createUrlShortenModel.LongUrl, 0).Err()
		if redis_err != nil {
			log.Println("Error setting short url to redis: ", redis_err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": redis_err.Error()})
			return
		}
		log.Println("Short url set to redis: ", shortUrl)
	}

	// if the original url is already shortened, return the short url
	if shortUrl != "" {
		c.JSON(http.StatusCreated, gin.H{"short_url": shortUrl})
	}

}

func GetUrlByShorten(c *gin.Context, shortUrl string) {
	// check if the short url is already exist in redis
	result, redis_err := redis.GetRedis().Get(context.Background(), shortUrl).Result()
	if redis_err != nil {
		log.Println("Error getting short url from redis: ", redis_err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": redis_err.Error()})
		return
	} else {
		if result == "" {
			url, err := repositories.GetUrlByShortUrl(shortUrl)
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
}
