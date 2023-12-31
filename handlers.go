package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/KunalDuran/weather-api/data"
	"github.com/KunalDuran/weather-api/models"
	"github.com/KunalDuran/weather-api/util"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		util.JSONResponse(w, http.StatusMethodNotAllowed, &models.Response{
			Status:  "error",
			Message: "Method not allowed.",
			Data:    nil,
		})
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid JSON provided.",
			Data:    nil,
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Username and password are required.",
			Data:    nil,
		})
		return
	}

	if !util.ValidateEmail(user.Username) {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid username, please provide a valid email address.",
			Data:    nil,
		})
		return
	}

	userRecord, err := data.GetUserByUsername(db, user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			util.JSONResponse(w, http.StatusUnauthorized, &models.Response{
				Status:  "error",
				Message: "Invalid credentials.",
				Data:    nil,
			})
			return
		} else {
			util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
				Status:  "error",
				Message: "Internal server error.",
				Data:    nil,
			})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRecord.Password), []byte(user.Password)); err != nil {
		util.JSONResponse(w, http.StatusUnauthorized, &models.Response{
			Status:  "error",
			Message: "Invalid credentials.",
			Data:    nil,
		})
		return
	}

	token, err := util.CreateToken(userRecord.ID, userRecord.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			Status:  "error",
			Message: "Failed to create token.",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	resp := &models.Response{
		Status:  "success",
		Message: "Logged in successfully.",
		Data:    map[string]string{"token": token},
	}
	util.JSONResponse(w, http.StatusOK, resp)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		util.JSONResponse(w, http.StatusMethodNotAllowed, &models.Response{
			Status:  "error",
			Message: "Method not allowed.",
			Data:    nil,
		})
		return
	}

	var user struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		BirthDate string `json:"birth_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid JSON provided.",
			Data:    nil,
		})
		return
	}

	if user.Username == "" || user.Password == "" || user.BirthDate == "" {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Username, password and birth date are required.",
			Data:    nil,
		})
		return
	}

	if !util.ValidateEmail(user.Username) {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid email address.",
			Data:    nil,
		})
		return
	}

	if !util.ValidatePassword(user.Password) {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one digit.",
			Data:    nil,
		})
		return
	}

	if !util.ValidateDateOfBirth(user.BirthDate) {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid birth date.",
			Data:    nil,
		})
		return
	}

	existingUser, err := data.GetUserByUsername(db, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
			Status:  "error",
			Message: "Internal server error.",
			Data:    nil,
		})
		return
	} else if existingUser != nil {
		util.JSONResponse(w, http.StatusConflict, &models.Response{
			Status:  "error",
			Message: "Username already exists.",
			Data:    nil,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
			Status:  "error",
			Message: "Internal server error.",
			Data:    nil,
		})
		return
	}

	parseBirthDate, err := time.Parse("2006-01-02", user.BirthDate)
	if err != nil {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid birth date.",
			Data:    nil,
		})
		return
	}

	id, err := data.CreateUser(db, user.Username, string(hashedPassword), parseBirthDate)
	if err != nil {
		util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
			Status:  "error",
			Message: "Internal server error.",
			Data:    nil,
		})
		return
	}

	token, err := util.CreateToken(id, user.Username)
	if err != nil {
		util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
			Status:  "error",
			Message: "Failed to create token.",
			Data:    nil,
		})
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	resp := &models.Response{
		Status:  "success",
		Message: "Registered successfully.",
		Data:    map[string]string{"token": token},
	}

	util.JSONResponse(w, http.StatusOK, resp)

}

func weatherHandler(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")
	if city == "" {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "City name is required.",
			Data:    nil,
		})
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, API_KEY)

	resp, err := util.WebRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp models.StandardResponse
		err = json.Unmarshal(body, &errorResp)
		if err != nil {
			log.Error(err)
		}

		util.JSONResponse(w, http.StatusOK, &models.Response{
			Status:  "error",
			Message: "City not found.",
			Data:    nil,
		})
		return
	}

	var weatherResponse models.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Error(err)
	}

	userID, err := util.GetUserIDFromToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Error(err)
	}

	insertedRowID, err := data.InsertWeatherHistory(db, weatherResponse, userID)
	if err != nil {
		log.Error(err)
	}

	weatherResponse.WeatherID = insertedRowID

	util.JSONResponse(w, http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Weather fetched successfully",
		Data:    weatherResponse,
	})

}

func getWeatherHistoryHandler(w http.ResponseWriter, r *http.Request) {

	userID, _ := util.GetUserIDFromToken(r.Header.Get("Authorization"))
	weatherData, err := data.FetchWeatherHistory(db, userID)
	if err != nil {
		util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
			Status:  "error",
			Message: "Failed to fetch weather history.",
			Data:    nil,
		})
		return
	}

	if weatherData == nil {
		util.JSONResponse(w, http.StatusOK, &models.Response{
			Status:  "info",
			Message: "No Search History Found.",
			Data:    nil,
		})
		return
	}
	util.JSONResponse(w, http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Search history fetched successfully.",
		Data:    weatherData,
	})
}

func deleteWeatherHistoryHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		util.JSONResponse(w, http.StatusMethodNotAllowed, &models.Response{
			Status:  "error",
			Message: "Method not allowed.",
			Data:    nil,
		})
		return
	}

	weatherID := r.URL.Query().Get("weatherID")

	weatherIDInt, err := strconv.Atoi(weatherID)
	if err != nil || weatherIDInt <= 0 || weatherID == "" {
		util.JSONResponse(w, http.StatusBadRequest, &models.Response{
			Status:  "error",
			Message: "Invalid weatherID.",
			Data:    nil,
		})
		return
	}

	affectedRows, err := data.DeleteWeather(db, weatherIDInt)
	if err != nil {
		util.JSONResponse(w, http.StatusInternalServerError, &models.Response{
			Status:  "error",
			Message: "Failed to delete weather.",
			Data:    nil,
		})
		return
	}

	if affectedRows == 0 {
		util.JSONResponse(w, http.StatusNotFound, &models.Response{
			Status:  "error",
			Message: "Weather not found with this ID.",
			Data:    nil,
		})
		return
	}

	util.JSONResponse(w, http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Successfully deleted weather",
		Data:    nil,
	})
}

func bulkDeleteWeatherHistoryHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		util.JSONResponse(w, http.StatusMethodNotAllowed, &models.Response{
			Status:  "error",
			Message: "Method not allowed.",
			Data:    nil,
		})
		return
	}

	userID, _ := util.GetUserIDFromToken(r.Header.Get("Authorization"))
	affectedRows, err := data.BulkDeleteWeathers(db, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := &models.Response{
			Status:  "error",
			Message: "Failed to delete weathers.",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if affectedRows > 0 {
		resp := &models.Response{
			Status:  "success",
			Message: "Successfully deleted weathers.",
			Data:    nil,
		}
		util.JSONResponse(w, http.StatusOK, resp)
	} else {
		resp := &models.Response{
			Status:  "info",
			Message: "No history to delete.",
			Data:    nil,
		}
		util.JSONResponse(w, http.StatusNoContent, resp)
	}
}
