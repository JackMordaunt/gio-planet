// Package value defines some `form.Value` implementations for common data types.
package value

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~jackmordaunt/gio-planet/form"
	"git.sr.ht/~jackmordaunt/gio-planet/form/parse"
)

// Int maps text to an integer number.
type Int struct {
	Value   *int
	Default int
}

func (v Int) To() (string, error) {
	var n = *v.Value
	if n == 0 {
		n = v.Default
	}
	return strconv.Itoa(n), nil
}

func (v Int) From(text string) (err error) {
	*v.Value, err = parse.Int(text)
	return err
}

func (v Int) Clear() {
	*v.Value = 0
}

// Float maps text to a floating point number.
type Float struct {
	Value *float64
}

func (v Float) To() (string, error) {
	return strconv.FormatFloat(*v.Value, 'f', 2, 64), nil
}

func (v Float) From(text string) (err error) {
	*v.Value, err = parse.Float(text)
	return err
}

func (v Float) Clear() {
	*v.Value = 0
}

// Text wraps a text value.
type Text struct {
	Value *string
}

func (v Text) To() (string, error) {
	return *v.Value, nil
}

func (v Text) From(text string) error {
	*v.Value = text
	return nil
}

func (v Text) Clear() {
	*v.Value = ""
}

// Days maps text to 24 hour units of time.
type Days struct {
	Value *time.Duration
}

func (v Days) To() (string, error) {
	days := (*v.Value) / (time.Hour * 24)
	return strconv.Itoa(int(days)), nil
}

func (v Days) From(text string) (err error) {
	*v.Value, err = parse.Day(text)
	return err
}

func (v Days) Clear() {
	*v.Value = time.Hour * 24
}

// Required errors when the field is empty.
type Required struct {
	form.Value
}

func (v Required) From(text string) error {
	if len(strings.TrimSpace(text)) == 0 {
		return fmt.Errorf("required")
	}
	return v.Value.From(text)
}

// Date maps text to a structured date.
type Date struct {
	Value   *time.Time
	Default time.Time
}

func (v Date) To() (string, error) {
	var t = *v.Value
	if t.IsZero() {
		t = v.Default
	}
	return fmt.Sprintf("%d/%d/%d", t.Day(), t.Month(), t.Year()), nil
}

func (v Date) From(text string) (err error) {
	*v.Value, err = parse.Date(text)
	return err
}

func (v Date) Clear() {
	*v.Value = time.Now()
}
