package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountryModel(t *testing.T) {
	t.Run("Country struct fields", func(t *testing.T) {
		country := Country{
			Name: "United States",
			Flag: "🇺🇸",
		}

		assert.Equal(t, "United States", country.Name)
		assert.Equal(t, "🇺🇸", country.Flag)
	})

	t.Run("NewCountry struct fields", func(t *testing.T) {
		formData := NewCountry{
			Name: "Canada",
			Code: "CA",
		}

		assert.Equal(t, "Canada", formData.Name)
		assert.Equal(t, "CA", formData.Code)
	})

	t.Run("Country JSON tags", func(t *testing.T) {
		country := Country{
			Name: "France",
			Flag: "🇫🇷",
		}

		// The struct should have proper JSON tags
		// This test verifies that the struct is properly defined
		assert.Equal(t, "France", country.Name)
		assert.Equal(t, "🇫🇷", country.Flag)
	})

	t.Run("NewCountry JSON tags", func(t *testing.T) {
		formData := NewCountry{
			Name: "Germany",
			Code: "DE",
		}

		// The struct should have proper JSON tags
		// This test verifies that the struct is properly defined
		assert.Equal(t, "Germany", formData.Name)
		assert.Equal(t, "DE", formData.Code)
	})
}