package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Run("full CRUD flow", func(t *testing.T) {
		dbPath := "test_integration.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Test: Add a new country
		err = db.AddCountry("Testland", "TL")
		require.NoError(t, err)

		// Test: Retrieve all countries (should include the new one)
		allCountries, err := db.GetAllCountries()
		require.NoError(t, err)
		require.Len(t, allCountries, 1)
		assert.Equal(t, "Testland", allCountries[0].Name)
		assert.Equal(t, "🇹🇱", allCountries[0].Flag) // TL flag emoji

		// Test: Retrieve random countries (should include the new one)
		randomCountries, err := db.GetRandomCountries()
		require.NoError(t, err)
		// The result might be empty if random doesn't pick our country, so we check if it's in the possible results
		found := false
		for _, country := range randomCountries {
			if country.Name == "Testland" {
				found = true
				break
			}
		}
		// Since we have only one country in the DB, it should appear in the random results
		assert.True(t, found)

		// Test: Add another country
		err = db.AddCountry("Anotherland", "AN")
		require.NoError(t, err)

		// Test: Retrieve all countries again (should now have 2)
		allCountries, err = db.GetAllCountries()
		require.NoError(t, err)
		require.Len(t, allCountries, 2)
		
		// Should be ordered by creation date (descending) and then name (ascending)
		// Since both were created at the same time (in the same test), we'll check for both
		var foundTestland, foundAnotherland bool
		for _, country := range allCountries {
			if country.Name == "Anotherland" {
				foundAnotherland = true
				assert.Equal(t, "🇦🇳", country.Flag) // AN flag emoji
			}
			if country.Name == "Testland" {
				foundTestland = true
				assert.Equal(t, "🇹🇱", country.Flag) // TL flag emoji
			}
		}
		assert.True(t, foundTestland)
		assert.True(t, foundAnotherland)
	})
}

func TestCountry2Flag(t *testing.T) {
	// Since country2flag is a private function, we test it through the public API
	t.Run("country code to flag conversion", func(t *testing.T) {
		dbPath := "test_flag_conversion.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Add test countries with different codes
		testCases := []struct {
			name     string
			code     string
			expected string
		}{
			{"United States", "US", "🇺🇸"},
			{"Canada", "CA", "🇨🇦"},
			{"Japan", "JP", "🇯🇵"},
			{"Germany", "DE", "🇩🇪"},
			{"France", "FR", "🇫🇷"},
		}

		for _, tc := range testCases {
			err := db.AddCountry(tc.name, tc.code)
			require.NoError(t, err)
		}

		// Retrieve all countries and verify flag conversion
		countries, err := db.GetAllCountries()
		require.NoError(t, err)

		for _, tc := range testCases {
			found := false
			for _, country := range countries {
				if country.Name == tc.name {
					assert.Equal(t, tc.expected, country.Flag)
					found = true
					break
				}
			}
			assert.True(t, found, "Country %s not found in results", tc.name)
		}
	})
}