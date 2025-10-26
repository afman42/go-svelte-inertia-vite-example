# Validation Package

A reliable and easy-to-use input validation package for Go applications.

## Features

- String validation (required, min/max length, pattern matching)
- Numeric validation (int, float with min/max constraints)
- Email validation
- URL validation
- Phone number validation
- Credit card validation
- Custom validation rules
- Comprehensive error handling
- Easy integration with web forms and APIs

## Installation

```bash
go get github.com/afman42/go-svelte-inertia/validation
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/afman42/go-svelte-inertia/validation"
)

func main() {
    // Create a new validator
    v := validation.New()
    
    // Validate different types of data
    v.ValidateString("name", "John", 2, 50, "") // min length 2, max length 50
    v.ValidateEmail("email", "john@example.com")
    v.ValidateInt("age", 25, 18, 120) // min 18, max 120
    
    if !v.Valid() {
        // Print validation errors
        for field, message := range v.Errors {
            fmt.Printf("Error in field %s: %s\n", field, message)
        }
    } else {
        fmt.Println("All validations passed!")
    }
}
```

## Available Validation Functions

### String Validation
```go
v.ValidateString(key, value string, minLen, maxLen int, pattern string)
```

### Email Validation
```go
v.ValidateEmail(key, value string)
```

### Integer Validation
```go
v.ValidateInt(key string, value, min, max int)
```

### Float Validation
```go
v.ValidateFloat(key string, value, min, max float64)
```

### URL Validation
```go
v.ValidateURL(key, value string)
```

### Phone Number Validation
```go
v.ValidatePhone(key, value string)
```

### Credit Card Validation
```go
v.ValidateCreditCard(key, value string)
```

### Field Matching
```go
v.ValidateMatch(key, value, matchValue, matchKey string)
```

## Standalone Validation Functions

The package also provides standalone validation functions:

```go
// Check if string is required (not empty after trimming)
validation.Required(value string) bool

// Check string length
validation.MinLength(value string, minLength int) bool
validation.MaxLength(value string, maxLength int) bool

// Check numeric values
validation.MinInt(value, min int) bool
validation.MaxInt(value, max int) bool
validation.MinFloat(value, min float64) bool
validation.MaxFloat(value, max float64) bool

// Check email format
validation.IsEmail(value string) bool

// Check if value is in a list
validation.In(value string, list ...string) bool

// Check against regex pattern
validation.MatchesPattern(value, pattern string) bool
```

## Error Handling

The validation package provides comprehensive error handling:

```go
v := validation.New()

// Check if validation passed
if v.Valid() {
    // No errors
} else {
    // Get all errors
    errors := v.GetErrors() // Returns []ValidationError
    for _, err := range errors {
        fmt.Printf("Field: %s, Message: %s\n", err.Field, err.Message)
    }
    
    // Or access the errors map directly
    for field, message := range v.Errors {
        fmt.Printf("Error in %s: %s\n", field, message)
    }
}
```

## Advanced Usage

### Conditional Validation
```go
v := validation.New()

// Use Check to validate conditions
v.Check(len(password) >= 8, "password", "Password must be at least 8 characters")
v.Check(validation.In(country, "US", "CA", "UK"), "country", "Invalid country")
```

### Custom Validation
```go
v := validation.New()

// Perform custom validation logic
if !isValidCustomFormat(input) {
    v.AddError("custom_field", "Custom validation failed")
}
```

## License

This project is licensed under the MIT License.