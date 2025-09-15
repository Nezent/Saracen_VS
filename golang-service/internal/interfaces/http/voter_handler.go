package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/voter"
	"github.com/gorilla/mux"
)

// VoterHandler handles HTTP requests for voter operations
type VoterHandler struct {
	service voter.Service
}

// NewVoterHandler creates a new voter HTTP handler
func NewVoterHandler(service voter.Service) *VoterHandler {
	return &VoterHandler{service: service}
}

// CreateVoter handles POST /api/voters
func (h *VoterHandler) CreateVoter(w http.ResponseWriter, r *http.Request) {
	var req voter.VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.CreateVoter(req)
	if err != nil {
		if err.Error() == fmt.Sprintf("voter with id: %d already exists", req.VoterID) {
			h.writeErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		if err.Error() == fmt.Sprintf("invalid age: %d, must be 18 or older", req.Age) {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, response)
}

// GetVoter handles GET /api/voters/{voter_id}
func (h *VoterHandler) GetVoter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voterID, err := strconv.Atoi(vars["voter_id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid voter ID")
		return
	}

	response, err := h.service.GetVoter(voterID)
	if err != nil {
		if err.Error() == fmt.Sprintf("voter with id: %d was not found", voterID) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetAllVoters handles GET /api/voters
func (h *VoterHandler) GetAllVoters(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.GetAllVoters()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// UpdateVoter handles PUT /api/voters/{voter_id}
func (h *VoterHandler) UpdateVoter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voterID, err := strconv.Atoi(vars["voter_id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid voter ID")
		return
	}

	var req voter.VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.UpdateVoter(voterID, req)
	if err != nil {
		if err.Error() == fmt.Sprintf("voter with id: %d was not found", voterID) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == fmt.Sprintf("invalid age: %d, must be 18 or older", req.Age) {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// DeleteVoter handles DELETE /api/voters/{voter_id}
func (h *VoterHandler) DeleteVoter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voterID, err := strconv.Atoi(vars["voter_id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid voter ID")
		return
	}

	err = h.service.DeleteVoter(voterID)
	if err != nil {
		if err.Error() == fmt.Sprintf("voter with id: %d was not found", voterID) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	successResponse := map[string]string{
		"message": fmt.Sprintf("voter with id: %d deleted successfully", voterID),
	}
	h.writeJSONResponse(w, http.StatusOK, successResponse)
}

// writeJSONResponse writes a JSON response
func (h *VoterHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes an error response
func (h *VoterHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := voter.ErrorResponse{Message: message}
	h.writeJSONResponse(w, statusCode, errorResponse)
}
