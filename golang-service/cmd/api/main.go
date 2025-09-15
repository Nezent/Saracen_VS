package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Nezent/Saracen_Voting_System/internal/application"
	"github.com/Nezent/Saracen_Voting_System/internal/infrastructure/database"
	httpHandler "github.com/Nezent/Saracen_Voting_System/internal/interfaces/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get database URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize repository
	voterRepo := database.NewPostgresVoterRepository(db)

	// Initialize service
	voterService := application.NewVoterService(voterRepo)

	// Initialize handler
	voterHandler := httpHandler.NewVoterHandler(voterService)

	// Setup routes
	router := mux.NewRouter()

	// Voter routes (Q1-Q5)
	router.HandleFunc("/api/voters", voterHandler.CreateVoter).Methods("POST")
	router.HandleFunc("/api/voters/{voter_id:[0-9]+}", voterHandler.GetVoter).Methods("GET")
	router.HandleFunc("/api/voters", voterHandler.GetAllVoters).Methods("GET")
	router.HandleFunc("/api/voters/{voter_id:[0-9]+}", voterHandler.UpdateVoter).Methods("PUT")
	router.HandleFunc("/api/voters/{voter_id:[0-9]+}", voterHandler.DeleteVoter).Methods("DELETE")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
