package main

import (
	"log"
	"os"

	"booking/internal/adapters/redis"
	"booking/internal/booking"

	"github.com/gin-gonic/gin"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379" // фолбэк для локальной разработки
	}

	rdb := redis.NewClient(redisURL)
	store := booking.NewRedisStore(rdb)
	svc := booking.NewService(store)

	r := gin.Default()
	h := booking.NewHandler(svc)

	r.GET("/movies", h.ListMovies)
	r.GET("/movies/:movieID/seats", h.ListSeats)
	r.POST("/movies/:movieID/seats/:seatID/hold", h.HoldSeat)
	r.PUT("/sessions/:sessionID/confirm", h.ConfirmSession)
	r.DELETE("/sessions/:sessionID", h.ReleaseSession)
	
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
