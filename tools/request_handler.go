package tools

import (
	"log"
	"net/http"
)

// ValidadeRequestHeader should validate the Bearer head
func ValidadeRequestHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Bearer")

		if token != "" {
			log.Println("Bearer token was found:", token)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden, cause: Bearer token was not found", http.StatusForbidden)
		}

	})
}

// LoggingRequest logging all requests
func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, " - ", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
