package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Log the method and the requested URL
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		}).Info("Incoming request")
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
		// Log how long it took
		log.WithFields(log.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"ip":       r.RemoteAddr,
			"duration": time.Since(start),
		}).Info("Completed request")
	})
}
