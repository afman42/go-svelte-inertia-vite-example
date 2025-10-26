// Package handlers provides HTTP handlers for the application routes.
// It includes handlers for authentication, user management, and data operations.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/models"
	inertia "github.com/romsar/gonertia/v2"
)

// Handler wraps the dependencies needed for handlers
type Handler struct {
	DB           *database.DB
	In           *inertia.Inertia
	SessionStore *auth.SessionStore
	Auth         *AuthHandler
}

// New creates a new Handler with the given dependencies
func New(db *database.DB, in *inertia.Inertia, sessionStore *auth.SessionStore) *Handler {
	authHandler := NewAuth(db, in, sessionStore)
	return &Handler{
		DB:           db,
		In:           in,
		SessionStore: sessionStore,
		Auth:         authHandler,
	}
}

// HomeHandler handles the home page request
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in HomeHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	props := make(inertia.Props)
	props = h.AddGlobalProps(r, props)

	time.Sleep(300 * time.Millisecond)

	err := h.In.Render(w, r, "Home", props)
	if err != nil {
		log.Printf("Error rendering Home page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RandomCountriesHandler handles the random countries request
func (h *Handler) RandomCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in RandomCountriesHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	countries, err := h.DB.GetRandomCountries()
	if err != nil {
		log.Printf("Error getting random countries: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	props := make(inertia.Props)
	props = h.AddGlobalProps(r, props)
	props["countries"] = countries

	err = h.In.Render(w, r, "Countries/Random", props)
	if err != nil {
		log.Printf("Error rendering random countries page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AllCountriesHandler handles the all countries request
func (h *Handler) AllCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in AllCountriesHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	countries, err := h.DB.GetAllCountries()
	if err != nil {
		log.Printf("Error getting all countries: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	props := make(inertia.Props)
	props = h.AddGlobalProps(r, props)
	props["countries"] = countries

	err = h.In.Render(w, r, "Countries/All", props)
	if err != nil {
		log.Printf("Error rendering all countries page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// NewCountriesHandler handles adding a new country
func (h *Handler) NewCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in NewCountriesHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user is authenticated
	userID := h.Auth.GetUserIDFromSession(r)
	if userID == 0 {
		h.In.Redirect(w, r, "/login")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var formData models.NewCountry

	err := decoder.Decode(&formData)
	if err != nil {
		log.Printf("JSON decode error in NewCountriesHandler: %v", err)
		http.Error(w, "Bad Request - Invalid JSON", http.StatusBadRequest)
		return
	}

	// Simple validation
	errors := make(map[string][]string)

	// Validate name
	if strings.TrimSpace(formData.Name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	} else if len(strings.TrimSpace(formData.Name)) < 2 {
		errors["name"] = append(errors["name"], "Name must be at least 2 characters")
	}

	// Validate code
	if strings.TrimSpace(formData.Code) == "" {
		errors["code"] = append(errors["code"], "Country code is required")
	} else if len(strings.TrimSpace(formData.Code)) != 2 {
		errors["code"] = append(errors["code"], "Country code must be 2 characters")
	} else if !isAlphaString(formData.Code) {
		errors["code"] = append(errors["code"], "Country code must contain only letters")
	}

	props := make(inertia.Props)
	props = h.AddGlobalProps(r, props)
	if len(errors) > 0 {

		// Return validation errors to the frontend
		props["errors"] = errors
		countries, err := h.DB.GetAllCountries()
		if err != nil {
			log.Printf("Error getting all countries: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		props["countries"] = countries
		err = h.In.Render(w, r, "Countries/All", props)
		if err != nil {
			log.Printf("Error rendering Countries/All page with errors: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = h.DB.AddCountry(formData.Name, formData.Code)
	if err != nil {
		log.Printf("Database insert error in NewCountriesHandler: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.In.Redirect(w, r, "/all")
}

// isAlphaString checks if a string contains only alphabetic characters
func isAlphaString(s string) bool {
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return true
}

// LoginViewHandler shows the login page
func (h *Handler) LoginViewHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in LoginViewHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user is already logged in
	if h.Auth.IsAuthenticated(r) {
		h.In.Redirect(w, r, "/")
		return
	}

	props := make(inertia.Props)
	props = h.AddGlobalProps(r, props)

	err := h.In.Render(w, r, "Auth/Login", props)
	if err != nil {
		log.Printf("Error rendering Login page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RegisterViewHandler shows the register page
func (h *Handler) RegisterViewHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in RegisterViewHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user is already logged in
	if h.Auth.IsAuthenticated(r) {
		h.In.Redirect(w, r, "/")
		return
	}

	props := make(inertia.Props)
	props = h.AddGlobalProps(r, props)
	err := h.In.Render(w, r, "Auth/Register", props)
	if err != nil {
		log.Printf("Error rendering Register page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AuthMiddleware is a middleware that checks if user is authenticated
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.Auth.IsAuthenticated(r) {
			h.In.Redirect(w, r, "/login")
			return
		}
		next.ServeHTTP(w, r)
	}
}

// AddGlobalProps adds global properties (like user data) to props
func (h *Handler) AddGlobalProps(r *http.Request, props inertia.Props) inertia.Props {
	// Add global user data
	user := h.Auth.AuthenticatedUser(r)
	if user != nil {
		props["user"] = &models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	} else {
		props["user"] = nil
	}

	return props
}
