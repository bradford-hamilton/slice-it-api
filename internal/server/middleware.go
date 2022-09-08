package server

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/go-chi/cors"
)

func limiterMiddleware() *limiter.Limiter {
	lmt := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute * 30,
	})
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	})
	return lmt
}

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300, // Maximum value not ignored by any major browsers
	})
}
