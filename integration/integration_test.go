package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/handlers"
	"github.com/afman42/go-svelte-inertia/models"
	"github.com/afman42/go-svelte-inertia/server"
	inertia "github.com/romsar/gonertia/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullApplicationFlow(t *testing.T) {
	// This test verifies the full application flow from database to handler to server
	db, err := database.New("test_full_flow.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	defer os.Remove("test_full_flow.db")
	
	// Initialize Inertia with a simple template for testing
	in, err := inertia.NewFromFile("../test_template.tmpl")
	if err != nil {
		t.Skip("test template not available")
		return
	}
	
	sessionStore := auth.NewSessionStore()
	handler := handlers.New(db, in, sessionStore)
	srv := server.New(in, handler)
	
	// Set up the HTTP multiplexer with routes
	mux := http.NewServeMux()
	srv.SetupRoutes(mux, false) // Production mode

	t.Run("full user registration and login flow", func(t *testing.T) {
		// Test registration
		registerData := models.UserRegister{
			Name:     "Integration User",
			Email:    "integration@example.com",
			Password: "securepassword123",
		}
		
		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Registration should return status 200
		assert.Equal(t, http.StatusOK, rr.Code)

		// Test login
		loginData := models.UserLogin{
			Email:    "integration@example.com",
			Password: "securepassword123",
		}
		
		jsonData, err = json.Marshal(loginData)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Should return OK status (Inertia redirect uses 200 status code) after successful login
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// Check for session cookie
		sessionCookieFound := false
		for _, cookie := range rr.Result().Cookies() {
			if cookie.Name == handlers.SessionCookieName {
				sessionCookieFound = true
				break
			}
		}
		assert.True(t, sessionCookieFound, "Session cookie should be set after login")
	})

	t.Run("protected route access after authentication", func(t *testing.T) {
		// First login to get a session
		loginData := models.UserLogin{
			Email:    "integration@example.com",
			Password: "securepassword123",
		}
		
		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Extract session cookie
		var sessionCookie *http.Cookie
		for _, cookie := range rr.Result().Cookies() {
			if cookie.Name == handlers.SessionCookieName {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie, "Session cookie should exist")
		
		// Now try to access a protected route (POST /countries) with the session
		countryData := models.NewCountry{
			Name: "Integration Country",
			Code: "IC",
		}
		
		jsonData, err = json.Marshal(countryData)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/countries", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie) // Add the session cookie
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Should redirect after successful country creation
		assert.Equal(t, http.StatusFound, rr.Code)
		
		// Verify the country was added to the database
		countries, err := db.GetAllCountries()
		require.NoError(t, err)
		
		found := false
		for _, country := range countries {
			if country.Name == "Integration Country" {
				found = true
				break
			}
		}
		assert.True(t, found, "Country should have been added to the database")
	})

	t.Run("public routes accessible without authentication", func(t *testing.T) {
		// Test home route
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Should be accessible (not a 404)
		assert.NotEqual(t, http.StatusNotFound, rr.Code)

		// Test random countries route
		req, err = http.NewRequest("GET", "/random", nil)
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Should be accessible (not a 404)
		assert.NotEqual(t, http.StatusNotFound, rr.Code)

		// Test all countries route
		req, err = http.NewRequest("GET", "/all", nil)
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Should be accessible (not a 404)
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("logout functionality", func(t *testing.T) {
		// First login to get a session
		loginData := models.UserLogin{
			Email:    "integration@example.com",
			Password: "securepassword123",
		}
		
		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Extract session cookie
		var sessionCookie *http.Cookie
		for _, cookie := range rr.Result().Cookies() {
			if cookie.Name == handlers.SessionCookieName {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie, "Session cookie should exist")
		
		// Now try to logout
		req, err = http.NewRequest("GET", "/logout", nil)
		req.AddCookie(sessionCookie) // Add the session cookie
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// Should redirect after logout
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "/")
	})
}

func TestApplicationInitializationIntegration(t *testing.T) {
	db, err := database.New("test_init_integration.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	defer os.Remove("test_init_integration.db")
	
	// Initialize Inertia with a simple template
	in, err := inertia.NewFromFile("../test_template.tmpl")
	if err != nil {
		t.Skip("test template not available")
		return
	}
	
	sessionStore := auth.NewSessionStore()
	handler := handlers.New(db, in, sessionStore)
	srv := server.New(in, handler)
	
	// Set up routes
	mux := http.NewServeMux()
	srv.SetupRoutes(mux, false) // Production mode

	// Test that all expected routes are registered
	endpoints := []string{
		"/",
		"/random",
		"/all",
		"/login",
		"/register",
		"/logout",
		"/profile",
	}
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("GET", endpoint, nil)
		require.NoError(t, err, "Failed to create request for endpoint: %s", endpoint)

		// Create a test recorder to check if it's a 404 or not
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		
		// All endpoints should be registered (not 404)
		assert.NotEqualf(t, http.StatusNotFound, rr.Code, "Endpoint %s should be registered", endpoint)
	}
}

func TestDatabaseHandlerIntegration(t *testing.T) {
	db, err := database.New("test_db_handler_integration.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	defer os.Remove("test_db_handler_integration.db")
	
	// Add some test data to the database
	err = db.AddCountry("Test Integration Country", "TIC")
	require.NoError(t, err)
	
	// Test that the data can be retrieved
	countries, err := db.GetAllCountries()
	require.NoError(t, err)
	assert.Len(t, countries, 1)
	assert.Equal(t, "Test Integration Country", countries[0].Name)
}

func TestAuthIntegration(t *testing.T) {
	db, err := database.New("test_auth_integration.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	defer os.Remove("test_auth_integration.db")
	
	// Initialize Inertia
	in, err := inertia.NewFromFile("../test_template.tmpl")
	if err != nil {
		t.Skip("test template not available")
		return
	}
	
	sessionStore := auth.NewSessionStore()
	authHandler := handlers.NewAuth(db, in, sessionStore)
	
	// Test the full authentication flow
	email := "auth.test@example.com"
	password := "testpassword123"
	
	// Register a user
	userRegister := models.UserRegister{
		Name:     "Auth Test User",
		Email:    email,
		Password: password,
	}
	
	hashedPassword, err := auth.HashPassword(userRegister.Password)
	require.NoError(t, err)
	
	err = db.AddUser(userRegister.Name, userRegister.Email, hashedPassword)
	require.NoError(t, err)
	
	// Verify user was added
	user, err := db.GetUserByEmail(email)
	require.NoError(t, err)
	assert.Equal(t, userRegister.Name, user.Name)
	assert.Equal(t, userRegister.Email, user.Email)
	
	// Test creating a session
	sessionID, err := sessionStore.CreateSession(user.ID, user.Email)
	require.NoError(t, err)
	assert.NotEmpty(t, sessionID)
	
	// Test getting the user from session
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	
	// Add the session cookie to the request
	req.AddCookie(&http.Cookie{
		Name:  handlers.SessionCookieName,
		Value: sessionID,
	})
	
	userID := authHandler.GetUserIDFromSession(req)
	assert.Equal(t, user.ID, userID)
	
	isAuthenticated := authHandler.IsAuthenticated(req)
	assert.True(t, isAuthenticated)
	
	authUser := authHandler.AuthenticatedUser(req)
	assert.NotNil(t, authUser)
	assert.Equal(t, user.ID, authUser.ID)
	assert.Equal(t, user.Name, authUser.Name)
	assert.Equal(t, user.Email, authUser.Email)
}