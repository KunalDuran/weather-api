package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDOB(t *testing.T) {
	date, err := ParseDOB("2023-07-27")
	assert.NoError(t, err)

	expectedDate := time.Date(2023, 07, 27, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedDate, date)
}

func TestParseDOBWithInvalidDate(t *testing.T) {
	_, err := ParseDOB("2023-13-32")
	assert.Error(t, err)
}

func TestParseTimestamp(t *testing.T) {
	date, err := ParseTimestamp("2023-07-27 15:04:05")
	assert.NoError(t, err)

	expectedDate := time.Date(2023, 07, 27, 15, 4, 5, 0, time.UTC)
	assert.Equal(t, expectedDate, date)
}

func TestParseTimestampWithInvalidTimestamp(t *testing.T) {
	_, err := ParseTimestamp("2023-13-32 25:61:62")
	assert.Error(t, err)
}
