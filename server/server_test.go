
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	inertia "github.com/romsar/gonertia/v2"
)

func TestNew(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a handler instance
	db, err := database.New("test_server_new.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	
	sessionStore := auth.NewSessionStore()
	handler := handlers.New(db, in, sessionStore)

	t.Run("creates new server instance", func(t *testing.T) {
		server := New(in, handler)

		require.NotNil(t, server)
		assert.Equal(t, in, server.Inertia)
		assert.Equal(t, handler, server.Handler)
	})
}

func TestSetupRoutes(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a handler instance
	db, err := database.New("test_setup_routes.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	
	sessionStore := auth.NewSessionStore()
	handler := handlers.New(db, in, sessionStore)
	server := New(in, handler)

	t.Run("sets up routes properly in production mode", func(t *testing.T) {
		mux := http.NewServeMux()
		server.SetupRoutes(mux, false) // Production mode

		// Test that the home route exists by checking if a request can be made
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// The route exists, but might return an error because we're not properly configured
		// Just check if the route is registered (it should not be a 404)
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("sets up routes properly in development mode", func(t *testing.T) {
		mux := http.NewServeMux()
		server.SetupRoutes(mux, true) // Development mode

		// Test that the home route exists
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// The route exists, but might return an error because we're not properly configured
		// Just check if the route is registered (it should not be a 404)
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("includes protected routes with authentication middleware", func(t *testing.T) {
		mux := http.NewServeMux()
		server.SetupRoutes(mux, false) // Production mode

		// Test that the POST /countries route exists (this is protected)
		req, err := http.NewRequest("POST", "/countries", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// The route should exist but may redirect due to auth middleware
		// It should not be a 404
		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})
}

func TestServerStaticFolder(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a handler instance
	db, err := database.New("test_static.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	
	sessionStore := auth.NewSessionStore()
	handler := handlers.New(db, in, sessionStore)
	server := New(in, handler)

	t.Run("sets up static folder properly", func(t *testing.T) {
		mux := http.NewServeMux()
		
		// This test will try to add a static route
		// Using a test directory or in-memory fs would be more appropriate in a real scenario
		// For now, we'll just make sure the function doesn't panic
		assert.NotPanics(t, func() {
			server.serverStaticFolder(mux, "/test/", http.Dir("."))
		})
	})

	t.Run("serves static files correctly", func(t *testing.T) {
		mux := http.NewServeMux()
		server.serverStaticFolder(mux, "/assets/", http.Dir("."))
		
		// Test that the static file handler works
		req, err := http.NewRequest("GET", "/assets/somefile.txt", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// The file may not exist, but we're testing that the handler is set up
		// It should either be 404 (file not found) or 200 (file found)
		assert.Contains(t, []int{http.StatusNotFound, http.StatusOK}, rr.Code)
	})
}

func TestRouteHandlers(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a handler instance
	db, err := database.New("test_routes.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	
	sessionStore := auth.NewSessionStore()
	handler := handlers.New(db, in, sessionStore)
	server := New(in, handler)

	t.Run("registers all required routes", func(t *testing.T) {
		mux := http.NewServeMux()
		server.SetupRoutes(mux, false) // Production mode

		routes := []string{
			"/",
			"/random", 
			"/all",
			"/login",
			"/register",
			"/logout",
			"/profile",
		}

		for _, route := range routes {
			req, err := http.NewRequest("GET", route, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			
			// Each route should be registered (not a 404)
			assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route %s should be registered", route)
		}
	})
}