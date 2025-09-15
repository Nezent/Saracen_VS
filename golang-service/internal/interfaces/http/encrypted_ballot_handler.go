package http

import (
	"encoding/json"
	"net/http"

	"github.com/Nezent/Saracen_Voting_System/internal/application"
	"github.com/Nezent/Saracen_Voting_System/internal/domain/ballot"
)

// EncryptedBallotHandler handles HTTP requests for encrypted ballots (Q16)
type EncryptedBallotHandler struct {
	service *application.EncryptedBallotService
}

// NewEncryptedBallotHandler creates a new encrypted ballot handler
func NewEncryptedBallotHandler(service *application.EncryptedBallotService) *EncryptedBallotHandler {
	return &EncryptedBallotHandler{service: service}
}

// CreateEncryptedBallot handles POST /api/ballots/encrypted
func (h *EncryptedBallotHandler) CreateEncryptedBallot(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req ballot.EncryptedBallotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid JSON format",
		})
		return
	}

	// Create encrypted ballot
	response, err := h.service.CreateEncryptedBallot(&req)
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

// GetEncryptedBallot handles GET /api/ballots/encrypted/{ballot_id}
func (h *EncryptedBallotHandler) GetEncryptedBallot(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract ballot ID from URL path
	ballotID := extractPathParameter(r.URL.Path, "/api/ballots/encrypted/")
	if ballotID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ballot_id is required",
		})
		return
	}

	// Get encrypted ballot
	encryptedBallot, err := h.service.GetEncryptedBallot(ballotID)
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

	// Return encrypted ballot
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(encryptedBallot)
}

// GetEncryptedBallotsByElection handles GET /api/ballots/encrypted?election_id={id}
func (h *EncryptedBallotHandler) GetEncryptedBallotsByElection(w http.ResponseWriter, r *http.Request) {
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

	// Get encrypted ballots for election
	ballots, err := h.service.GetEncryptedBallotsByElection(electionID)
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

// Helper functions for error classification
func containsValidationError(errMsg string) bool {
	validationKeywords := []string{
		"validation failed",
		"required",
		"invalid",
		"must be",
		"base64",
		"hexadecimal",
		"format",
	}

	for _, keyword := range validationKeywords {
		if contains(errMsg, keyword) {
			return true
		}
	}
	return false
}

func containsDuplicateError(errMsg string) bool {
	return contains(errMsg, "duplicate") || contains(errMsg, "already") || contains(errMsg, "nullifier")
}

func containsNotFoundError(errMsg string) bool {
	return contains(errMsg, "not found") || contains(errMsg, "does not exist")
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		(len(str) > len(substr) &&
			(containsAt(str, substr, 0) ||
				containsAt(str, substr, len(str)-len(substr)) ||
				containsInMiddle(str, substr))))
}

func containsAt(str, substr string, pos int) bool {
	if pos < 0 || pos+len(substr) > len(str) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if str[pos+i] != substr[i] {
			return false
		}
	}
	return true
}

func containsInMiddle(str, substr string) bool {
	for i := 1; i < len(str)-len(substr); i++ {
		if containsAt(str, substr, i) {
			return true
		}
	}
	return false
}

func extractPathParameter(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	if path[:len(prefix)] != prefix {
		return ""
	}
	return path[len(prefix):]
}
