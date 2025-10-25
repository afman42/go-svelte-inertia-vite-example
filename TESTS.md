# Unit Tests for Svelte + Inertia + Go + Vite Application

This directory contains comprehensive unit tests for the entire application, covering all major components:

## Test Coverage

### 1. Models (`models/models_test.go`)
- Tests for `Country` struct fields and JSON serialization
- Tests for `NewCountry` struct fields and JSON serialization
- Tests for `User` struct fields and JSON serialization
- Tests for `UserLogin` and `UserRegister` struct fields
- Verification of proper struct field definitions

### 2. Authentication (`auth/auth_test.go`)
- **HashPassword**: Tests password hashing functionality and error handling
- **CheckPasswordHash**: Tests password verification against hashed passwords
- **GenerateSessionID**: Tests session ID generation
- **NewSessionStore**: Tests session store initialization
- **CreateSession**: Tests session creation with proper expiration times
- **GetSession**: Tests session retrieval with validity checks
- **DeleteSession**: Tests session deletion
- **ClearExpiredSessions**: Tests cleanup of expired sessions

### 3. Database (`database/db_test.go`)
- **New function**: Tests database connection initialization and error handling
- **initDB function**: Tests table and index creation, proper table structure validation
- **GetRandomCountries**: Tests random country retrieval with proper flag emojis, empty databases, limits, and error cases
- **GetAllCountries**: Tests all country retrieval with proper ordering and flag conversion
- **AddCountry**: Tests country insertion with timestamps, constraint handling, and error cases
- **Table structure**: Tests for proper column definitions and indexes

### 4. Handlers (`handlers/handlers_test.go`)
- **New function**: Tests handler initialization with all dependencies
- **HomeHandler**: Tests home page rendering and error handling for nil Inertia
- **RandomCountriesHandler**: Tests random countries page rendering and database error handling
- **AllCountriesHandler**: Tests all countries page rendering and database error handling  
- **NewCountriesHandler**: Tests country creation, JSON validation, field validation, and database insertion
- **AuthMiddleware**: Tests authentication middleware functionality
- **Auth handlers**: Tests register, login, logout, profile, and authentication functions
- **IsAuthenticated**: Tests authentication status checking
- **AuthenticatedUser**: Tests retrieval of authenticated user information
- **GetUserIDFromSession**: Tests user ID retrieval from session

### 5. Server (`server/server_test.go`)
- **New function**: Tests server initialization
- **SetupRoutes**: Tests route configuration for both dev and prod modes
- **serverStaticFolder**: Tests static file serving setup
- **Route registration**: Tests that all routes are properly registered
- **Protected routes**: Tests that protected routes are configured with auth middleware

### 6. Main Application (`main_test.go`)
- **Initialization**: Tests application component initialization
- **Integration**: Tests that components work together properly
- **Main function**: Tests main function execution without panicking

### 7. Integration Tests (`integration_test.go`)
- **Full user registration and login flow**: Tests complete auth workflow
- **Protected route access**: Tests that authenticated users can access protected routes
- **Public routes**: Tests that public routes are accessible without auth
- **Static file serving**: Tests static file serving functionality
- **Logout functionality**: Tests complete logout workflow
- **Application initialization**: Tests that all routes are properly registered
- **Database-handler integration**: Tests data flow between database and handlers
- **Authentication integration**: Tests full auth flow with session management

## Key Features of the Test Suite

1. **Comprehensive Error Handling**: Tests for error scenarios like database connection failures, invalid JSON, nil dependencies, and authentication failures.

2. **Database Isolation**: Each test uses separate temporary databases to ensure isolation.

3. **Mock Dependencies**: Where appropriate, tests use mock templates and dependencies to isolate specific functionality.

4. **Integration Testing**: Includes integration tests that verify the complete data flow through the application.

5. **Edge Cases**: Tests for empty databases, constraint violations, expired sessions, and other edge cases.

6. **Authentication Flow Testing**: Full testing of user registration, login, logout, and profile access workflows.

7. **Protected Route Testing**: Verification that protected routes require authentication.

8. **Session Management**: Tests for proper session creation, validation, and cleanup.

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
go test ./auth
go test ./models
```

To run integration tests specifically:
```bash
go test -run Integration ./integration_test.go
```

## Test-Driven Improvements Made

During test development, several improvements were made to the main application code:

1. **Error Handling**: Added nil checks for Inertia instances in all handlers to prevent panics
2. **Template Pathing**: Fixed filesystem pathing issues for compatibility with Go's fs.FS interface
3. **Robust Error Reporting**: Enhanced error messages and handling in database and handler functions
4. **Authentication Security**: Improved password hashing and session management practices
5. **Session Validation**: Added proper session expiration and validation checks
6. **Input Validation**: Enhanced validation for user registration and login forms

These tests ensure the application is robust, handles errors gracefully, authentication is secure, and all components work together correctly.