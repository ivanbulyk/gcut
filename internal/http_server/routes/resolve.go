package routes

import (
	"errors"
	"github.com/ivanbulyk/gcut/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

// ResolveURL resolves url
func ResolveURL(ctx *gin.Context) {
	// get the short from the url
	url := ctx.Param("url")

	// query the redis db to find the original URL, if a match is found
	// increment the redirect counter and redirect to the original URL
	// else return error message
	r := storage.CreateRedisClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "couldn't close redis client",
			})
		}
	}(r)

	value, err := r.Get(storage.Ctx, url).Result()
	if errors.Is(err, redis.Nil) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "short not found in redis database",
		})
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't connect to redis database",
		})
	}

	// increment the counter
	rIncr := storage.CreateRedisClient(1)
	defer func(rIncr *redis.Client) {
		err := rIncr.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "couldn't close redis client",
			})
		}
	}(rIncr)
	_ = rIncr.Incr(storage.Ctx, "counter")

	// redirect to original URL
	ctx.Redirect(http.StatusMovedPermanently, value)
}
