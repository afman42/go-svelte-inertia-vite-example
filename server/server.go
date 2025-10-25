package server

import (
	"net/http"

	"github.com/afman42/go-svelte-inertia/handlers"
	inertia "github.com/romsar/gonertia/v2"
)

// Server represents the HTTP server with its dependencies
type Server struct {
	Inertia *inertia.Inertia
	Handler *handlers.Handler
}

// New creates a new Server instance
func New(in *inertia.Inertia, h *handlers.Handler) *Server {
	return &Server{
		Inertia: in,
		Handler: h,
	}
}

// SetupRoutes configures the HTTP routes
func (s *Server) SetupRoutes(mux *http.ServeMux, isDev bool) {
	// Setup static file serving
	if isDev {
		s.serverStaticFolder(mux, "/src/assets/", http.Dir("frontend/src/assets"))
	} else {
		s.serverStaticFolder(mux, "/assets/", http.Dir("frontend/dist/assets"))
	}

	// Register endpoints with middleware
	endpoints := map[string]http.HandlerFunc{
		"/{$}":            s.Handler.HomeHandler,
		"/random":         s.Handler.RandomCountriesHandler,
		"/all":            s.Handler.AllCountriesHandler,
		"POST /countries": s.Handler.NewCountriesHandler,
	}

	for endpoint, f := range endpoints {
		mux.Handle(endpoint, s.Inertia.Middleware(http.HandlerFunc(f)))
	}
}

// serverStaticFolder sets up static file serving for a specific path
func (s *Server) serverStaticFolder(mux *http.ServeMux, path string, fsys http.FileSystem) {
	mux.Handle(path, http.StripPrefix(path, http.FileServer(fsys)))
}