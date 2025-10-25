# Unit Tests for Svelte + Inertia + Go + Vite Application

This directory contains comprehensive unit tests for the entire application, covering all major components:

## Test Coverage

### 1. Models (`models/models_test.go`)
- Tests for `Country` struct fields and JSON serialization
- Tests for `NewCountry` struct fields and JSON serialization
- Verification of proper struct field definitions

### 2. Database (`database/db_test.go` and `database/integration_test.go`)
- **New function**: Tests database connection initialization and error handling
- **initDB function**: Tests table and index creation
- **GetRandomCountries**: Tests random country retrieval with proper flag emojis, empty databases, and error cases
- **GetAllCountries**: Tests all country retrieval with proper ordering and flag conversion
- **AddCountry**: Tests country insertion with timestamps and error handling
- **Integration tests**: Full CRUD flow testing to verify the entire data flow works correctly
- **Flag conversion**: Tests that country codes are properly converted to flag emojis

### 3. Handlers (`handlers/handlers_test.go`)
- **HomeHandler**: Tests home page rendering and error handling for nil Inertia
- **RandomCountriesHandler**: Tests random countries page rendering and database error handling
- **AllCountriesHandler**: Tests all countries page rendering and database error handling  
- **NewCountriesHandler**: Tests country creation, JSON validation, and database insertion
- **New function**: Tests handler initialization

### 4. Server (`server/server_test.go`)
- **New function**: Tests server initialization
- **SetupRoutes**: Tests route configuration
- **serverStaticFolder**: Tests static file serving setup

### 5. Main Application (`main_test.go`)
- Basic application initialization verification

## Key Features of the Test Suite

1. **Comprehensive Error Handling**: Tests for error scenarios like database connection failures, invalid JSON, and nil dependencies.

2. **Database Isolation**: Each test uses separate temporary databases to ensure isolation.

3. **Mock Dependencies**: Where appropriate, tests use mock templates and dependencies to isolate specific functionality.

4. **Integration Testing**: Includes integration tests that verify the complete data flow through the application.

5. **Edge Cases**: Tests for empty databases, constraint violations, and other edge cases.

## Test Execution

To run all tests:
```bash
go test ./...
```

To run with verbose output:
```bash
go test -v ./...
```

To run specific package tests:
```bash
go test ./database
go test ./handlers  
go test ./server
go test ./models
```

## Test-Driven Improvements Made

During test development, several improvements were made to the main application code:

1. **Error Handling**: Added nil checks for Inertia instances in all handlers to prevent panics
2. **Template Pathing**: Fixed filesystem pathing issues for compatibility with Go's fs.FS interface
3. **Robust Error Reporting**: Enhanced error messages and handling in database and handler functions

These tests ensure the application is robust, handles errors gracefully, and all components work together correctly.