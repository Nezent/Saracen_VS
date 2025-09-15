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

	// Initialize repositories
	voterRepo := database.NewPostgresVoterRepository(db)
	voteRepo := database.NewPostgresVoteRepository(db)
	encryptedBallotRepo := database.NewEncryptedBallotRepository(db)
	rankedBallotRepo := database.NewRankedBallotRepository(db)

	// Initialize services
	voterService := application.NewVoterService(voterRepo)
	voteService := application.NewVoteService(voteRepo, voterRepo)
	encryptedBallotService := application.NewEncryptedBallotService(encryptedBallotRepo, voterRepo)
	rankedBallotService := application.NewRankedBallotService(rankedBallotRepo, voterRepo)

	// Initialize handlers
	voterHandler := httpHandler.NewVoterHandler(voterService)
	voteHandler := httpHandler.NewVoteHandler(voteService)
	encryptedBallotHandler := httpHandler.NewEncryptedBallotHandler(encryptedBallotService)
	rankedBallotHandler := httpHandler.NewRankedBallotHandler(rankedBallotService)

	// Setup routes
	router := mux.NewRouter()

	// Voter routes (Q1-Q5)
	router.HandleFunc("/api/voters", voterHandler.CreateVoter).Methods("POST")
	router.HandleFunc("/api/voters/{voter_id:[0-9]+}", voterHandler.GetVoter).Methods("GET")
	router.HandleFunc("/api/voters", voterHandler.GetAllVoters).Methods("GET")
	router.HandleFunc("/api/voters/{voter_id:[0-9]+}", voterHandler.UpdateVoter).Methods("PUT")
	router.HandleFunc("/api/voters/{voter_id:[0-9]+}", voterHandler.DeleteVoter).Methods("DELETE")

	// Vote routes (Q13, Q14, Q15)
	router.HandleFunc("/api/votes/timeline", voteHandler.GetVoteTimeline).Methods("GET")
	router.HandleFunc("/api/votes/weighted", voteHandler.CastWeightedVote).Methods("POST")
	router.HandleFunc("/api/votes/range", voteHandler.GetRangeVotes).Methods("GET")

	// Encrypted Ballot routes (Q16)
	router.HandleFunc("/api/ballots/encrypted", encryptedBallotHandler.CreateEncryptedBallot).Methods("POST")
	router.HandleFunc("/api/ballots/encrypted/{ballot_id}", encryptedBallotHandler.GetEncryptedBallot).Methods("GET")
	router.HandleFunc("/api/ballots/encrypted", encryptedBallotHandler.GetEncryptedBallotsByElection).Methods("GET")

	// Ranked Ballot routes (Q19)
	router.HandleFunc("/api/ballots/ranked", rankedBallotHandler.CreateRankedBallot).Methods("POST")
	router.HandleFunc("/api/ballots/ranked/{ballot_id}", rankedBallotHandler.GetRankedBallot).Methods("GET")
	router.HandleFunc("/api/ballots/ranked", rankedBallotHandler.GetRankedBallotsByElection).Methods("GET")
	router.HandleFunc("/api/ballots/ranked/results", rankedBallotHandler.GetSchulzeResults).Methods("GET")
	router.HandleFunc("/api/ballots/ranked/voter/{voter_id:[0-9]+}", rankedBallotHandler.GetVoterBallots).Methods("GET")

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
