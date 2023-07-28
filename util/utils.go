package util

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/KunalDuran/weather-api/models"
	"github.com/golang-jwt/jwt/v4"
)

func WebRequest(method string, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func JSONResponse(w http.ResponseWriter, statusCode int, resp *models.Response) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(resp)

}

func CreateToken(id int, username string) (string, error) {
	claims := jwt.MapClaims{
		"Issuer":    "my-app",
		"Subject":   strconv.Itoa(id),
		"Username":  username,
		"ExpiresAt": time.Now().Add(time.Hour * 24),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte("my-secret-key")
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetUserIDFromToken(token string) (string, error) {
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("my-secret-key"), nil
	})
	if err != nil {
		return "", err
	}

	id := claims.Claims.(jwt.MapClaims)["Subject"].(string)

	return id, nil
}

func ParseDOB(dob string) (time.Time, error) {
	birthDate, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return time.Time{}, err
	}
	return birthDate, nil
}

func ParseTimestamp(timestamp string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ValidateEmail(email string) bool {
	// Regular expression for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return false
	}
	return matched
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	// Password should have at least one uppercase letter, one lowercase letter, and one digit
	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		if unicode.IsUpper(ch) {
			hasUpper = true
		} else if unicode.IsLower(ch) {
			hasLower = true
		} else if unicode.IsDigit(ch) {
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func ValidateDateOfBirth(dateOfBirth string) bool {
	// Define the expected date format
	dateLayout := "2006-01-02"

	// Parse the date string into a time.Time value
	dob, err := time.Parse(dateLayout, dateOfBirth)
	if err != nil {
		return false
	}

	// Check if the date is in the past (to avoid future dates of birth)
	now := time.Now()
	return dob.Before(now)
}
