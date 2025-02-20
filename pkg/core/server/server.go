package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	apps "shorten_url/pkg/apps"
	models "shorten_url/pkg/models"
)

var (
	app        *gin.Engine
	apiVersion string
	basePath   string
	apiGroup   *gin.RouterGroup
)

// InitServer is a function that initializes the server
func InitServer() {
	ENV := os.Getenv("ENV")
	if ENV == "local" {
		currentDir, _ := os.Getwd()
		serverEnvPath := filepath.Join(currentDir, "../.serverEnv")

		_ = godotenv.Load(serverEnvPath)
	}

	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")
	apiVersion = os.Getenv("API_VERSION")

	basePath = "/" + apiVersion

	env := os.Getenv("env")
	if env == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app = gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	apiGroup = app.Group(basePath)
	{
		// root route
		apiGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
		})

		apiGroup.POST("/shorten", ValidateRequest[models.CreateURLShortenModel](), func(c *gin.Context) {
			apps.CreateURLShorten(c)
		})

		apiGroup.GET("/:short_url", func(c *gin.Context) {
			shortURL := c.Param("short_url")
			apps.GetURLByShorten(c, shortURL)
		})
	}

	log.Printf("Server starting at %s:%s", apiHost, apiPort)
	if err := app.Run(fmt.Sprintf("%s:%s", apiHost, apiPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
