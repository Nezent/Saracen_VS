package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/vote"
	"github.com/Nezent/Saracen_Voting_System/internal/domain/voter"
)

// VoteHandler handles HTTP requests for vote operations
type VoteHandler struct {
	service vote.Service
}

// NewVoteHandler creates a new vote HTTP handler
func NewVoteHandler(service vote.Service) *VoteHandler {
	return &VoteHandler{service: service}
}

// GetVoteTimeline handles GET /api/votes/timeline?candidate_id={id}
func (h *VoteHandler) GetVoteTimeline(w http.ResponseWriter, r *http.Request) {
	// Get candidate_id from query parameter
	candidateIDStr := r.URL.Query().Get("candidate_id")
	if candidateIDStr == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "candidate_id query parameter is required")
		return
	}

	candidateID, err := strconv.Atoi(candidateIDStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid candidate_id")
		return
	}

	response, err := h.service.GetVoteTimeline(candidateID)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// CastWeightedVote handles POST /api/votes/weighted
func (h *VoteHandler) CastWeightedVote(w http.ResponseWriter, r *http.Request) {
	var req vote.WeightedVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.CastWeightedVote(req)
	if err != nil {
		if err.Error() == fmt.Sprintf("voter with id: %d has already voted", req.VoterID) {
			h.writeErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		if err.Error() == fmt.Sprintf("voter with id: %d was not found", req.VoterID) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, response)
}

// GetRangeVotes handles GET /api/votes/range?candidate_id={id}&from={t1}&to={t2}
func (h *VoteHandler) GetRangeVotes(w http.ResponseWriter, r *http.Request) {
	// Get candidate_id from query parameter
	candidateIDStr := r.URL.Query().Get("candidate_id")
	if candidateIDStr == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "candidate_id query parameter is required")
		return
	}

	candidateID, err := strconv.Atoi(candidateIDStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid candidate_id")
		return
	}

	// Get from and to query parameters
	from := r.URL.Query().Get("from")
	if from == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "from query parameter is required")
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "to query parameter is required")
		return
	}

	response, err := h.service.GetRangeVotes(candidateID, from, to)
	if err != nil {
		if err.Error() == "invalid interval: from > to" {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if fmt.Sprintf("%v", err)[0:7] == "invalid" { // Check for time format errors
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// writeJSONResponse writes a JSON response
func (h *VoteHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes an error response
func (h *VoteHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := voter.ErrorResponse{Message: message}
	h.writeJSONResponse(w, statusCode, errorResponse)
}
