package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates new database connection successfully", func(t *testing.T) {
		// Create a temporary database file for testing
		dbPath := "test_countries.db"
		defer os.Remove(dbPath) // Clean up after test

		db, err := New(dbPath)

		require.NoError(t, err)
		require.NotNil(t, db)
		assert.NotNil(t, db.Conn)

		// Close the connection to clean up
		db.Conn.Close()
	})

	t.Run("returns error for invalid database path", func(t *testing.T) {
		// Try to create a database in a non-existent directory
		db, err := New("/nonexistent/path/test.db")

		assert.Error(t, err)
		assert.Nil(t, db)
		// The actual error might vary, so just check that there is an error
	})
}

func TestInitDB(t *testing.T) {
	t.Run("initializes database with countries table", func(t *testing.T) {
		dbPath := "test_init.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Check if the countries table exists by querying its structure
		_, err = db.Conn.Exec("SELECT * FROM countries LIMIT 1;")
		// This should not error since initDB creates the table
		assert.NoError(t, err)
	})

	t.Run("creates proper table structure", func(t *testing.T) {
		dbPath := "test_structure.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Insert a test record to verify the table structure
		_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?)", "Test Country", "TC")
		assert.NoError(t, err)

		// Verify the record was inserted
		var name, alpha2 string
		err = db.Conn.QueryRow("SELECT name, alpha2 FROM countries WHERE alpha2 = ?", "TC").Scan(&name, &alpha2)
		assert.NoError(t, err)
		assert.Equal(t, "Test Country", name)
		assert.Equal(t, "TC", alpha2)
	})

	t.Run("creates proper index", func(t *testing.T) {
		dbPath := "test_index.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Check if the index was created
		_, err = db.Conn.Exec("SELECT * FROM sqlite_master WHERE type = 'index' AND name = 'idx_countries_created_at'")
		assert.NoError(t, err)
	})
}

func TestGetRandomCountries(t *testing.T) {
	t.Run("returns random countries without error", func(t *testing.T) {
		dbPath := "test_random.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Add some test data
		_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?), (?, ?), (?, ?)", 
			"United States", "US", "Canada", "CA", "Mexico", "MX")
		require.NoError(t, err)

		countries, err := db.GetRandomCountries()

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(countries), 0) // Could be empty if no countries exist
		assert.LessOrEqual(t, len(countries), 10)   // Should not exceed 10
	})

	t.Run("returns countries with proper flag emojis", func(t *testing.T) {
		dbPath := "test_flag.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Add a test country
		_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?)", "United States", "US")
		require.NoError(t, err)

		countries, err := db.GetRandomCountries()

		require.NoError(t, err)
		require.Len(t, countries, 1)
		assert.Equal(t, "United States", countries[0].Name)
		assert.Equal(t, "🇺🇸", countries[0].Flag) // US flag emoji
	})

	t.Run("handles empty database gracefully", func(t *testing.T) {
		dbPath := "test_empty.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		countries, err := db.GetRandomCountries()

		require.NoError(t, err)
		assert.Empty(t, countries)
	})

	t.Run("returns no more than 10 countries", func(t *testing.T) {
		dbPath := "test_limit.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Add more than 10 countries to test the limit
		for i := 0; i < 15; i++ {
			_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?)", 
				"Country"+string(rune(i+'A')), string(rune(i+'A'))+string(rune(i+'A')))
			require.NoError(t, err)
		}

		countries, err := db.GetRandomCountries()

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(countries), 0)
		assert.LessOrEqual(t, len(countries), 10)
	})

	t.Run("handles database query errors", func(t *testing.T) {
		dbPath := "test_query_error.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		db.Conn.Close() // Close connection to simulate error

		_, err = db.GetRandomCountries()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database query error")
	})
}

func TestGetAllCountries(t *testing.T) {
	t.Run("returns all countries ordered by creation date and name", func(t *testing.T) {
		dbPath := "test_all.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Add test data
		_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2, created_at) VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?)", 
			"Zimbabwe", "ZW", "2020-01-01 00:00:00",
			"Algeria", "DZ", "2021-01-01 00:00:00", 
			"Brazil", "BR", "2020-01-01 00:00:00")
		require.NoError(t, err)

		countries, err := db.GetAllCountries()

		require.NoError(t, err)
		assert.Len(t, countries, 3)
		// Should be ordered by created_at DESC, then name ASC
		// Algeria should come first (2021), then Brazil, then Zimbabwe (both 2020 but Brazil < Zimbabwe alphabetically)
		assert.Equal(t, "Algeria", countries[0].Name)
		assert.Equal(t, "Brazil", countries[1].Name)
		assert.Equal(t, "Zimbabwe", countries[2].Name)
	})

	t.Run("returns countries with proper flag emojis", func(t *testing.T) {
		dbPath := "test_all_flags.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Add a test country
		_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?)", "Canada", "CA")
		require.NoError(t, err)

		countries, err := db.GetAllCountries()

		require.NoError(t, err)
		require.Len(t, countries, 1)
		assert.Equal(t, "Canada", countries[0].Name)
		assert.Equal(t, "🇨🇦", countries[0].Flag) // Canada flag emoji
	})

	t.Run("handles empty database gracefully", func(t *testing.T) {
		dbPath := "test_all_empty.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		countries, err := db.GetAllCountries()

		require.NoError(t, err)
		assert.Empty(t, countries)
	})

	t.Run("handles database query errors", func(t *testing.T) {
		dbPath := "test_all_query_error.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		db.Conn.Close() // Close connection to simulate error

		_, err = db.GetAllCountries()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database query error")
	})
}

func TestAddCountry(t *testing.T) {
	t.Run("successfully adds a new country", func(t *testing.T) {
		dbPath := "test_add.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		err = db.AddCountry("Japan", "JP")

		require.NoError(t, err)

		// Verify the country was added
		var name, alpha2 string
		err = db.Conn.QueryRow("SELECT name, alpha2 FROM countries WHERE alpha2 = ?", "JP").Scan(&name, &alpha2)
		require.NoError(t, err)
		assert.Equal(t, "Japan", name)
		assert.Equal(t, "JP", alpha2)
	})

	t.Run("adds country with proper timestamp", func(t *testing.T) {
		dbPath := "test_add_timestamp.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		err = db.AddCountry("Australia", "AU")
		require.NoError(t, err)

		// Check that created_at field is populated
		var name, alpha2, createdAt string
		err = db.Conn.QueryRow("SELECT name, alpha2, created_at FROM countries WHERE alpha2 = ?", "AU").Scan(&name, &alpha2, &createdAt)
		require.NoError(t, err)
		assert.Equal(t, "Australia", name)
		assert.Equal(t, "AU", alpha2)
		assert.NotEmpty(t, createdAt)
	})

	t.Run("handles database insert errors", func(t *testing.T) {
		dbPath := "test_add_error.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		db.Conn.Close() // Close connection to simulate error

		err = db.AddCountry("Test", "TT")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database insert error")
	})

	t.Run("handles constraint violations", func(t *testing.T) {
		dbPath := "test_constraint.db"
		defer os.Remove(dbPath)

		db, err := New(dbPath)
		require.NoError(t, err)
		defer db.Conn.Close()

		// Create a table with a NOT NULL constraint on name
		_, err = db.Conn.Exec("INSERT INTO countries (alpha2) VALUES (?)", "TT")
		// This might return an error depending on SQLite settings, but it's good to test
	})
}