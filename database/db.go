package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/afman42/go-svelte-inertia/models"
	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB connection
type DB struct {
	Conn *sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &DB{Conn: db}
	
	if err = database.initDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return database, nil
}

// initDB creates the countries table if it doesn't exist
func (db *DB) initDB() error {
	// Create the countries table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS countries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		alpha2 TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create countries table: %w", err)
	}

	// Add an index on created_at for better performance for the all countries query
	indexQuery := `CREATE INDEX IF NOT EXISTS idx_countries_created_at ON countries (created_at DESC);`
	_, err = db.Conn.Exec(indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on countries table: %w", err)
	}

	return nil
}

// GetRandomCountries returns 10 random countries
func (db *DB) GetRandomCountries() ([]models.Country, error) {
	rows, err := db.Conn.Query("SELECT name, alpha2 FROM countries order by random() limit 10")
	if err != nil {
		return nil, fmt.Errorf("database query error in GetRandomCountries: %w", err)
	}
	defer rows.Close()

	countries := make([]models.Country, 0, 10)

	for rows.Next() {
		country := models.Country{}
		var alpha2 string
		if err := rows.Scan(&country.Name, &alpha2); err != nil {
			return nil, fmt.Errorf("database scan error in GetRandomCountries: %w", err)
		}

		country.Flag = country2flag(alpha2)
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database rows error in GetRandomCountries: %w", err)
	}

	return countries, nil
}

// GetAllCountries returns all countries ordered by creation date and name
func (db *DB) GetAllCountries() ([]models.Country, error) {
	rows, err := db.Conn.Query("SELECT name, alpha2 FROM countries order by created_at desc, name asc")
	if err != nil {
		return nil, fmt.Errorf("database query error in GetAllCountries: %w", err)
	}
	defer rows.Close()

	countries := make([]models.Country, 0)

	for rows.Next() {
		country := models.Country{}
		var alpha2 string
		if err := rows.Scan(&country.Name, &alpha2); err != nil {
			return nil, fmt.Errorf("database scan error in GetAllCountries: %w", err)
		}

		country.Flag = country2flag(alpha2)
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database rows error in GetAllCountries: %w", err)
	}

	return countries, nil
}

// AddCountry inserts a new country into the database
func (db *DB) AddCountry(name, code string) error {
	_, err := db.Conn.Exec("insert into countries (name, alpha2, created_at) values (?, ?, datetime())", name, code)
	if err != nil {
		return fmt.Errorf("database insert error in AddCountry: %w", err)
	}

	return nil
}

// country2flag converts a country code to a flag emoji
func country2flag(countryCode string) string {
	var flagEmoji strings.Builder
	countryCode = strings.ToUpper(countryCode)
	for _, char := range countryCode {
		flagEmoji.WriteRune(rune(char) + 0x1F1A5)
	}
	return flagEmoji.String()
}