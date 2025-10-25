package models

import "time"

// Country represents a country with name and flag
type Country struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// NewCountry represents form data for creating a new country
type NewCountry struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// User represents a user with authentication details
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserLogin represents login form data
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegister represents registration form data
type UserRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}