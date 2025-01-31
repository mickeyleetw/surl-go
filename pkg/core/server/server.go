package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	apps "shorten_url/pkg/apps"
	models "shorten_url/pkg/models"
)

var (
	app         *gin.Engine
	API_VERSION string
	BASE_PATH   string
	apiGroup    *gin.RouterGroup
)

func InitServer() {
	ENV := os.Getenv("ENV")
	if ENV == "local" {
		currentDir, _ := os.Getwd()
		serverEnvPath := filepath.Join(currentDir, "../.serverEnv")

		_ = godotenv.Load(serverEnvPath)
	}

	API_HOST := os.Getenv("API_HOST")
	API_PORT := os.Getenv("API_PORT")
	API_VERSION = os.Getenv("API_VERSION")

	BASE_PATH = "/" + API_VERSION

	env := os.Getenv("env")
	if env == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app = gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	apiGroup = app.Group(BASE_PATH)
	{
		// root route
		apiGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
		})

		apiGroup.POST("/shorten", ValidateRequest[models.CreateUrlShortenModel](), func(c *gin.Context) {
			apps.CreateUrlShorten(c)
		})

		apiGroup.GET("/:short_url", func(c *gin.Context) {
			shortUrl := c.Param("short_url")
			apps.GetUrlByShorten(c, shortUrl)
		})
	}

	app.Run(fmt.Sprintf("%s:%s", API_HOST, API_PORT))
}
