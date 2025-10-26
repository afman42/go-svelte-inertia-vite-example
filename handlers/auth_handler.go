// Package handlers provides HTTP handlers for the application routes.
// It includes handlers for authentication, user management, and data operations.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/models"
	"github.com/afman42/go-svelte-inertia/validation"
	inertia "github.com/romsar/gonertia/v2"
)

const (
	SessionCookieName = "session_id"
)

// AuthHandler wraps the dependencies needed for authentication handlers
type AuthHandler struct {
	DB           *database.DB
	In           *inertia.Inertia
	SessionStore *auth.SessionStore
}

// NewAuth creates a new AuthHandler with the given dependencies
func NewAuth(db *database.DB, in *inertia.Inertia, sessionStore *auth.SessionStore) *AuthHandler {
	return &AuthHandler{
		DB:           db,
		In:           in,
		SessionStore: sessionStore,
	}
}

// RegisterHandler handles user registration
func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var formData models.UserRegister
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&formData)
	if err != nil {
		log.Printf("JSON decode error in RegisterHandler: %v", err)
		http.Error(w, "Bad Request - Invalid JSON", http.StatusBadRequest)
		return
	}

	v := validation.New()
	// Validate name
	v.Check(v.Required(formData.Name), "name", "Name is required")
	v.Check(v.MinLength(formData.Name, 2), "name", "Name must be at least 2 characters")

	// Validate email
	v.Check(v.Required(formData.Email), "email", "Email is required")
	v.ValidateEmail("email", formData.Email)

	// Validate password
	v.Check(v.Required(formData.Password), "password", "Password is required")
	v.Check(v.MinLength(formData.Password, 6), "password", "Password must be at least 6 characters")

	// Validate password confirmation
	v.Check(v.Required(formData.PasswordConfirmation), "password_confirmation", "Password Confirmation is required")
	v.Check(formData.Password != formData.PasswordConfirmation, "password_confirmation", "Passwords do not match")
	v.Check(v.MinLength(formData.PasswordConfirmation, 6), "password_confirmation", "Password Confirmation must be at least 6 characters")

	props := make(inertia.Props)
	if !v.Valid() {

		// Return validation errors to the frontend for HTML
		props["errors"] = v.GetErrors()
		props["old"] = formData // Include the submitted data to preserve form values

		err := a.In.Render(w, r, "Auth/Register", props)
		if err != nil {
			log.Printf("Error rendering Register page with errors: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Check if user already exists
	existingUser, _ := a.DB.GetUserByEmail(formData.Email)
	if existingUser != nil {

		// Return error to frontend for HTML
		v.AddError("email", "User with this email already exists")
		props["errors"] = v.GetErrors()
		props["old"] = formData

		err := a.In.Render(w, r, "Auth/Register", props)
		if err != nil {
			log.Printf("Error rendering Register page with duplicate email error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(formData.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the user in the database
	err = a.DB.AddUser(formData.Name, formData.Email, hashedPassword)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return success response - redirect to login for HTML requests
	a.In.Redirect(w, r, "/login")
}

// LoginHandler handles user login
func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var formData models.UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&formData)
	if err != nil {
		log.Printf("JSON decode error in LoginHandler: %v", err)
		http.Error(w, "Bad Request - Invalid JSON", http.StatusBadRequest)
		return
	}

	v := validation.New()

	// Validate email
	v.Check(v.Required(formData.Email), "email", "Email is required")
	v.ValidateEmail("email", formData.Email)
	v.Check(v.MaxLength(formData.Email, 255), "email", "Email must be at most 255 characters")

	// Validate password
	v.Check(v.Required(formData.Password), "password", "Password is required")
	v.Check(v.MinLength(formData.Password, 6), "password", "Password must be at least 6 characters")
	v.Check(v.MaxLength(formData.Password, 72), "password", "Password must be at most 72 characters")

	props := make(inertia.Props)
	if !v.Valid() {
		// Return validation errors to the frontend for HTML
		props["errors"] = v.GetErrors()
		props["old"] = formData // Include the submitted data to preserve form values

		err := a.In.Render(w, r, "Auth/Login", props)
		if err != nil {
			log.Printf("Error rendering Login page with errors: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Get user from the database
	user, err := a.DB.GetUserByEmail(formData.Email)
	if err != nil || user == nil {

		// Return error to frontend for HTML
		v.AddError("email", "Invalid check email")
		props["errors"] = v.GetErrors()
		props["old"] = formData

		err := a.In.Render(w, r, "Auth/Login", props)
		if err != nil {
			log.Printf("Error rendering Login page with invalid credentials error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Check password
	if !auth.CheckPasswordHash(formData.Password, user.Password) {

		// Return error to frontend for HTML
		v.AddError("password", "Invalid check password")
		props["errors"] = v.GetErrors()
		props["old"] = formData

		err := a.In.Render(w, r, "Auth/Login", props)
		if err != nil {
			log.Printf("Error rendering Login page with invalid password error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Create session
	sessionID, err := a.SessionStore.CreateSession(user.ID, user.Email)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	a.In.Redirect(w, r, "GET /", http.StatusOK)
}

// LogoutHandler handles user logout
func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		// No session, just redirect to home
		a.In.Redirect(w, r, "/")
		return
	}

	// Remove session from store
	sessionID := cookie.Value
	a.SessionStore.DeleteSession(sessionID)

	// Delete session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Redirect to home
	a.In.Redirect(w, r, "/")
}

// ProfileHandler handles user profile page
func (a *AuthHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if !a.IsAuthenticated(r) {
		// For regular HTML requests, redirect to login
		http.Redirect(w, r, "/login", http.StatusFound) // 302 redirect
		return
	}

	props := make(inertia.Props)
	props = a.AddGlobalProps(r, props)
	err := a.In.Render(w, r, "Auth/Profile", props)
	if err != nil {
		log.Printf("Error rendering profile page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetUserIDFromSession retrieves user ID from session cookie
func (a *AuthHandler) GetUserIDFromSession(r *http.Request) int {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return 0
	}

	session, err := a.SessionStore.GetSession(cookie.Value)
	if err != nil {
		return 0
	}

	return session.UserID
}

// IsAuthenticated checks if the user is authenticated
func (a *AuthHandler) IsAuthenticated(r *http.Request) bool {
	return a.GetUserIDFromSession(r) != 0
}

// AuthenticatedUser returns the authenticated user or nil if not authenticated
func (a *AuthHandler) AuthenticatedUser(r *http.Request) *models.User {
	userID := a.GetUserIDFromSession(r)
	if userID == 0 {
		return nil
	}

	user, err := a.DB.GetUserByID(userID)
	if err != nil {
		return nil
	}

	return user
}
