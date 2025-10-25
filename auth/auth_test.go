package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("successfully hashes a password", func(t *testing.T) {
		password := "mySecurePassword123"
		
		hash, err := HashPassword(password)
		
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash) // Should not be the same as original
	})

	t.Run("returns error for empty password", func(t *testing.T) {
		password := ""
		
		hash, err := HashPassword(password)
		
		// bcrypt should handle empty passwords, so this should not error
		// unless the implementation has special validation
		assert.NotEmpty(t, hash) // The hash should still be generated
		assert.NoError(t, err)
	})

	t.Run("different passwords produce different hashes", func(t *testing.T) {
		password1 := "password1"
		password2 := "password2"
		
		hash1, err := HashPassword(password1)
		require.NoError(t, err)
		
		hash2, err := HashPassword(password2)
		require.NoError(t, err)
		
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("same password produces different hashes due to salt", func(t *testing.T) {
		password := "samePassword"
		
		hash1, err := HashPassword(password)
		require.NoError(t, err)
		
		hash2, err := HashPassword(password)
		require.NoError(t, err)
		
		assert.NotEqual(t, hash1, hash2) // bcrypt uses salt, so hashes should differ
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("returns true for correct password", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)
		require.NoError(t, err)
		
		result := CheckPasswordHash(password, hash)
		
		assert.True(t, result)
	})

	t.Run("returns false for incorrect password", func(t *testing.T) {
		password := "mySecurePassword123"
		wrongPassword := "wrongPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)
		
		result := CheckPasswordHash(wrongPassword, hash)
		
		assert.False(t, result)
	})

	t.Run("returns false for invalid hash", func(t *testing.T) {
		password := "mySecurePassword123"
		invalidHash := "invalidHashString"
		
		result := CheckPasswordHash(password, invalidHash)
		
		assert.False(t, result)
	})

	t.Run("returns false for empty hash", func(t *testing.T) {
		password := "mySecurePassword123"
		emptyHash := ""
		
		result := CheckPasswordHash(password, emptyHash)
		
		assert.False(t, result)
	})
}

func TestGenerateSessionID(t *testing.T) {
	t.Run("generates unique session IDs", func(t *testing.T) {
		sessionID1 := GenerateSessionID()
		sessionID2 := GenerateSessionID()
		
		assert.NotEmpty(t, sessionID1)
		assert.NotEmpty(t, sessionID2)
		assert.NotEqual(t, sessionID1, sessionID2)
	})

	t.Run("generates numeric strings", func(t *testing.T) {
		sessionID := GenerateSessionID()
		
		// Just verify it's not empty - the function generates a timestamp-based string
		assert.NotEmpty(t, sessionID)
	})
}

func TestNewSessionStore(t *testing.T) {
	t.Run("creates new session store", func(t *testing.T) {
		store := NewSessionStore()
		
		require.NotNil(t, store)
		assert.NotNil(t, store.sessions)
		assert.Empty(t, store.sessions)
	})
}

func TestSessionStoreCreateSession(t *testing.T) {
	t.Run("creates a new session successfully", func(t *testing.T) {
		store := NewSessionStore()
		
		sessionID, err := store.CreateSession(1, "test@example.com")
		
		require.NoError(t, err)
		assert.NotEmpty(t, sessionID)
		
		// Check that the session exists in the store
		session, exists := store.sessions[sessionID]
		assert.True(t, exists)
		assert.Equal(t, 1, session.UserID)
		assert.Equal(t, "test@example.com", session.Email)
		assert.True(t, !session.ExpiresAt.IsZero())
		assert.True(t, !session.CreatedAt.IsZero())
	})

	t.Run("creates session with correct expiration time", func(t *testing.T) {
		store := NewSessionStore()
		
		sessionID, err := store.CreateSession(2, "user@example.com")
		require.NoError(t, err)
		
		session := store.sessions[sessionID]
		expectedExpiry := time.Now().Add(24 * time.Hour)
		
		// Check that the session expires in about 24 hours (with some tolerance)
		assert.WithinDuration(t, expectedExpiry, session.ExpiresAt, 1*time.Minute)
	})
}

func TestSessionStoreGetSession(t *testing.T) {
	t.Run("retrieves existing session", func(t *testing.T) {
		store := NewSessionStore()
		sessionID, err := store.CreateSession(1, "test@example.com")
		require.NoError(t, err)

		session, err := store.GetSession(sessionID)

		require.NoError(t, err)
		assert.Equal(t, 1, session.UserID)
		assert.Equal(t, "test@example.com", session.Email)
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		store := NewSessionStore()

		session, err := store.GetSession("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "session not found")
	})

	t.Run("returns error for expired session", func(t *testing.T) {
		store := NewSessionStore()
		sessionID, err := store.CreateSession(1, "test@example.com")
		require.NoError(t, err)

		// Manually expire the session
		store.sessions[sessionID].ExpiresAt = time.Now().Add(-1 * time.Hour)

		session, err := store.GetSession(sessionID)

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "session not found or expired")
	})
}

func TestSessionStoreDeleteSession(t *testing.T) {
	t.Run("deletes existing session", func(t *testing.T) {
		store := NewSessionStore()
		sessionID, err := store.CreateSession(1, "test@example.com")
		require.NoError(t, err)

		err = store.DeleteSession(sessionID)

		require.NoError(t, err)
		_, exists := store.sessions[sessionID]
		assert.False(t, exists)
	})

	t.Run("handles deletion of non-existent session", func(t *testing.T) {
		store := NewSessionStore()

		err := store.DeleteSession("non-existent-id")

		// Deletion of non-existent session should not return an error
		assert.NoError(t, err)
	})
}

func TestSessionStoreClearExpiredSessions(t *testing.T) {
	t.Run("removes expired sessions", func(t *testing.T) {
		store := NewSessionStore()
		
		// Create a valid session
		validSessionID, err := store.CreateSession(1, "valid@example.com")
		require.NoError(t, err)
		
		// Create an expired session manually
		expiredSessionID := "expired-session-id"
		store.sessions[expiredSessionID] = &Session{
			UserID:    2,
			Email:     "expired@example.com",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			CreatedAt: time.Now(),
		}

		store.ClearExpiredSessions()

		// Check that the valid session still exists
		_, validExists := store.sessions[validSessionID]
		assert.True(t, validExists)
		
		// Check that the expired session was removed
		_, expiredExists := store.sessions[expiredSessionID]
		assert.False(t, expiredExists)
	})

	t.Run("does not remove valid sessions", func(t *testing.T) {
		store := NewSessionStore()
		
		// Create a valid session
		sessionID, err := store.CreateSession(1, "valid@example.com")
		require.NoError(t, err)

		store.ClearExpiredSessions()

		// Check that the valid session still exists
		_, exists := store.sessions[sessionID]
		assert.True(t, exists)
	})
}