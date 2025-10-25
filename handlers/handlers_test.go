package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	inertia "github.com/romsar/gonertia/v2"
)

func TestHomeHandler(t *testing.T) {
	// Create a mock Inertia instance using a test template
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a handler with a mock database
	// For tests, we'll create a temporary database
	db, err := database.New("test_home_handler.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	handler := New(db, in)

	t.Run("handles home page request successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.HomeHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("handles internal server error when Inertia fails", func(t *testing.T) {
		// Create a handler with a nil Inertia to test error handling
		failingHandler := &Handler{
			DB:  db,
			In:  nil, // This will cause an error
		}

		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.HomeHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestRandomCountriesHandler(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a temporary database for testing
	db, err := database.New("test_random_handler.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	// Add some test data
	_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?), (?, ?)", 
		"United States", "US", "Canada", "CA")
	require.NoError(t, err)

	handler := New(db, in)

	t.Run("returns random countries successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/random", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.RandomCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("handles error when database query fails", func(t *testing.T) {
		// Create a handler with a database that will fail
		failingDB, _ := database.New("test_random_fail.db")
		failingDB.Conn.Close() // Close connection to cause error

		failingHandler := &Handler{
			DB:  failingDB,
			In:  in,
		}

		req, err := http.NewRequest("GET", "/random", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.RandomCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestAllCountriesHandler(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a temporary database for testing
	db, err := database.New("test_all_handler.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	// Add some test data
	_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?), (?, ?)", 
		"United States", "US", "Canada", "CA")
	require.NoError(t, err)

	handler := New(db, in)

	t.Run("returns all countries successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/all", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.AllCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("handles error when database query fails", func(t *testing.T) {
		// Create a handler with a database that will fail
		failingDB, _ := database.New("test_all_fail.db")
		failingDB.Conn.Close() // Close connection to cause error

		failingHandler := &Handler{
			DB:  failingDB,
			In:  in,
		}

		req, err := http.NewRequest("GET", "/all", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.AllCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestNewCountriesHandler(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a temporary database for testing
	db, err := database.New("test_new_handler.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	handler := New(db, in)

	t.Run("adds new country successfully", func(t *testing.T) {
		// Create form data
		formData := models.NewCountry{
			Name: "Japan",
			Code: "JP",
		}
		
		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check that we get a redirect response (302)
		assert.Equal(t, http.StatusFound, rr.Code)
		
		// Verify the country was added to the database
		var name, alpha2 string
		err = db.Conn.QueryRow("SELECT name, alpha2 FROM countries WHERE alpha2 = ?", "JP").Scan(&name, &alpha2)
		require.NoError(t, err)
		assert.Equal(t, "Japan", name)
		assert.Equal(t, "JP", alpha2)
	})

	t.Run("handles invalid JSON in request body", func(t *testing.T) {
		invalidJSON := []byte(`{"name": "Test", "code":}`) // Invalid JSON

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("handles database insert error", func(t *testing.T) {
		// Create a handler with a database that will fail
		failingDB, _ := database.New("test_new_fail.db")
		failingDB.Conn.Close() // Close connection to cause error

		failingHandler := &Handler{
			DB:  failingDB,
			In:  in,
		}

		// Create form data
		formData := models.NewCountry{
			Name: "Test Country",
			Code: "TC",
		}
		
		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.NewCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("accepts valid country data", func(t *testing.T) {
		// Create form data
		formData := models.NewCountry{
			Name: "Australia",
			Code: "AU",
		}
		
		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check that we get a redirect response (302)
		assert.Equal(t, http.StatusFound, rr.Code)
		
		// Verify the country was added to the database
		var name, alpha2 string
		err = db.Conn.QueryRow("SELECT name, alpha2 FROM countries WHERE alpha2 = ?", "AU").Scan(&name, &alpha2)
		require.NoError(t, err)
		assert.Equal(t, "Australia", name)
		assert.Equal(t, "AU", alpha2)
	})
}

func TestNewHandler(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a database instance for testing
	db, err := database.New("test_new.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	t.Run("creates new handler instance", func(t *testing.T) {
		handler := New(db, in)

		require.NotNil(t, handler)
		assert.Equal(t, db, handler.DB)
		assert.Equal(t, in, handler.In)
	})
}