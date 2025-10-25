package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/afman42/go-svelte-inertia/auth"
	"github.com/afman42/go-svelte-inertia/database"
	"github.com/afman42/go-svelte-inertia/handlers"
	"github.com/afman42/go-svelte-inertia/server"
	"github.com/olivere/vite"
	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	isDev := flag.Bool("dev", false, "run in development mode")
	flag.Parse()

	// Setup Vite integration
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:           os.DirFS("frontend/dist"),
		IsDev:        *isDev,
		ViteURL:      "http://localhost:5174",
		ViteEntry:    "src/main.ts",
		ViteTemplate: 14,
	})
	if err != nil {
		log.Fatal("Failed to create Vite HTML fragment:", err)
		return
	}

	// Initialize Inertia
	in, err := inertia.NewFromFile("frontend/index.tmpl")
	if err != nil {
		log.Fatal("Failed to initialize Inertia:", err)
		return
	}

	in.ShareTemplateData("Vite", viteFragment)

	// Initialize database
	db, err := database.New("countries.sqlite")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
		return
	}
	defer db.Conn.Close()

	// Initialize session store
	sessionStore := auth.NewSessionStore()

	// Initialize handlers
	handler := handlers.New(db, in, sessionStore)

	// Initialize server
	srv := server.New(in, handler)

	// Setup HTTP multiplexer
	mux := http.NewServeMux()

	// Setup routes
	srv.SetupRoutes(mux, *isDev)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
