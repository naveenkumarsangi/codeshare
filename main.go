package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	fb, err := getFirebaseApp()
	if err != nil {
		log.Fatalf("getFirebaseApp: %v", err)
	}

	bucket, err = getDefaultBucket(fb)
	if err != nil {
		log.Fatalf("getDefaultBucket: %v", err)
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Make paste
	r.POST("/paste", func(c *gin.Context) {
		content := c.PostForm("content")
		if content == "" {
			c.String(http.StatusBadRequest, "content in form is empty")
			return
		}

		id := uuid.New()
		if err := uploadFile(content, bucket, id.String()); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("unable to upload content, error: %s", err))
		} else {
			c.String(http.StatusOK, fmt.Sprintf("http://localhost:8080/s/%s", id))
		}
	})

	// Read shared content
	r.GET("/s/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		downloadFile(c.Writer, bucket, id)
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Println(err)
	}
	if storageClient != nil {
		storageClient.Close()
	}
}
