package models

import "time"

// User represents the user data
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
}

// WeatherResponse represents the weather data received from the OpenWeatherMap API
type WeatherResponse struct {
	WeatherID int    `json:"weather_id"`
	UserID    string `json:"-"`
	Coord     struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weathers []Weather `json:"weather"`
	Base     string    `json:"base"`
	Main     struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone  int       `json:"timezone"`
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Cod       int       `json:"cod"`
	CreatedAt time.Time `json:"created_at"`
}

// StandardResponse represents the standard response from the OpenWeatherMap API
type StandardResponse struct {
	COD     string `json:"cod"`
	Message string `json:"message"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
