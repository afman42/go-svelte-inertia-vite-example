package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

// Error returns a string representation of all validation errors
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Validator contains validation rules and errors
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the Errors map is empty
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error to the Errors map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error if the condition is false
func (v *Validator) Check(condition bool, key, message string) {
	if !condition {
		v.AddError(key, message)
	}
}

// GetErrors returns all validation errors as a ValidationErrors slice
func (v *Validator) GetErrors() ValidationErrors {
	var errors ValidationErrors
	for field, message := range v.Errors {
		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return errors
}

// In returns true if the value is present in the list
func In(value string, list ...string) bool {
	for _, item := range list {
		if value == item {
			return true
		}
	}
	return false
}

// MatchesPattern returns true if the value matches the given regular expression pattern
func (v *Validator) MatchesPattern(value, pattern string) bool {
	matched, err := regexp.MatchString(pattern, value)
	return err == nil && matched
}

// Required returns true if the value is not empty
func (v *Validator) Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinLength returns true if the value has at least minLength characters
func (v *Validator) MinLength(value string, minLength int) bool {
	return utf8.RuneCountInString(value) >= minLength
}

// MaxLength returns true if the value has at most maxLength characters
func (v *Validator) MaxLength(value string, maxLength int) bool {
	return utf8.RuneCountInString(value) <= maxLength
}

// MinInt returns true if the value is greater than or equal to min
func (v *Validator) MinInt(value, min int) bool {
	return value >= min
}

// MaxInt returns true if the value is less than or equal to max
func (v *Validator) MaxInt(value, max int) bool {
	return value <= max
}

// MinFloat returns true if the value is greater than or equal to min
func (v *Validator) MinFloat(value, min float64) bool {
	return value >= min
}

// MaxFloat returns true if the value is less than or equal to max
func (v *Validator) MaxFloat(value, max float64) bool {
	return value <= max
}

// EmailRX is a regular expression for email validation
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEmail returns true if the value matches the email regular expression
func (v *Validator) IsEmail(value string) bool {
	return EmailRX.MatchString(value)
}

// IsInt returns true if the value is a valid integer string
func (v *Validator) IsInt(value string) bool {
	_, err := fmt.Sscanf(value, "%d")
	return err == nil
}

// IsFloat returns true if the value is a valid float string
func (v *Validator) IsFloat(value string) bool {
	_, err := fmt.Sscanf(value, "%f")
	return err == nil
}

// IsBool returns true if the value is a valid boolean string
func (v *Validator) IsBool(value string) bool {
	_, err := fmt.Sscanf(value, "%t")
	return err == nil
}

// ValidateString validates a string value based on multiple criteria
func (v *Validator) ValidateString(key, value string, minLen int, maxLen int, pattern string) {
	if !v.Required(value) {
		v.AddError(key, "This field cannot be blank")
		return
	}

	if !v.MinLength(value, minLen) {
		v.AddError(key, fmt.Sprintf("This field must be at least %d characters long", minLen))
	}

	if !v.MaxLength(value, maxLen) {
		v.AddError(key, fmt.Sprintf("This field must be at most %d characters long", maxLen))
	}

	if pattern != "" && !v.MatchesPattern(value, pattern) {
		v.AddError(key, "This field format is invalid")
	}
}

// ValidateInt validates an integer value based on min/max constraints
func (v *Validator) ValidateInt(key string, value, min, max int) {
	if !v.MinInt(value, min) {
		v.AddError(key, fmt.Sprintf("This field must be at least %d", min))
	}

	if !v.MaxInt(value, max) {
		v.AddError(key, fmt.Sprintf("This field must be at most %d", max))
	}
}

// ValidateFloat validates a float64 value based on min/max constraints
func (v *Validator) ValidateFloat(key string, value, min, max float64) {
	if !v.MinFloat(value, min) {
		v.AddError(key, fmt.Sprintf("This field must be at least %.2f", min))
	}

	if !v.MaxFloat(value, max) {
		v.AddError(key, fmt.Sprintf("This field must be at most %.2f", max))
	}
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(key, value string) {
	if !v.IsEmail(value) {
		v.AddError(key, "This field must be a valid email address")
	}
}

// ValidateURL validates if the value is a valid URL format
func (v *Validator) ValidateURL(key, value string) {
	// Simple URL validation using regex
	urlPattern := `^(https?://)?([a-zA-Z0-9.-]+\.[a-zA-Z]{2,4})(/[^\s]*)?$`
	if !v.MatchesPattern(value, urlPattern) {
		v.AddError(key, "This field must be a valid URL")
	}
}

// ValidatePhone validates if the value is a phone number format
func (v *Validator) ValidatePhone(key, value string) {
	// Simple phone number validation
	phonePattern := `^[\+]?[1-9][\d]{0,15}$`
	if !v.MatchesPattern(value, phonePattern) {
		v.AddError(key, "This field must be a valid phone number")
	}
}

// ValidateCreditCard validates if the value is a credit card format
func (v *Validator) ValidateCreditCard(key, value string) {
	// Simple credit card validation (just numbers with spaces/dashes)
	cardPattern := `^[\d\s-]{13,19}$`
	if !v.MatchesPattern(value, cardPattern) {
		v.AddError(key, "This field must be a valid credit card number")
	}
}

// ValidateMatch checks if two values match
func (v *Validator) ValidateMatch(key, value, matchValue, matchKey string) {
	if value != matchValue {
		v.AddError(key, fmt.Sprintf("This field must match the %s field", matchKey))
	}
}
