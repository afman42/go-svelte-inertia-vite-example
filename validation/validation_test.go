package validation

import (
	"testing"
)

func TestValidator_Required(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid non-empty string", "hello", true},
		{"empty string", "", false},
		{"whitespace only", "   ", false},
		{"tab only", "\t", false},
		{"newline only", "\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.Required(tt.input)
			if result != tt.expected {
				t.Errorf("Required(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidator_MinLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		min      int
		expected bool
	}{
		{"string longer than min", "hello", 3, true},
		{"string equal to min", "hel", 3, true},
		{"string shorter than min", "he", 3, false},
		{"empty string with positive min", "", 5, false},
		{"zero min with empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.MinLength(tt.input, tt.min)
			if result != tt.expected {
				t.Errorf("MinLength(%q, %d) = %v; want %v", tt.input, tt.min, result, tt.expected)
			}
		})
	}
}

func TestValidator_MaxLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		max      int
		expected bool
	}{
		{"string shorter than max", "hi", 5, true},
		{"string equal to max", "hello", 5, true},
		{"string longer than max", "hello world", 5, false},
		{"zero max with empty string", "", 0, true},
		{"zero max with non-empty string", "hello", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.MaxLength(tt.input, tt.max)
			if result != tt.expected {
				t.Errorf("MaxLength(%q, %d) = %v; want %v", tt.input, tt.max, result, tt.expected)
			}
		})
	}
}

func TestValidator_IsEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "user@sub.domain.com", true},
		{"invalid email no domain", "test@", false},
		{"invalid email no @", "test.example.com", false},
		{"valid email with numbers", "user123@test-domain.com", true},
		{"invalid email special chars", "us@er@domain.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.IsEmail(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmail(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidator_In(t *testing.T) {
	list := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"value in list", "banana", true},
		{"value not in list", "orange", false},
		{"exact match only", "ban", false},
		{"case sensitive", "Banana", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := In(tt.input, list...)
			if result != tt.expected {
				t.Errorf("In(%q, %v) = %v; want %v", tt.input, list, result, tt.expected)
			}
		})
	}
}

func TestValidator_MinInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		min      int
		expected bool
	}{
		{"value greater than min", 10, 5, true},
		{"value equal to min", 5, 5, true},
		{"value less than min", 3, 5, false},
		{"negative values", -5, -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.MinInt(tt.input, tt.min)
			if result != tt.expected {
				t.Errorf("MinInt(%d, %d) = %v; want %v", tt.input, tt.min, result, tt.expected)
			}
		})
	}
}

func TestValidator_MaxInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		max      int
		expected bool
	}{
		{"value less than max", 3, 5, true},
		{"value equal to max", 5, 5, true},
		{"value greater than max", 10, 5, false},
		{"negative values", -10, -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.MaxInt(tt.input, tt.max)
			if result != tt.expected {
				t.Errorf("MaxInt(%d, %d) = %v; want %v", tt.input, tt.max, result, tt.expected)
			}
		})
	}
}

func TestValidator_MatchesPattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected bool
	}{
		{"valid pattern match", "abc123", "[a-z]+[0-9]+", true},
		{"invalid pattern match", "123abc", "[a-z]+[0-9]+", false},
		{"exact match", "hello", "^hello$", true},
		{"partial no match", "hello world", "^hello$", false},
		{"invalid regex", "test", "[invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result := v.MatchesPattern(tt.input, tt.pattern)
			if result != tt.expected {
				t.Errorf("MatchesPattern(%q, %q) = %v; want %v", tt.input, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestValidator_ValidateString(t *testing.T) {
	v := New()
	
	// Test valid string
	v.ValidateString("name", "John", 2, 50, "")
	if !v.Valid() {
		t.Errorf("ValidateString with valid input should not produce errors")
	}
	
	// Test invalid strings
	v = New()
	v.ValidateString("name", "", 2, 50, "")  // required check
	if v.Valid() {
		t.Errorf("ValidateString with empty string should produce errors")
	}
	
	v = New()
	v.ValidateString("name", "A", 2, 50, "")  // min length check
	if v.Valid() {
		t.Errorf("ValidateString with too short string should produce errors")
	}
	
	v = New()
	v.ValidateString("name", "A", 2, 1, "")  // max length check
	if v.Valid() {
		t.Errorf("ValidateString with too long string should produce errors")
	}
}

func TestValidator_ValidateEmail(t *testing.T) {
	v := New()
	
	// Test valid email
	v.ValidateEmail("email", "test@example.com")
	if !v.Valid() {
		t.Errorf("ValidateEmail with valid email should not produce errors")
	}
	
	// Test invalid email
	v = New()
	v.ValidateEmail("email", "invalid-email")
	if v.Valid() {
		t.Errorf("ValidateEmail with invalid email should produce errors")
	}
}

func TestValidator_Valid(t *testing.T) {
	v := New()
	if !v.Valid() {
		t.Errorf("New validator should be valid initially")
	}
	
	v.AddError("field", "error message")
	if v.Valid() {
		t.Errorf("Validator with errors should not be valid")
	}
}

func TestValidator_AddError(t *testing.T) {
	v := New()
	v.AddError("field1", "error1")
	
	if len(v.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(v.Errors))
	}
	
	if v.Errors["field1"] != "error1" {
		t.Errorf("Expected error1, got %s", v.Errors["field1"])
	}
	
	// Adding same field should not overwrite
	v.AddError("field1", "error2")
	if v.Errors["field1"] != "error1" {
		t.Errorf("AddError should not overwrite existing errors")
	}
}

func TestValidator_Check(t *testing.T) {
	v := New()
	v.Check(false, "field", "error message")
	
	if len(v.Errors) != 1 || v.Errors["field"] != "error message" {
		t.Errorf("Check with false condition should add error")
	}
	
	v = New()
	v.Check(true, "field", "error message")
	
	if len(v.Errors) != 0 {
		t.Errorf("Check with true condition should not add error")
	}
}

func TestValidator_GetErrors(t *testing.T) {
	v := New()
	v.AddError("field1", "message1")
	v.AddError("field2", "message2")
	
	errors := v.GetErrors()
	
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}
	
	// Check that errors are properly mapped
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.Field] = err.Message
	}
	
	if errorMap["field1"] != "message1" || errorMap["field2"] != "message2" {
		t.Errorf("GetErrors didn't return correct error mappings")
	}
}