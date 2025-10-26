package handlers

import (
	"net/http"

	"github.com/afman42/go-svelte-inertia/models"
	inertia "github.com/romsar/gonertia/v2"
)

// AddGlobalProps adds global properties (like user data) to props
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
