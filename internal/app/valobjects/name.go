package valobjects

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// Name errors
var (
	ErrInvalidName  = errors.New("Invalid name")
	ErrNameTooShort = fmt.Errorf("%w: min length allowed is %d", ErrInvalidName, minLength)
	ErrNameTooLong  = fmt.Errorf("%w: max length allowed is %d", ErrInvalidName, maxeLength)
	nilName         = ""
)

const (
	minLength  = 3
	maxeLength = 20
)

// Name represents a valid name
type Name string

// NewName creates a new valid name object
func NewName(n string) (Name, error) {
	switch len := len(strings.TrimSpace(n)); {
	case len < minLength:
		return "", ErrNameTooShort
	case len > maxeLength:
		return "", ErrNameTooLong
	default:
		return Name(n), nil
	}
}

// String implements the fmt.Stringer interface.
func (this Name) String() string {
	return string(this)
}

// Equals checks that two `names` are the same
func (this Name) Equals(value Value) bool {
	// that, ok := value.(Name)
	// return ok && this.name == that.name
	return this.String() == value.String()
}

// MarshalText used to serialize the object
func (this Name) MarshalText() ([]byte, error) {
	return []byte(this.String()), nil
}

// UnmarshalText deserializes the object and returns an error if it's invalid.
func (this *Name) UnmarshalText(data []byte) error {
	var err error
	*this, err = NewName(string(data))
	return err
}

// Value - Implementation of valuer for database/sql
func (this Name) Value() (driver.Value, error) {
	// ensuring value is a base driver.Value type, in this case `string`.
	return string(this), nil
}
