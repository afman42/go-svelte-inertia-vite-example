package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

	// Check if user already exists
	existingUser, _ := a.DB.GetUserByEmail(formData.Email)
	if existingUser != nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
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

	// Get user from the database
	user, err := a.DB.GetUserByEmail(formData.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if !auth.CheckPasswordHash(formData.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
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

	// Return success response
	// resp := map[string]interface{}{
	// 	"success": true,
	// 	"user": map[string]interface{}{
	// 		"id":    user.ID,
	// 		"name":  user.Name,
	// 		"email": user.Email,
	// 	},
	// }

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
	userID := a.GetUserIDFromSession(r)
	if userID == 0 {
		a.In.Redirect(w, r, "/login")
		return
	}

	user, err := a.DB.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		a.In.Redirect(w, r, "/login")
		return
	}

	err = a.In.Render(w, r, "Auth/Profile", inertia.Props{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
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
