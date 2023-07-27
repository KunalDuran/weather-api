package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Middleware function to log requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Processing request")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the request header.
		token := r.Header.Get("Authorization")
		if token == "" {
			// No JWT token found in the request header.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse the JWT token.
		claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Check the signature of the token.
			return []byte("my-secret-key"), nil
		})

		if err != nil {
			// Invalid JWT token.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Get the user ID from the claims.
		id := claims.Claims.(jwt.MapClaims)["Subject"].(string)

		// Check if the user is authenticated.
		if id == "" {
			// User is not authenticated.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Continue with the request.
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware is a middleware function that adds the necessary CORS headers to the response.
func CorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers to allow cross-origin requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// If it's a preflight request, send an empty response with the necessary headers and return
		if r.Method == http.MethodOptions {
			return
		}

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
