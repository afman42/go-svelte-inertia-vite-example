package models

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