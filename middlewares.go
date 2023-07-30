package main

import (
	"net/http"

	"github.com/KunalDuran/weather-api/models"
	"github.com/KunalDuran/weather-api/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

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
		token := r.Header.Get("Authorization")
		if token == "" {
			util.JSONResponse(w, http.StatusUnauthorized, &models.Response{
				Status:  "error",
				Message: "No token provided",
				Data:    nil,
			})
			return
		}

		claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("my-secret-key"), nil
		})

		if err != nil {
			util.JSONResponse(w, http.StatusUnauthorized, &models.Response{
				Status:  "error",
				Message: "Invalid token",
				Data:    nil,
			})
			return
		}

		id := claims.Claims.(jwt.MapClaims)["Subject"].(string)

		if id == "" {
			util.JSONResponse(w, http.StatusUnauthorized, &models.Response{
				Status:  "error",
				Message: "Invalid token",
				Data:    nil,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware is a middleware function that adds the necessary CORS headers to the response.
func CorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// If it's a preflight request, send an empty response with the necessary headers and return
		if r.Method == http.MethodOptions {
			return
		}

		handler.ServeHTTP(w, r)
	})
}
