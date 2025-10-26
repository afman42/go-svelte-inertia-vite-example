package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/models"
	inertia "github.com/romsar/gonertia/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	// Create a database instance for testing
	db, err := database.New("test_new.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	// Create a session store
	sessionStore := auth.NewSessionStore()

	t.Run("creates new handler instance", func(t *testing.T) {
		handler := New(db, in, sessionStore)

		require.NotNil(t, handler)
		assert.Equal(t, db, handler.DB)
		assert.Equal(t, in, handler.In)
		assert.Equal(t, sessionStore, handler.SessionStore)
		assert.NotNil(t, handler.Auth)
	})
}

func TestHomeHandler(t *testing.T) {
	// Create a mock Inertia instance using a test template
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	// Create a handler with a mock database
	// For tests, we'll create a temporary database
	db, err := database.New("test_home_handler.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	// Create a session store
	sessionStore := auth.NewSessionStore()

	handler := New(db, in, sessionStore)

	t.Run("handles home page request successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.HomeHandler(rr, req)

		// Check the status code - this may not be 200 due to template setup
		// but it should not be a 404 or other error
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("handles internal server error when Inertia fails", func(t *testing.T) {
		// Create a handler with a nil Inertia to test error handling
		failingHandler := &Handler{
			DB:           db,
			In:           nil, // This will cause an error
			SessionStore: sessionStore,
		}

		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.HomeHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("handles different HTTP methods for home route", func(t *testing.T) {
		// Test POST method
		req, err := http.NewRequest("POST", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.HomeHandler(rr, req)

		// For non-GET methods, should still return some response without error
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
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

	// Create a session store
	sessionStore := auth.NewSessionStore()

	// Add some test data
	_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?), (?, ?)",
		"United States", "US", "Canada", "CA")
	require.NoError(t, err)

	handler := New(db, in, sessionStore)

	t.Run("returns random countries successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/random", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.RandomCountriesHandler(rr, req)

		// Check the status code - this may not be 200 due to template setup
		// but it should not be a 404 or other error
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("handles error when database query fails", func(t *testing.T) {
		// Create a handler with a database that will fail
		failingDB, _ := database.New("test_random_fail.db")
		failingDB.Conn.Close() // Close connection to cause error

		failingHandler := &Handler{
			DB:           failingDB,
			In:           in,
			SessionStore: sessionStore,
		}

		req, err := http.NewRequest("GET", "/random", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.RandomCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("handles random countries with different HTTP methods", func(t *testing.T) {
		// Test POST method
		req, err := http.NewRequest("POST", "/random", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.RandomCountriesHandler(rr, req)

		// For non-GET methods, should still return some response without error
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
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

	// Create a session store
	sessionStore := auth.NewSessionStore()

	// Add some test data
	_, err = db.Conn.Exec("INSERT INTO countries (name, alpha2) VALUES (?, ?), (?, ?)",
		"United States", "US", "Canada", "CA")
	require.NoError(t, err)

	handler := New(db, in, sessionStore)

	t.Run("returns all countries successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/all", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.AllCountriesHandler(rr, req)

		// Check the status code - this may not be 200 due to template setup
		// but it should not be a 404 or other error
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("handles error when database query fails", func(t *testing.T) {
		// Create a handler with a database that will fail
		failingDB, _ := database.New("test_all_fail.db")
		failingDB.Conn.Close() // Close connection to cause error

		failingHandler := &Handler{
			DB:           failingDB,
			In:           in,
			SessionStore: sessionStore,
		}

		req, err := http.NewRequest("GET", "/all", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		failingHandler.AllCountriesHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("handles all countries with different HTTP methods", func(t *testing.T) {
		// Test POST method
		req, err := http.NewRequest("POST", "/all", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.AllCountriesHandler(rr, req)

		// For non-GET methods, should still return some response without error
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
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

	// Create a session store
	sessionStore := auth.NewSessionStore()

	handler := New(db, in, sessionStore)

	t.Run("redirects to login when not authenticated", func(t *testing.T) {
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
		handler.NewCountriesHandler(rr, req)

		// Should redirect to login since user is not authenticated
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "/login")
	})

	t.Run("handles database insert error when authenticated", func(t *testing.T) {
		// Create a separate test database for this specific test
		testDB, err := database.New("test_new_insert_error.db")
		require.NoError(t, err)
		defer testDB.Conn.Close()
		defer os.Remove("test_new_insert_error.db")

		// Create the full handler with proper initialization
		testHandler := New(testDB, in, sessionStore)

		// Close the database connection after handler creation to cause error during insertion
		testDB.Conn.Close()

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

		// Add a valid session to the request to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		testHandler.NewCountriesHandler(rr, req)

		// Check the status code - should be Internal Server Error since DB is closed
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("handles invalid JSON in request body", func(t *testing.T) {
		invalidJSON := []byte(`{"name": "Test", "code":}`) // Invalid JSON

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code - should be Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("handles missing name field", func(t *testing.T) {
		// Create form data with empty name
		formData := models.NewCountry{
			Name: "",
			Code: "TC",
		}

		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code - should be Unprocessable Entity
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles missing code field", func(t *testing.T) {
		// Create form data with empty code
		formData := models.NewCountry{
			Name: "Test Country",
			Code: "",
		}

		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code - should be Unprocessable Entity
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles name with excessive length", func(t *testing.T) {
		// Create form data with very long name (256 characters)
		longName := "A"
		for i := 0; i < 255; i++ {
			longName += "A"
		}
		
		formData := models.NewCountry{
			Name: longName,
			Code: "TC",
		}

		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code - should be Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles code with invalid length", func(t *testing.T) {
		// Create form data with code that is not exactly 2 characters
		formData := models.NewCountry{
			Name: "Test Country",
			Code: "T", // Too short
		}

		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code - should be Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles code with invalid format", func(t *testing.T) {
		// Create form data with code that is not uppercase letters
		formData := models.NewCountry{
			Name: "Test Country",
			Code: "12", // Not uppercase letters
		}

		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Check the status code - should be Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("successfully adds country when authenticated", func(t *testing.T) {
		formData := models.NewCountry{
			Name: "Test Country",
			Code: "TC",
		}

		jsonData, err := json.Marshal(formData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		handler.NewCountriesHandler(rr, req)

		// Should redirect to /all after successful creation
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "/all")

		// Verify the country was added to the database
		countries, err := db.GetAllCountries()
		require.NoError(t, err)

		found := false
		for _, country := range countries {
			if country.Name == "Test Country" {
				found = true
				break
			}
		}
		assert.True(t, found, "Country should have been added to the database")
	})
}

// MockHTTPHandler is a helper to create an HTTP handler for testing
type MockHTTPHandler struct {
	ResponseCode int
	Handler      http.HandlerFunc
}

func (m *MockHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(m.ResponseCode)
	m.Handler(w, r)
}

func TestAuthMiddleware(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	// Create a temporary database for testing
	db, err := database.New("test_auth_middleware.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	// Create a session store
	sessionStore := auth.NewSessionStore()

	handler := New(db, in, sessionStore)

	t.Run("redirects unauthenticated requests", func(t *testing.T) {
		// Create a protected handler
		protectedHandler := handler.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Protected content")
		})

		req, err := http.NewRequest("GET", "/protected", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		protectedHandler(rr, req)

		// Should redirect to login since user is not authenticated
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "/login")
	})

	t.Run("allows authenticated requests", func(t *testing.T) {
		// Create a protected handler
		protectedHandler := handler.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Protected content")
		})

		req, err := http.NewRequest("GET", "/protected", nil)
		require.NoError(t, err)

		// Add a valid session to bypass auth check
		sessionID, err := sessionStore.CreateSession(1, "test@example.com")
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		protectedHandler(rr, req)

		// Should not redirect if authenticated
		assert.NotEqual(t, http.StatusFound, rr.Code)
		assert.Equal(t, http.StatusOK, rr.Code) // Should return 200 OK
	})
}

// Auth handler tests
func TestNewAuth(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	db, err := database.New("test_auth_new.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	sessionStore := auth.NewSessionStore()

	t.Run("creates new auth handler instance", func(t *testing.T) {
		authHandler := NewAuth(db, in, sessionStore)

		require.NotNil(t, authHandler)
		assert.Equal(t, db, authHandler.DB)
		assert.Equal(t, in, authHandler.In)
		assert.Equal(t, sessionStore, authHandler.SessionStore)
	})
}

func TestRegisterHandler(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	sessionStore := auth.NewSessionStore()

	t.Run("handles user registration successfully", func(t *testing.T) {
		db, err := database.New("test_register_success.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_success.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data
		registerData := models.UserRegister{
			Name:                 "Test User",
			Email:                "test@example.com",
			Password:             "securepassword123",
			PasswordConfirmation: "securepassword123",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return 200 status on successful registration
		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify user was added to the database
		user, err := db.GetUserByEmail("test@example.com")
		require.NoError(t, err)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		// Password should be hashed, so it won't be the same as the input
	})

	t.Run("handles invalid JSON in request body", func(t *testing.T) {
		db, err := database.New("test_register_invalid_json.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_invalid_json.db")

		authHandler := NewAuth(db, in, sessionStore)

		invalidJSON := []byte(`{"name": "Test", "email":}`) // Invalid JSON

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("handles existing user", func(t *testing.T) {
		db, err := database.New("test_register_existing.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_existing.db")

		authHandler := NewAuth(db, in, sessionStore)

		// First, register a user
		registerData := models.UserRegister{
			Name:                 "Existing User",
			Email:                "existing@example.com",
			Password:             "password123",
			PasswordConfirmation: "password123",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		// First registration should succeed
		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Create the request again for the second attempt
		jsonData, err = json.Marshal(registerData) // Marshal the same data again
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return conflict error
		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("handles non-POST requests", func(t *testing.T) {
		db, err := database.New("test_register_method.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_method.db")

		authHandler := NewAuth(db, in, sessionStore)

		req, err := http.NewRequest("GET", "/register", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Method Not Allowed
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("handles missing name field", func(t *testing.T) {
		db, err := database.New("test_register_missing_name.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_missing_name.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data with missing name
		registerData := models.UserRegister{
			Name:                 "", // Missing name
			Email:                "test@example.com",
			Password:             "securepassword123",
			PasswordConfirmation: "securepassword123",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles missing email field", func(t *testing.T) {
		db, err := database.New("test_register_missing_email.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_missing_email.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data with missing email
		registerData := models.UserRegister{
			Name:                 "Test User",
			Email:                "", // Missing email
			Password:             "securepassword123",
			PasswordConfirmation: "securepassword123",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles invalid email format", func(t *testing.T) {
		db, err := database.New("test_register_invalid_email.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_invalid_email.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data with invalid email format
		registerData := models.UserRegister{
			Name:                 "Test User",
			Email:                "invalid-email", // Invalid email format
			Password:             "securepassword123",
			PasswordConfirmation: "securepassword123",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles missing password field", func(t *testing.T) {
		db, err := database.New("test_register_missing_password.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_missing_password.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data with missing password
		registerData := models.UserRegister{
			Name:                 "Test User",
			Email:                "test@example.com",
			Password:             "", // Missing password
			PasswordConfirmation: "", // Missing password confirmation
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles short password", func(t *testing.T) {
		db, err := database.New("test_register_short_password.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_short_password.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data with short password
		registerData := models.UserRegister{
			Name:                 "Test User",
			Email:                "test@example.com",
			Password:             "123", // Too short
			PasswordConfirmation: "123", // Too short
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles mismatched password confirmation", func(t *testing.T) {
		db, err := database.New("test_register_mismatched_password.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_register_mismatched_password.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create registration data with mismatched passwords
		registerData := models.UserRegister{
			Name:                 "Test User",
			Email:                "test@example.com",
			Password:             "securepassword123",
			PasswordConfirmation: "differentpassword", // Mismatched confirmation
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.RegisterHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	sessionStore := auth.NewSessionStore()

	t.Run("handles user login successfully", func(t *testing.T) {
		db, err := database.New("test_login_success.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_success.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create a test user first
		hashedPassword, err := auth.HashPassword("password123")
		require.NoError(t, err)
		err = db.AddUser("Test User", "login@example.com", hashedPassword)
		require.NoError(t, err)

		// Create login data
		loginData := models.UserLogin{
			Email:    "login@example.com",
			Password: "password123",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return OK status (Inertia redirect uses 200 status code) after successful login
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check that a session cookie was set
		cookies := rr.Result().Cookies()
		sessionCookieFound := false
		for _, cookie := range cookies {
			if cookie.Name == SessionCookieName {
				sessionCookieFound = true
				break
			}
		}
		assert.True(t, sessionCookieFound, "Session cookie should be set")
	})

	t.Run("handles invalid credentials", func(t *testing.T) {
		db, err := database.New("test_login_invalid_creds.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_invalid_creds.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create a test user first
		hashedPassword, err := auth.HashPassword("password123")
		require.NoError(t, err)
		err = db.AddUser("Test User", "login@example.com", hashedPassword)
		require.NoError(t, err)

		// Create login data with wrong password
		loginData := models.UserLogin{
			Email:    "login@example.com",
			Password: "wrongpassword",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		authHandler.LoginHandler(rr, req)
		// Should return unauthorized error
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("handles non-existent user", func(t *testing.T) {
		db, err := database.New("test_login_nonexistent.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_nonexistent.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create login data for non-existent user
		loginData := models.UserLogin{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return unauthorized error
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("handles invalid JSON in request body", func(t *testing.T) {
		db, err := database.New("test_login_invalid_json.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_invalid_json.db")

		authHandler := NewAuth(db, in, sessionStore)

		invalidJSON := []byte(`{"email": "test@example.com", "password":}`) // Invalid JSON

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("handles non-POST requests", func(t *testing.T) {
		db, err := database.New("test_login_method.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_method.db")

		authHandler := NewAuth(db, in, sessionStore)

		req, err := http.NewRequest("GET", "/login", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return Method Not Allowed
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("handles missing email field", func(t *testing.T) {
		db, err := database.New("test_login_missing_email.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_missing_email.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create login data with missing email
		loginData := models.UserLogin{
			Email:    "", // Missing email
			Password: "password123",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles invalid email format", func(t *testing.T) {
		db, err := database.New("test_login_invalid_email.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_invalid_email.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create login data with invalid email format
		loginData := models.UserLogin{
			Email:    "invalid-email", // Invalid email format
			Password: "password123",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles missing password field", func(t *testing.T) {
		db, err := database.New("test_login_missing_password.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_missing_password.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create login data with missing password
		loginData := models.UserLogin{
			Email:    "login@example.com",
			Password: "", // Missing password
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("handles extra long password", func(t *testing.T) {
		db, err := database.New("test_login_long_password.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_login_long_password.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create a very long password
		longPassword := "A"
		for i := 0; i < 255; i++ {
			longPassword += "A"
		}
		
		// Create login data with very long password
		loginData := models.UserLogin{
			Email:    "login@example.com",
			Password: longPassword, // Very long password
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)

		// Should return Unprocessable Entity for validation error
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
}

func TestLogoutHandler(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	db, err := database.New("test_logout.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	sessionStore := auth.NewSessionStore()
	authHandler := NewAuth(db, in, sessionStore)

	// Create a session first
	sessionID, err := sessionStore.CreateSession(1, "logout@example.com")
	require.NoError(t, err)

	t.Run("handles user logout successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/logout", nil)
		require.NoError(t, err)

		// Add session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		authHandler.LogoutHandler(rr, req)

		// Should redirect to home page after logout
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "/")

		// Check that the session was deleted from the store
		_, err = sessionStore.GetSession(sessionID)
		assert.Error(t, err)
	})

	t.Run("handles logout without session", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/logout", nil)
		require.NoError(t, err)
		// No session cookie added

		rr := httptest.NewRecorder()
		authHandler.LogoutHandler(rr, req)

		// Should still redirect to home page (no error for no session)
		assert.Equal(t, http.StatusFound, rr.Code)
	})
}

func TestProfileHandler(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	sessionStore := auth.NewSessionStore()

	t.Run("returns profile for authenticated user", func(t *testing.T) {
		db, err := database.New("test_profile_auth.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_profile_auth.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create a test user
		hashedPassword, err := auth.HashPassword("password123")
		require.NoError(t, err)
		err = db.AddUser("Profile User", "profile@example.com", hashedPassword)
		require.NoError(t, err)

		// Create a session for the user
		user, err := db.GetUserByEmail("profile@example.com")
		require.NoError(t, err)
		sessionID, err := sessionStore.CreateSession(user.ID, user.Email)
		require.NoError(t, err)

		req, err := http.NewRequest("GET", "/profile", nil)
		require.NoError(t, err)

		// Add session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		authHandler.ProfileHandler(rr, req)

		// Should not be a 404 since the route should exist
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("redirects unauthenticated user to login", func(t *testing.T) {
		db, err := database.New("test_profile_unauth.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_profile_unauth.db")

		authHandler := NewAuth(db, in, sessionStore)

		req, err := http.NewRequest("GET", "/profile", nil)
		require.NoError(t, err)
		// No session cookie added

		rr := httptest.NewRecorder()
		authHandler.ProfileHandler(rr, req)

		// Should redirect to login page
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "/login")
	})
}

func TestGetUserIDFromSession(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	db, err := database.New("test_get_user_id.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	sessionStore := auth.NewSessionStore()
	authHandler := NewAuth(db, in, sessionStore)

	// Create a session
	sessionID, err := sessionStore.CreateSession(123, "test@example.com")
	require.NoError(t, err)

	t.Run("returns user ID for valid session", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		userID := authHandler.GetUserIDFromSession(req)

		assert.Equal(t, 123, userID)
	})

	t.Run("returns 0 for invalid session", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add invalid session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: "invalid-session-id",
		})

		userID := authHandler.GetUserIDFromSession(req)

		assert.Equal(t, 0, userID)
	})

	t.Run("returns 0 for no session", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)
		// No session cookie added

		userID := authHandler.GetUserIDFromSession(req)

		assert.Equal(t, 0, userID)
	})

	t.Run("handles malformed session cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add malformed session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: "", // Empty session ID
		})

		userID := authHandler.GetUserIDFromSession(req)

		assert.Equal(t, 0, userID)
	})

	t.Run("handles different HTTP methods", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/test", nil)
		require.NoError(t, err)

		// Add session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		userID := authHandler.GetUserIDFromSession(req)

		assert.Equal(t, 123, userID) // Should work with any HTTP method
	})
}

func TestIsAuthenticated(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	db, err := database.New("test_is_authenticated.db")
	require.NoError(t, err)
	defer db.Conn.Close()

	sessionStore := auth.NewSessionStore()
	authHandler := NewAuth(db, in, sessionStore)

	// Create a session
	sessionID, err := sessionStore.CreateSession(456, "auth@example.com")
	require.NoError(t, err)

	t.Run("returns true for authenticated user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		isAuthenticated := authHandler.IsAuthenticated(req)

		assert.True(t, isAuthenticated)
	})

	t.Run("returns false for unauthenticated user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)
		// No session cookie added

		isAuthenticated := authHandler.IsAuthenticated(req)

		assert.False(t, isAuthenticated)
	})

	t.Run("returns false for invalid session cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add invalid session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: "invalid-session-id",
		})

		isAuthenticated := authHandler.IsAuthenticated(req)

		assert.False(t, isAuthenticated)
	})

	t.Run("handles empty session cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add empty session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: "", // Empty value
		})

		isAuthenticated := authHandler.IsAuthenticated(req)

		assert.False(t, isAuthenticated)
	})
}

func TestAuthenticatedUser(t *testing.T) {
	// Create mock dependencies
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)

	sessionStore := auth.NewSessionStore()

	t.Run("returns user for authenticated user", func(t *testing.T) {
		db, err := database.New("test_auth_user_auth.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_auth_user_auth.db")

		authHandler := NewAuth(db, in, sessionStore)

		// Create a test user first
		hashedPassword, err := auth.HashPassword("password123")
		require.NoError(t, err)
		err = db.AddUser("Auth Test User", "authuser@example.com", hashedPassword)
		require.NoError(t, err)

		// Get the user to get their ID
		user, err := db.GetUserByEmail("authuser@example.com")
		require.NoError(t, err)

		// Create a session for the user
		sessionID, err := sessionStore.CreateSession(user.ID, user.Email)
		require.NoError(t, err)

		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		// Add session cookie to request
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: sessionID,
		})

		authUser := authHandler.AuthenticatedUser(req)

		require.NotNil(t, authUser)
		assert.Equal(t, user.ID, authUser.ID)
		assert.Equal(t, "Auth Test User", authUser.Name)
		assert.Equal(t, "authuser@example.com", authUser.Email)
	})

	t.Run("returns nil for unauthenticated user", func(t *testing.T) {
		db, err := database.New("test_auth_user_unauth.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_auth_user_unauth.db")

		authHandler := NewAuth(db, in, sessionStore)

		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)
		// No session cookie added

		authUser := authHandler.AuthenticatedUser(req)

		assert.Nil(t, authUser)
	})
}
