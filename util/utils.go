package util

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/KunalDuran/weather-api/models"
	"github.com/golang-jwt/jwt/v4"
)

// WebRequest function is used to make an HTTP request
// to the specified URL and return the response body.
// It returns an error if the request fails.
// The caller is responsible for closing the response body.
// Usage:
// resp, err := WebRequest("GET", "https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}", nil)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// defer resp.Body.Close()
// body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// fmt.Println(string(body))
func WebRequest(method string, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// JSONResponse function is used to send a JSON response
// to the client.
// Usage:
// JSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
func JSONResponse(w http.ResponseWriter, statusCode int, resp *models.Response) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(resp)

}

func CreateToken(id int, username string) (string, error) {
	// Create a JWT claims object.
	claims := jwt.MapClaims{
		"Issuer":    "my-app",
		"Subject":   strconv.Itoa(id),
		"Username":  username,
		"ExpiresAt": time.Now().Add(time.Hour * 24),
	}

	// Create a JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	secret := []byte("my-secret-key")
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetUserIDFromToken(token string) (string, error) {

	// Parse the JWT token.
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check the signature of the token.
		return []byte("my-secret-key"), nil
	})
	if err != nil {
		// Invalid JWT token.
		return "", err
	}

	// Get the user ID from the claims.
	id := claims.Claims.(jwt.MapClaims)["Subject"].(string)

	// Return the user name.
	return id, nil
}

// Parse DOB
func ParseDOB(dob string) (time.Time, error) {
	// Parse the date of birth.
	birthDate, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return time.Time{}, err
	}
	return birthDate, nil
}

// ParseTimestamp
func ParseTimestamp(timestamp string) (time.Time, error) {
	// Parse the timestamp.
	t, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
