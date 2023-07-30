package util

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTokenAndParseUserID(t *testing.T) {
	// Test data
	id := 123
	username := "testuser"

	// Create a token
	token, err := CreateToken(id, username)
	assert.NoError(t, err)

	// Parse the user ID from the token
	parsedID, err := GetUserIDFromToken(token)
	assert.NoError(t, err)

	// Ensure the parsed user ID matches the original ID
	assert.Equal(t, strconv.Itoa(id), parsedID)
}

func TestParseDOB(t *testing.T) {
	// Test data
	dobString := "1990-05-15"
	expectedDOB, _ := time.Parse("2006-01-02", dobString)

	// Parse the date of birth
	parsedDOB, err := ParseDOB(dobString)
	assert.NoError(t, err)

	// Ensure the parsed date of birth matches the expected value
	assert.Equal(t, expectedDOB, parsedDOB)
}

func TestParseTimestamp(t *testing.T) {
	// Test data
	timestampString := "2023-07-30 12:34:56"
	expectedTimestamp, _ := time.Parse("2006-01-02 15:04:05", timestampString)

	// Parse the timestamp
	parsedTimestamp, err := ParseTimestamp(timestampString)
	assert.NoError(t, err)

	// Ensure the parsed timestamp matches the expected value
	assert.Equal(t, expectedTimestamp, parsedTimestamp)
}

func TestValidateEmail(t *testing.T) {
	// Test valid email addresses
	validEmails := []string{
		"test@example.com",
		"user123@testdomain.co",
		"john.doe@example.co.uk",
	}

	for _, email := range validEmails {
		assert.True(t, ValidateEmail(email))
	}

	// Test invalid email addresses
	invalidEmails := []string{
		"invalid_email.com",
		"user@test..com",
		"invalid_email@",
		"@test.com",
	}

	for _, email := range invalidEmails {
		assert.False(t, ValidateEmail(email))
	}
}

func TestValidatePassword(t *testing.T) {
	validPasswords := []string{
		"Abcdefg1",   // Minimum length and all requirements satisfied
		"Password12", // Minimum length and all requirements satisfied
	}

	for _, password := range validPasswords {
		assert.True(t, ValidatePassword(password))
	}

	invalidPasswords := []string{
		"12345678",    // Missing uppercase letter and lowercase letter
		"password",    // Missing uppercase letter, lowercase letter, and digit
		"ABCDEFGH",    // Missing lowercase letter and digit
		"AbcdEfgh",    // Missing digit
		"Abcdefg",     // Missing uppercase letter and digit
		"abcdefgh",    // Missing uppercase letter and digit
		"ABC123",      // Missing lowercase letter
		"abcdefghij",  // Missing uppercase letter and lowercase letter
		"abcdefgh123", // Missing uppercase letter
		"ABCDEFG123",  // Missing lowercase letter
		"abcdefghijk", // Missing uppercase letter and lowercase letter and digit
		"abcdefghijK", // Missing digit
		"abcdefghij1", // Missing uppercase letter
	}

	for _, password := range invalidPasswords {
		assert.False(t, ValidatePassword(password))
	}
}

func TestValidateDateOfBirth(t *testing.T) {
	// Test valid date of birth (in the past)
	validDOB := "2000-01-01"
	assert.True(t, ValidateDateOfBirth(validDOB))

	// Test invalid date of birth (in the future)
	invalidDOB := "2050-01-01"
	assert.False(t, ValidateDateOfBirth(invalidDOB))
}

func TestValidateDateOfBirthInvalidFormat(t *testing.T) {
	// Test invalid date of birth (wrong format)
	invalidDOB := "01-01-2000"
	assert.False(t, ValidateDateOfBirth(invalidDOB))
}

func TestValidateDateOfBirthInvalidDate(t *testing.T) {
	// Test invalid date of birth (invalid date, e.g., 31st February)
	invalidDOB := "2000-02-31"
	assert.False(t, ValidateDateOfBirth(invalidDOB))
}

func TestValidateDateOfBirthFutureDate(t *testing.T) {
	// Test invalid date of birth (future date)
	futureDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	assert.False(t, ValidateDateOfBirth(futureDate))
}

func TestValidateDateOfBirthCurrentDate(t *testing.T) {
	// Test valid date of birth (current date)
	currentDate := time.Now().Format("2006-01-02")
	assert.True(t, ValidateDateOfBirth(currentDate))
}

func TestValidateDateOfBirthPastDate(t *testing.T) {
	// Test valid date of birth (past date)
	pastDate := time.Now().AddDate(-20, 0, 0).Format("2006-01-02")
	assert.True(t, ValidateDateOfBirth(pastDate))
}
