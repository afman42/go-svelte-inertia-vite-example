// Package auth provides authentication and session management functionality.
// It includes password hashing, session store management, and user authentication utilities.
package auth

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hash), nil
}

// CheckPasswordHash compares a plain text password with a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Session represents a user session
type Session struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateSessionID creates a simple session ID (in production, use a more secure method)
func GenerateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// SessionStore is an in-memory store for sessions (in production, use Redis or database)
type SessionStore struct {
	sessions map[string]*Session
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

// CreateSession creates a new session for a user
func (s *SessionStore) CreateSession(userID int, email string) (string, error) {
	sessionID := GenerateSessionID()
	session := &Session{
		UserID:    userID,
		Email:     email,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Session expires in 24 hours
		CreatedAt: time.Now(),
	}
	s.sessions[sessionID] = session
	return sessionID, nil
}

// GetSession retrieves a session by ID
func (s *SessionStore) GetSession(sessionID string) (*Session, error) {
	session, exists := s.sessions[sessionID]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session not found or expired")
	}
	return session, nil
}

// DeleteSession removes a session
func (s *SessionStore) DeleteSession(sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

// ClearExpiredSessions removes expired sessions (should be run periodically)
func (s *SessionStore) ClearExpiredSessions() {
	now := time.Now()
	for sessionID, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, sessionID)
		}
	}
}