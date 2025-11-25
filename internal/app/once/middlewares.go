package once

import (
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var now = time.Now()
		log.Printf(
			`time=%s path=%s method=%s body_bytes_received=%v`,
			time.Now().Format(time.RFC3339), r.URL.Path, r.Method, r.ContentLength,
		)

		next.ServeHTTP(w, r)

		var requestDuration = time.Since(now)
		log.Printf(
			`time=%s path=%s method=%s request_duration=%v content_type=%s`,
			time.Now().Format(time.RFC3339), r.URL.Path, r.Method, requestDuration, r.Header.Get("Content-Type"),
		)
	})
}
