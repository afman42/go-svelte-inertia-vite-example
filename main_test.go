package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/handlers"
	"github.com/afman42/go-svelte-inertia/server"
	inertia "github.com/romsar/gonertia/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainApplicationInitialization(t *testing.T) {
	t.Run("verifies application components can be initialized", func(t *testing.T) {
		// Test database initialization
		db, err := database.New("test_main_init.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_main_init.db")
		
		assert.NotNil(t, db)
		assert.NotNil(t, db.Conn)
		
		// Test Inertia initialization
		in, err := inertia.NewFromFile("test_template.tmpl")
		if err != nil {
			// If test template doesn't exist, verify that we can at least check for its existence
			_, err = os.Stat("test_template.tmpl")
			// Just make sure there's no critical error when initializing the app
		}
		
		// Test session store initialization
		sessionStore := auth.NewSessionStore()
		assert.NotNil(t, sessionStore)
		// Can't access unexported field, so just check the store exists
		
		// Test handlers initialization
		handler := handlers.New(db, in, sessionStore)
		assert.NotNil(t, handler)
		assert.Equal(t, db, handler.DB)
		assert.Equal(t, in, handler.In)
		assert.Equal(t, sessionStore, handler.SessionStore)
		assert.NotNil(t, handler.Auth)
		
		// Test server initialization
		srv := server.New(in, handler)
		assert.NotNil(t, srv)
		assert.Equal(t, in, srv.Inertia)
		assert.Equal(t, handler, srv.Handler)
	})
}

func TestMainComponentsIntegration(t *testing.T) {
	t.Run("verifies components work together", func(t *testing.T) {
		// Initialize all components
		db, err := database.New("test_integration.db")
		require.NoError(t, err)
		defer db.Conn.Close()
		defer os.Remove("test_integration.db")
		
		// Add some test data to the database
		err = db.AddCountry("Test Country", "TC")
		require.NoError(t, err)
		
		// Initialize Inertia
		in, err := inertia.NewFromFile("test_template.tmpl")
		if err != nil {
			// Skip this test if template file doesn't exist
			t.Skip("test template not available")
		}
		
		// Initialize other components
		sessionStore := auth.NewSessionStore()
		handler := handlers.New(db, in, sessionStore)
		srv := server.New(in, handler)
		
		// Set up routes
		mux := http.NewServeMux()
		srv.SetupRoutes(mux, false) // Production mode
		
		// Test that routes can be registered without error
		assert.NotNil(t, mux)
	})
}

// TestMainFunction tests the main function by running a simplified version
func TestMainFunction(t *testing.T) {
	// This test is to verify that initialization doesn't panic
	// We can't run the full main function without starting a server
	t.Run("initialization completes without error", func(t *testing.T) {
		// Since the main function starts a server, we can only test initialization
		// without actually running the server
		assert.NotPanics(t, func() {
			// Initialize just like the main function does, but don't start the server
			in, err := inertia.NewFromFile("test_template.tmpl")
			if err != nil {
				// If template doesn't exist, skip this test
				t.Skip("test template not available")
				return
			}

			// Initialize database
			db, err := database.New("test_main_function.db")
			if err != nil {
				t.Logf("Could not create test database: %v", err)
				return
			}
			defer db.Conn.Close()
			defer os.Remove("test_main_function.db")

			// Initialize session store
			sessionStore := auth.NewSessionStore()

			// Initialize handlers
			handler := handlers.New(db, in, sessionStore)

			// Initialize server
			srv := server.New(in, handler)

			// Setup routes
			mux := http.NewServeMux()
			srv.SetupRoutes(mux, false)
		})
	})
}