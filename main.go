package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KunalDuran/weather-api/data"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
var API_KEY string

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	API_KEY = os.Getenv("API_KEY")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db, err = data.InitDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Warn(err)
		fmt.Println("Error connecting to database")
		return
	}

	// to keep the connection alive
	go func() {
		for {
			time.Sleep(time.Second * 15)
			err := db.Ping()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/weather", AuthMiddleware(weatherHandler))
	http.HandleFunc("/api/history", AuthMiddleware(getWeatherHistoryHandler))
	http.HandleFunc("/api/history/delete", AuthMiddleware(deleteWeatherHistoryHandler))
	http.HandleFunc("/api/history/bulkdelete", AuthMiddleware(bulkDeleteWeatherHistoryHandler))

	log.Println("Server started on port 8080")

	mux := http.DefaultServeMux

	loggedMux := loggingMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8080", CorsMiddleware(loggedMux)))
}
