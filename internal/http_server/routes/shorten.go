package routes

import (
	"errors"
	"github.com/ivanbulyk/gcut/internal/lib/utils"
	"github.com/ivanbulyk/gcut/internal/storage"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"net/http"
	"strconv"
	"time"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL shortens url
func ShortenURL(ctx *gin.Context) {

	// parsing incoming request body
	body := new(Request)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot parse JSON",
			"error":   err.Error(),
		})
		return
	}

	// implement rate limiting
	// everytime a user queries, check if the IP is already in database,
	// if yes, decrement the calls remaining by one, else add the IP to database
	// with expiry of `30mins`. So in this case the user will be able to send 10
	// requests every 30 minutes
	r2 := storage.CreateRedisClient(1)
	defer func(r2 *redis.Client) {
		err := r2.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "couldn't close redis client",
			})
		}
	}(r2)
	value, err := r2.Get(storage.Ctx, ctx.RemoteIP()).Result()
	if errors.Is(err, redis.Nil) {
		_ = r2.Set(storage.Ctx, ctx.RemoteIP(), os.Getenv("RATE_LIMIT"), 30*60*time.Second).Err()
	} else {
		value, _ = r2.Get(storage.Ctx, ctx.RemoteIP()).Result()
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := r2.TTL(storage.Ctx, ctx.RemoteIP()).Result()
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
			return

		}
	}

	// check if an input is valid url
	if !govalidator.IsURL(body.URL) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid URL",
		})
		return
	}

	// enforce http
	// all url will be converted to https before storing in database
	body.URL = utils.EnforceHTTP(body.URL)

	// check for domain error
	if !utils.RemoveDomainError(body.URL) {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "invalid input",
		})
		return
	}

	// check if the user has provided any custom short urls
	// if yes, proceed,
	// else, create a new short using the first 6 digits of uuid
	slug := ""
	if body.CustomShort == "" {
		slug = uuid.New().String()[:6]
	} else {
		slug = body.CustomShort
	}

	r := storage.CreateRedisClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "couldn't close redis client",
			})
		}
	}(r)

	value, _ = r.Get(storage.Ctx, slug).Result()

	// check if the short provided by user  is already in use
	if value != "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "URL short already in use",
		})
		return
	}
	if body.Expiry == 0 {
		body.Expiry = 24 // default expiry of 24 hours
	}
	err = r.Set(storage.Ctx, slug, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to connect to server",
		})
		return
	}

	// respond with the url, short, expiry in hours, calls remaining and time to reset
	resp := Response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	r2.Decr(storage.Ctx, ctx.RemoteIP())
	value, _ = r2.Get(storage.Ctx, ctx.RemoteIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(value)
	ttl, _ := r2.TTL(storage.Ctx, ctx.RemoteIP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + slug
	ctx.JSON(http.StatusOK, gin.H{
		"payload": resp,
	})
	return
}
