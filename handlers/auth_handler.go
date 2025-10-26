// Package handlers provides HTTP handlers for the application routes.
// It includes handlers for authentication, user management, and data operations.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/models"
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

	// Simple validation
	errors := make(map[string][]string)

	// Validate name
	if strings.TrimSpace(formData.Name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	} else if len(strings.TrimSpace(formData.Name)) < 2 {
		errors["name"] = append(errors["name"], "Name must be at least 2 characters")
	}

	// Validate email
	if strings.TrimSpace(formData.Email) == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else if !strings.Contains(formData.Email, "@") || !strings.Contains(formData.Email, ".") {
		// Simple email validation without regex
		errors["email"] = append(errors["email"], "Email format is invalid")
	}

	// Validate password
	if formData.Password == "" {
		errors["password"] = append(errors["password"], "Password is required")
	} else if len(formData.Password) < 6 {
		errors["password"] = append(errors["password"], "Password must be at least 6 characters")
	}

	// Validate password confirmation
	if formData.PasswordConfirmation == "" {
		errors["password_confirmation"] = append(errors["password_confirmation"], "Password confirmation is required")
	} else if formData.Password != formData.PasswordConfirmation {
		errors["password_confirmation"] = append(errors["password_confirmation"], "Passwords do not match")
	}

	props := make(inertia.Props)
	if len(errors) > 0 {
		// Return validation errors to the frontend
		props["errors"] = errors
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
		// Return error to frontend
		props["errors"] = map[string][]string{"email": {"User with this email already exists"}}
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

	// Return success response - redirect to login
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

	// Simple validation
	errors := make(map[string][]string)

	// Validate email
	if strings.TrimSpace(formData.Email) == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else if !strings.Contains(formData.Email, "@") || !strings.Contains(formData.Email, ".") {
		// Simple email validation
		errors["email"] = append(errors["email"], "Email format is invalid")
	}

	// Validate password
	if formData.Password == "" {
		errors["password"] = append(errors["password"], "Password is required")
	}

	props := make(inertia.Props)
	if len(errors) > 0 {
		// Return validation errors to the frontend
		props["errors"] = errors
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
		// Return error to frontend
		props["errors"] = map[string][]string{"email": {"Invalid email or password"}}
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
		// Return error to frontend
		props["errors"] = map[string][]string{"password": {"Invalid email or password"}}
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

func (a *AuthHandler) AddGlobalProps(r *http.Request, props inertia.Props) inertia.Props {
	// Add global user data
	user := a.AuthenticatedUser(r)
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
