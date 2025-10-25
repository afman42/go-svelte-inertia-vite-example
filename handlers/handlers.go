package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/models"
	inertia "github.com/romsar/gonertia/v2"
)

// Handler wraps the dependencies needed for handlers
type Handler struct {
	DB  *database.DB
	In  *inertia.Inertia
}

// New creates a new Handler with the given dependencies
func New(db *database.DB, in *inertia.Inertia) *Handler {
	return &Handler{
		DB:  db,
		In:  in,
	}
}

// HomeHandler handles the home page request
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if h.In == nil {
		log.Printf("Inertia is nil in HomeHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	time.Sleep(300 * time.Millisecond)
	err := h.In.Render(w, r, "Home", inertia.Props{
		"user": "data",
	})
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

	err = h.In.Render(w, r, "Countries/Random", inertia.Props{
		"countries": countries,
	})
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

	err = h.In.Render(w, r, "Countries/All", inertia.Props{
		"countries": countries,
	})
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
	
	decoder := json.NewDecoder(r.Body)
	var formData models.NewCountry

	err := decoder.Decode(&formData)
	if err != nil {
		log.Printf("JSON decode error in NewCountriesHandler: %v", err)
		http.Error(w, "Bad Request - Invalid JSON", http.StatusBadRequest)
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