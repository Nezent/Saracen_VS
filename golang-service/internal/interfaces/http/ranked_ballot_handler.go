package http

import (
	"encoding/json"
	"net/http"

	"github.com/Nezent/Saracen_Voting_System/internal/application"
	"github.com/Nezent/Saracen_Voting_System/internal/domain/ballot"
)

// RankedBallotHandler handles HTTP requests for ranked ballots (Q19)
type RankedBallotHandler struct {
	service *application.RankedBallotService
}

// NewRankedBallotHandler creates a new ranked ballot handler
func NewRankedBallotHandler(service *application.RankedBallotService) *RankedBallotHandler {
	return &RankedBallotHandler{service: service}
}

// CreateRankedBallot handles POST /api/ballots/ranked
func (h *RankedBallotHandler) CreateRankedBallot(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req ballot.RankedBallotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid JSON format",
		})
		return
	}

	// Create ranked ballot
	response, err := h.service.CreateRankedBallot(&req)
	if err != nil {
		// Determine appropriate HTTP status code
		statusCode := http.StatusInternalServerError
		if containsValidationError(err.Error()) {
			statusCode = http.StatusBadRequest
		} else if containsDuplicateError(err.Error()) {
			statusCode = http.StatusConflict
		} else if containsNotFoundError(err.Error()) {
			statusCode = http.StatusNotFound
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetRankedBallot handles GET /api/ballots/ranked/{ballot_id}
func (h *RankedBallotHandler) GetRankedBallot(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract ballot ID from URL path
	ballotID := extractPathParameter(r.URL.Path, "/api/ballots/ranked/")
	if ballotID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ballot_id is required",
		})
		return
	}

	// Get ranked ballot
	rankedBallot, rankings, err := h.service.GetRankedBallot(ballotID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if containsNotFoundError(err.Error()) {
			statusCode = http.StatusNotFound
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Combine ballot and rankings for response
	response := map[string]interface{}{
		"ballot":   rankedBallot,
		"rankings": rankings,
	}

	// Return ranked ballot with rankings
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetRankedBallotsByElection handles GET /api/ballots/ranked?election_id={id}
func (h *RankedBallotHandler) GetRankedBallotsByElection(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract election ID from query parameters
	electionID := r.URL.Query().Get("election_id")
	if electionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "election_id query parameter is required",
		})
		return
	}

	// Get ranked ballots for election
	ballots, err := h.service.GetRankedBallotsByElection(electionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return ballots
	response := map[string]interface{}{
		"election_id": electionID,
		"ballots":     ballots,
		"count":       len(ballots),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetSchulzeResults handles GET /api/ballots/ranked/results?election_id={id}
func (h *RankedBallotHandler) GetSchulzeResults(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract election ID from query parameters
	electionID := r.URL.Query().Get("election_id")
	if electionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "election_id query parameter is required",
		})
		return
	}

	// Calculate Schulze results
	results, err := h.service.CalculateSchulzeWinner(electionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return Schulze results
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// GetVoterBallots handles GET /api/ballots/ranked/voter/{voter_id}
func (h *RankedBallotHandler) GetVoterBallots(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract voter ID from URL path
	voterIDStr := extractPathParameter(r.URL.Path, "/api/ballots/ranked/voter/")
	if voterIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "voter_id is required",
		})
		return
	}

	// Convert voter ID to int
	voterID := parseIntFromString(voterIDStr)
	if voterID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "voter_id must be a positive integer",
		})
		return
	}

	// Get voter ballots
	ballots, err := h.service.GetVoterBallots(voterID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if containsNotFoundError(err.Error()) {
			statusCode = http.StatusNotFound
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return voter ballots
	response := map[string]interface{}{
		"voter_id": voterID,
		"ballots":  ballots,
		"count":    len(ballots),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// parseIntFromString converts string to int, returns 0 if invalid
func parseIntFromString(s string) int {
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}
