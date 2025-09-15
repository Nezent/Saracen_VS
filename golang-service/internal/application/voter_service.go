package application

import (
	"fmt"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/voter"
)

// VoterService implements the voter.Service interface
type VoterService struct {
	repo voter.Repository
}

// NewVoterService creates a new voter service
func NewVoterService(repo voter.Repository) voter.Service {
	return &VoterService{repo: repo}
}

// CreateVoter creates a new voter with validation
func (s *VoterService) CreateVoter(req voter.VoterRequest) (*voter.VoterResponse, error) {
	// Check if voter already exists
	if req.VoterID != 0 {
		exists, err := s.repo.ExistsByID(req.VoterID)
		if err != nil {
			return nil, fmt.Errorf("error checking voter existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("voter with id: %d already exists", req.VoterID)
		}
	}

	// Create voter model
	v := &voter.Voter{
		VoterID: req.VoterID,
		Name:    req.Name,
		Age:     req.Age,
	}

	// Validate age
	if err := v.ValidateAge(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Create(v); err != nil {
		return nil, fmt.Errorf("failed to create voter: %w", err)
	}

	// Return response
	return &voter.VoterResponse{
		VoterID:  v.VoterID,
		Name:     v.Name,
		Age:      v.Age,
		HasVoted: v.HasVoted,
	}, nil
}

// GetVoter retrieves a voter by ID
func (s *VoterService) GetVoter(voterID int) (*voter.VoterResponse, error) {
	v, err := s.repo.GetByID(voterID)
	if err != nil {
		return nil, err
	}

	return &voter.VoterResponse{
		VoterID:  v.VoterID,
		Name:     v.Name,
		Age:      v.Age,
		HasVoted: v.HasVoted,
	}, nil
}

// GetAllVoters retrieves all voters
func (s *VoterService) GetAllVoters() (*voter.VotersListResponse, error) {
	voters, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var voterItems []voter.VoterListItem
	for _, v := range voters {
		voterItems = append(voterItems, voter.VoterListItem{
			VoterID: v.VoterID,
			Name:    v.Name,
			Age:     v.Age,
		})
	}

	return &voter.VotersListResponse{
		Voters: voterItems,
	}, nil
}

// UpdateVoter updates an existing voter
func (s *VoterService) UpdateVoter(voterID int, req voter.VoterRequest) (*voter.VoterResponse, error) {
	// Get existing voter
	existingVoter, err := s.repo.GetByID(voterID)
	if err != nil {
		return nil, err
	}

	// Create updated voter model
	updatedVoter := &voter.Voter{
		VoterID:   voterID,
		Name:      req.Name,
		Age:       req.Age,
		HasVoted:  existingVoter.HasVoted,  // Preserve has_voted status
		CreatedAt: existingVoter.CreatedAt, // Preserve created_at
	}

	// Validate age
	if err := updatedVoter.ValidateAge(); err != nil {
		return nil, err
	}

	// Update in repository
	if err := s.repo.Update(updatedVoter); err != nil {
		return nil, fmt.Errorf("failed to update voter: %w", err)
	}

	// Return response
	return &voter.VoterResponse{
		VoterID:  updatedVoter.VoterID,
		Name:     updatedVoter.Name,
		Age:      updatedVoter.Age,
		HasVoted: updatedVoter.HasVoted,
	}, nil
}

// DeleteVoter deletes a voter by ID
func (s *VoterService) DeleteVoter(voterID int) error {
	return s.repo.Delete(voterID)
}
