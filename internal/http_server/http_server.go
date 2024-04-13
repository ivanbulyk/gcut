package http_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ivanbulyk/gcut/internal/http_server/routes"
	"github.com/ivanbulyk/gcut/internal/storage"
	"log"
	"os"
)

func Init() {
	pingResult, err := storage.CreateRedisClient(0).Ping(storage.Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to create redis client: %v", err)
	}

	fmt.Println("PING:", pingResult)

	router := gin.Default()
	router.Use(gin.Logger())

	SetUpRotes(router)

	router.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))

}

func SetUpRotes(router *gin.Engine) {
	router.GET("/:url", routes.ResolveURL)
	router.POST("/api/v1", routes.ShortenURL)
	router.GET("/", routes.IndexHandler)
}
