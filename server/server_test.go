
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	
	handler := handlers.New(db, in)

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
	
	handler := handlers.New(db, in)
	server := New(in, handler)

	t.Run("sets up routes properly", func(t *testing.T) {
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
}

func TestServerStaticFolder(t *testing.T) {
	// Create a mock Inertia instance
	in, err := inertia.NewFromFile("../test_template.tmpl")
	require.NoError(t, err)
	
	// Create a handler instance
	db, err := database.New("test_static.db")
	require.NoError(t, err)
	defer db.Conn.Close()
	
	handler := handlers.New(db, in)
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
}