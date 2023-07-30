package data

import (
	"database/sql"
	"time"

	"github.com/KunalDuran/weather-api/models"
	"github.com/KunalDuran/weather-api/util"
)

func InsertWeatherHistory(db *sql.DB, weather models.WeatherResponse, userID string) (int, error) {
	var insertedID int64
	stmt := "INSERT INTO weather_history (city_name, user_id, coord_lon, coord_lat, weather_id, weather_main, weather_description, weather_icon, base, temp, feels_like, temp_min, temp_max, pressure, humidity, visibility, wind_speed, wind_deg, clouds_all, dt, sys_type, sys_id, sys_country, sys_sunrise, sys_sunset, timezone) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := db.Exec(stmt,
		weather.Name,
		userID,
		weather.Coord.Lon,
		weather.Coord.Lat,
		weather.Weathers[0].ID,
		weather.Weathers[0].Main,
		weather.Weathers[0].Description,
		weather.Weathers[0].Icon,
		weather.Base,
		weather.Main.Temp,
		weather.Main.FeelsLike,
		weather.Main.TempMin,
		weather.Main.TempMax,
		weather.Main.Pressure,
		weather.Main.Humidity,
		weather.Visibility,
		weather.Wind.Speed,
		weather.Wind.Deg,
		weather.Clouds.All,
		weather.Dt,
		weather.Sys.Type,
		weather.Sys.ID,
		weather.Sys.Country,
		weather.Sys.Sunrise,
		weather.Sys.Sunset,
		weather.Timezone,
	)
	if err != nil {
		return 0, err
	}

	insertedID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(insertedID), nil
}

func DeleteWeather(db *sql.DB, id int) error {

	stmt := "DELETE FROM weather_history WHERE id = ?"

	_, err := db.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateWeatherIfExists(db *sql.DB, weather models.WeatherResponse, userID int) error {

	sqlStatement := `UPDATE weather_history SET
	  coord_lon = ?,
	  coord_lat = ?,
	  weather_id = ?,
	  weather_main = ?,
	  weather_description = ?,
	  weather_icon = ?,
	  base = ?,
	  temp = ?,
	  feels_like = ?,
	  temp_min = ?,
	  temp_max = ?,
	  pressure = ?,
	  humidity = ?,
	  visibility = ?,
	  wind_speed = ?,
	  wind_deg = ?,
	  clouds_all = ?,
	  dt = ?,
	  sys_type = ?,
	  sys_id = ?,
	  sys_country = ?,
	  sys_sunrise = ?,
	  sys_sunset = ?,
	  timezone = ?
	WHERE id = ?`

	_, err := db.Exec(sqlStatement,
		weather.Coord.Lon,
		weather.Coord.Lat,
		weather.Weathers[0].ID,
		weather.Weathers[0].Main,
		weather.Weathers[0].Description,
		weather.Weathers[0].Icon,
		weather.Base,
		weather.Main.Temp,
		weather.Main.FeelsLike,
		weather.Main.TempMin,
		weather.Main.TempMax,
		weather.Main.Pressure,
		weather.Main.Humidity,
		weather.Visibility,
		weather.Wind.Speed,
		weather.Wind.Deg,
		weather.Clouds.All,
		weather.Dt,
		weather.Sys.Type,
		weather.Sys.ID,
		weather.Sys.Country,
		weather.Sys.Sunrise,
		weather.Sys.Sunset,
		weather.Timezone,
		userID)

	if err != nil {
		return err
	}

	return nil
}

func FetchWeatherHistory(db *sql.DB, userID string) ([]models.WeatherResponse, error) {

	var weatherHistory []models.WeatherResponse

	stmt := "SELECT * FROM weather_history where user_id = ?"

	rows, err := db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		weather := models.WeatherResponse{}
		cityWeather := &models.Weather{}

		var createdAt string
		err := rows.Scan(
			&weather.WeatherID,
			&weather.Name,
			&weather.UserID,
			&weather.Coord.Lon,
			&weather.Coord.Lat,
			&cityWeather.ID,
			&cityWeather.Main,
			&cityWeather.Description,
			&cityWeather.Icon,
			&weather.Base,
			&weather.Main.Temp,
			&weather.Main.FeelsLike,
			&weather.Main.TempMin,
			&weather.Main.TempMax,
			&weather.Main.Pressure,
			&weather.Main.Humidity,
			&weather.Visibility,
			&weather.Wind.Speed,
			&weather.Wind.Deg,
			&weather.Clouds.All,
			&weather.Dt,
			&weather.Sys.Type,
			&weather.Sys.ID,
			&weather.Sys.Country,
			&weather.Sys.Sunrise,
			&weather.Sys.Sunset,
			&weather.Timezone,
			&createdAt)
		if err != nil {
			return nil, err
		}

		weather.CreatedAt, _ = util.ParseTimestamp(createdAt)
		weatherHistory = append(weatherHistory, weather)
	}

	return weatherHistory, nil
}

func CreateUser(db *sql.DB, username string, password string, birthDate time.Time) (int, error) {
	stmt := "INSERT INTO users (username, password, date_of_birth) VALUES (?, ?, ?)"

	result, err := db.Exec(stmt, username, password, birthDate)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {

	stmt := "SELECT * FROM users WHERE username = ?"

	rows := db.QueryRow(stmt, username)

	user := &models.User{}

	var birthDate, createdAt string
	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&birthDate,
		&createdAt,
	)

	user.DateOfBirth, _ = util.ParseDOB(birthDate)
	user.CreatedAt, _ = util.ParseTimestamp(createdAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func BulkDeleteWeathers(db *sql.DB, userID string) (int, error) {
	stmt := "DELETE FROM weather_history WHERE user_id = ?"
	result, err := db.Exec(stmt, userID)
	if err != nil {
		return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(affectedRows), nil
}

func GetWeatherByID(db *sql.DB, id int) (*models.WeatherResponse, error) {

	stmt := "SELECT * FROM weather_history WHERE id = ?"

	rows := db.QueryRow(stmt, id)

	weather := &models.WeatherResponse{}
	cityWeather := &models.Weather{}

	var createdAt string
	err := rows.Scan(
		&weather.WeatherID,
		&weather.Name,
		&weather.UserID,
		&weather.Coord.Lon,
		&weather.Coord.Lat,
		&cityWeather.ID,
		&cityWeather.Main,
		&cityWeather.Description,
		&cityWeather.Icon,
		&weather.Base,
		&weather.Main.Temp,
		&weather.Main.FeelsLike,
		&weather.Main.TempMin,
		&weather.Main.TempMax,
		&weather.Main.Pressure,
		&weather.Main.Humidity,
		&weather.Visibility,
		&weather.Wind.Speed,
		&weather.Wind.Deg,
		&weather.Clouds.All,
		&weather.Dt,
		&weather.Sys.Type,
		&weather.Sys.ID,
		&weather.Sys.Country,
		&weather.Sys.Sunrise,
		&weather.Sys.Sunset,
		&weather.Timezone,
		&createdAt,
	)

	weather.CreatedAt, _ = util.ParseTimestamp(createdAt)
	weather.Weathers = append(weather.Weathers, *cityWeather)

	if err != nil {
		return nil, err
	}

	return weather, nil
}

func UpdateWeather(db *sql.DB, weather *models.WeatherResponse) error {

	stmt := "UPDATE weather_history SET city_name = ?, coord_lon = ?, coord_lat = ?, weather_id = ?, weather_main = ?, weather_description = ?, weather_icon = ?, base = ?, temp = ?, feels_like = ?, temp_min = ?, temp_max = ?, pressure = ?, humidity = ?, visibility = ?, wind_speed = ?, wind_deg = ?, clouds_all = ?, dt = ?, sys_type = ?, sys_id = ?, sys_country = ?, sys_sunrise = ?, sys_sunset = ?, timezone = ? WHERE id = ?"

	_, err := db.Exec(stmt,
		weather.Name,
		weather.Coord.Lon,
		weather.Coord.Lat,
		weather.Weathers[0].ID,
		weather.Weathers[0].Main,
		weather.Weathers[0].Description,
		weather.Weathers[0].Icon,
		weather.Base,
		weather.Main.Temp,
		weather.Main.FeelsLike,
		weather.Main.TempMin,
		weather.Main.TempMax,
		weather.Main.Pressure,
		weather.Main.Humidity,
		weather.Visibility,
		weather.Wind.Speed,
		weather.Wind.Deg,
		weather.Clouds.All,
		weather.Dt,
		weather.Sys.Type,
		weather.Sys.ID,
		weather.Sys.Country,
		weather.Sys.Sunrise,
		weather.Sys.Sunset,
		weather.Timezone,
		weather.WeatherID,
	)
	if err != nil {
		return err
	}

	return nil
}
