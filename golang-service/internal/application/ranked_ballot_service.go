package application

import (
	"fmt"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/ballot"
	"github.com/Nezent/Saracen_Voting_System/internal/domain/voter"
)

// RankedBallotService handles ranked ballot business logic and Schulze calculations
type RankedBallotService struct {
	rankedBallotRepo ballot.RankedBallotRepository
	voterRepo        voter.Repository
}

// NewRankedBallotService creates a new ranked ballot service
func NewRankedBallotService(
	rankedBallotRepo ballot.RankedBallotRepository,
	voterRepo voter.Repository,
) *RankedBallotService {
	return &RankedBallotService{
		rankedBallotRepo: rankedBallotRepo,
		voterRepo:        voterRepo,
	}
}

// CreateRankedBallot creates a new ranked ballot
func (s *RankedBallotService) CreateRankedBallot(req *ballot.RankedBallotRequest) (*ballot.RankedBallotResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// Verify voter exists
	voterEntity, err := s.voterRepo.GetByID(req.VoterID)
	if err != nil {
		return nil, fmt.Errorf("voter not found: %v", err)
	}

	// Check if voter has already voted in this election
	existingBallots, err := s.rankedBallotRepo.GetByVoterID(req.VoterID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing ballots: %v", err)
	}

	for _, existingBallot := range existingBallots {
		if existingBallot.ElectionID == req.ElectionID {
			return nil, fmt.Errorf("voter %d has already voted in election %s", req.VoterID, req.ElectionID)
		}
	}

	// Validate election ID
	if err := s.ValidateElectionID(req.ElectionID); err != nil {
		return nil, fmt.Errorf("invalid election: %v", err)
	}

	// TODO: In a real implementation, validate that all candidate IDs in ranking exist
	// For now, we'll assume they're valid

	// Convert request to domain model
	rankedBallot, rankings, err := req.ToRankedBallot()
	if err != nil {
		return nil, fmt.Errorf("failed to create ranked ballot: %v", err)
	}

	// Store the ranked ballot with its rankings
	if err := s.rankedBallotRepo.Create(rankedBallot, rankings); err != nil {
		return nil, fmt.Errorf("failed to store ranked ballot: %v", err)
	}

	// Update voter's has_voted status
	voterEntity.HasVoted = true
	if err := s.voterRepo.Update(voterEntity); err != nil {
		// Log the error but don't fail the ballot creation
		fmt.Printf("Warning: failed to update voter has_voted status: %v\n", err)
	}

	return rankedBallot.ToResponse(), nil
}

// GetRankedBallot retrieves a ranked ballot by ID with its rankings
func (s *RankedBallotService) GetRankedBallot(ballotID string) (*ballot.RankedBallot, []ballot.BallotRanking, error) {
	if ballotID == "" {
		return nil, nil, fmt.Errorf("ballot_id is required")
	}

	return s.rankedBallotRepo.GetByBallotID(ballotID)
}

// GetRankedBallotsByElection retrieves all ranked ballots for an election
func (s *RankedBallotService) GetRankedBallotsByElection(electionID string) ([]ballot.RankedBallotWithRankings, error) {
	if electionID == "" {
		return nil, fmt.Errorf("election_id is required")
	}

	return s.rankedBallotRepo.GetByElectionID(electionID)
}

// CalculateSchulzeWinner calculates the Schulze method winner for an election
func (s *RankedBallotService) CalculateSchulzeWinner(electionID string) (*ballot.SchulzeResult, error) {
	if electionID == "" {
		return nil, fmt.Errorf("election_id is required")
	}

	// Get all ranked ballots for the election
	ballots, err := s.rankedBallotRepo.GetByElectionID(electionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ranked ballots: %v", err)
	}

	if len(ballots) == 0 {
		return &ballot.SchulzeResult{
			ElectionID: electionID,
			Winners:    []int{},
			Rankings:   []ballot.SchulzeCandidateRank{},
		}, nil
	}

	// Calculate Schulze winner using the algorithm
	result := ballot.CalculateSchulze(ballots)
	result.ElectionID = electionID

	return result, nil
}

// GetVoterBallots retrieves all ballots for a specific voter
func (s *RankedBallotService) GetVoterBallots(voterID int) ([]*ballot.RankedBallot, error) {
	if voterID <= 0 {
		return nil, fmt.Errorf("voter_id must be positive")
	}

	// Verify voter exists
	_, err := s.voterRepo.GetByID(voterID)
	if err != nil {
		return nil, fmt.Errorf("voter not found: %v", err)
	}

	return s.rankedBallotRepo.GetByVoterID(voterID)
}

// ValidateElectionID validates election ID format and timing
func (s *RankedBallotService) ValidateElectionID(electionID string) error {
	if electionID == "" {
		return fmt.Errorf("election_id is required")
	}

	// Basic validation - in production, check against election table
	// and verify election is active and accepting votes
	if len(electionID) < 3 {
		return fmt.Errorf("invalid election_id format")
	}

	// TODO: In production, verify:
	// 1. Election exists
	// 2. Election is currently active (within voting period)
	// 3. Election accepts ranked ballots

	return nil
}

// GetElectionResults provides comprehensive election results including Schulze analysis
func (s *RankedBallotService) GetElectionResults(electionID string) (*ballot.SchulzeResult, error) {
	if electionID == "" {
		return nil, fmt.Errorf("election_id is required")
	}

	// Calculate Schulze results
	results, err := s.CalculateSchulzeWinner(electionID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate election results: %v", err)
	}

	return results, nil
}
