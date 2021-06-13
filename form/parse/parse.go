// Package parse implements text parsing for common data types.
package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Int parses an integer from digit characters.
func Int(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}
	return n, nil
}

// Float parses a floating point number from digit characters.
func Float(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}
	return n, nil
}

// Uint parses an unsigned integer from digit characters.
func Uint(s string) (uint, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	} else if n < 1 {
		return 0, fmt.Errorf("must be an amount greater than 0")
	}
	return uint(n), nil
}

// Day parses a day from digit characters.
func Day(s string) (time.Duration, error) {
	n, err := Uint(s)
	if err != nil {
		return time.Duration(0), err
	}
	return time.Hour * 24 * time.Duration(n), nil
}

// FieldRequired ensures that a string is not empty.
func FieldRequired(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("required")
	}
	return s, nil
}

// FormatTime formats a time object into a string.
func FormatTime(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Day(), t.Month(), t.Year())
}

// Date parses a time object from a textual dd/mm/yyyy format.
func Date(s string) (date time.Time, err error) {
	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return date, fmt.Errorf("must be dd/mm/yyyy")
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return date, fmt.Errorf("year not a number: %s", parts[2])
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return date, fmt.Errorf("month not a number: %s", parts[2])
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return date, fmt.Errorf("day not a number: %s", parts[2])
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}
