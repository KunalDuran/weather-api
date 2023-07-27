package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/KunalDuran/weather-api/data"
	"github.com/KunalDuran/weather-api/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
var API_KEY string

type App struct {
	Logger *log.Logger
	DB     *sql.DB
}

func main() {

	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	API_KEY = os.Getenv("API_KEY")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Connect the database
	db, err = data.InitDB(dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		log.Panic(err)
	}

	// Define API routes
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/weather", util.AuthMiddleware(weatherHandler))
	http.HandleFunc("/api/history", util.AuthMiddleware(getWeatherHistoryHandler))
	http.HandleFunc("/api/history/update", util.AuthMiddleware(updateWeatherHistoryHandler))
	http.HandleFunc("/api/history/delete", util.AuthMiddleware(deleteWeatherHistoryHandler))
	http.HandleFunc("/api/history/bulkdelete", util.AuthMiddleware(bulkDeleteWeatherHistoryHandler))

	// Start the server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", util.CorsMiddleware(http.DefaultServeMux)))
}
