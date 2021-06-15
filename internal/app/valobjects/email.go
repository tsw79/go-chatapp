package valobjects

import (
	"database/sql/driver"
	"errors"
	"regexp"
)

var ErrInvalidEmailAddress = errors.New("Invalid email address")

// Emails represents a valid email address
type Email string

// NewEmail creates a new Email instance
func NewEmail(email string) (Email, error) {
	valid, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email)
	if !valid {
		return "", ErrInvalidEmailAddress
	}
	return Email(email), nil
}

// String implements the fmt.Stringer interface
func (this Email) String() string {
	return string(this)
}

// Equals checks that two emails are the same
func (this Email) Equals(value Value) bool {
	return this.String() == value.String()
}

// Value - Implementation of valuer for database/sql
func (this Email) Value() (driver.Value, error) {
	// ensuring value is a base driver.Value type, in this case `string`.
	return string(this), nil
}

// TODO Is this really necessary? Look at type Name, it works without it!
func (this *Email) Scan(value interface{}) error {
	if valueStr, err := driver.String.ConvertValue(value); err == nil {
		// Make sure this is a string
		if val, ok := valueStr.(string); ok {
			// Set the value of the pointer this to Email(val)
			*this = Email(val)
			return nil
		}
	}
	return errors.New("Failed to scan type Email")
}
