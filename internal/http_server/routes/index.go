package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// IndexHandler serves as healthcheck
func IndexHandler(ctx *gin.Context) {

	log.Println("Loading gCut Service...")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully loaded gCut Service!",
		"port":    os.Getenv("PORT"),
	})
}
