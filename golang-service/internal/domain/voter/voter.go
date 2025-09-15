package voter

import (
	"errors"
	"strconv"
	"time"
)

// Voter represents the domain model for a voter
type Voter struct {
	VoterID   int       `json:"voter_id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	HasVoted  bool      `json:"has_voted"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// VoterRequest represents the request payload for creating/updating a voter
type VoterRequest struct {
	VoterID int    `json:"voter_id,omitempty"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

// VoterResponse represents the response payload for voter operations
type VoterResponse struct {
	VoterID  int    `json:"voter_id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	HasVoted bool   `json:"has_voted"`
}

// VoterListItem represents a voter item in the voters list (without has_voted)
type VoterListItem struct {
	VoterID int    `json:"voter_id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

// VotersListResponse represents the response for listing all voters
type VotersListResponse struct {
	Voters []VoterListItem `json:"voters"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// ValidateAge validates that the voter is at least 18 years old
func (v *Voter) ValidateAge() error {
	if v.Age < 18 {
		return errors.New("invalid age: " + strconv.Itoa(v.Age) + ", must be 18 or older")
	}
	return nil
}

// Repository defines the interface for voter data operations
type Repository interface {
	Create(voter *Voter) error
	GetByID(voterID int) (*Voter, error)
	GetAll() ([]*Voter, error)
	Update(voter *Voter) error
	Delete(voterID int) error
	ExistsByID(voterID int) (bool, error)
}

// Service defines the interface for voter business logic
type Service interface {
	CreateVoter(req VoterRequest) (*VoterResponse, error)
	GetVoter(voterID int) (*VoterResponse, error)
	GetAllVoters() (*VotersListResponse, error)
	UpdateVoter(voterID int, req VoterRequest) (*VoterResponse, error)
	DeleteVoter(voterID int) error
}
