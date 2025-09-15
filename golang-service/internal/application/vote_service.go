package application

import (
	"fmt"
	"time"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/vote"
	"github.com/Nezent/Saracen_Voting_System/internal/domain/voter"
)

// VoteService implements the vote.Service interface
type VoteService struct {
	voteRepo  vote.Repository
	voterRepo voter.Repository
}

// NewVoteService creates a new vote service
func NewVoteService(voteRepo vote.Repository, voterRepo voter.Repository) vote.Service {
	return &VoteService{
		voteRepo:  voteRepo,
		voterRepo: voterRepo,
	}
}

// GetVoteTimeline retrieves the timeline of votes for a specific candidate
func (s *VoteService) GetVoteTimeline(candidateID int) (*vote.VoteTimelineResponse, error) {
	votes, err := s.voteRepo.GetTimelineByCandidateID(candidateID)
	if err != nil {
		return nil, err
	}

	var timeline []vote.VoteTimelineItem
	for _, v := range votes {
		timeline = append(timeline, vote.VoteTimelineItem{
			VoteID:    v.VoteID,
			Timestamp: v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339 format
		})
	}

	return &vote.VoteTimelineResponse{
		CandidateID: candidateID,
		Timeline:    timeline,
	}, nil
}

// CastWeightedVote casts a weighted vote based on voter profile update status
func (s *VoteService) CastWeightedVote(req vote.WeightedVoteRequest) (*vote.WeightedVoteResponse, error) {
	// Check if voter has already voted
	hasVoted, err := s.voteRepo.HasVoted(req.VoterID)
	if err != nil {
		return nil, fmt.Errorf("error checking if voter has voted: %w", err)
	}
	if hasVoted {
		return nil, fmt.Errorf("voter with id: %d has already voted", req.VoterID)
	}

	// Get voter information to determine weight based on profile update status
	voterInfo, err := s.voterRepo.GetByID(req.VoterID)
	if err != nil {
		return nil, fmt.Errorf("voter with id: %d was not found", req.VoterID)
	}

	// Determine weight based on profile update status
	// If voter has recent activity (updated_at different from created_at), give weight 2
	// Otherwise, give weight 1
	weight := 1
	if voterInfo.UpdatedAt.After(voterInfo.CreatedAt.Add(time.Minute)) {
		weight = 2 // Higher weight for voters who updated their profile
	}

	// Create the vote
	now := time.Now()
	v := &vote.Vote{
		VoterID:     req.VoterID,
		CandidateID: req.CandidateID,
		Weight:      weight,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save the vote
	err = s.voteRepo.CreateWeightedVote(v)
	if err != nil {
		return nil, fmt.Errorf("failed to cast weighted vote: %w", err)
	}

	// Get the actual vote from database to ensure we return the stored weight
	storedVote, err := s.voteRepo.GetByID(v.VoteID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created vote: %w", err)
	}

	// Update voter's has_voted status
	voterInfo.HasVoted = true
	err = s.voterRepo.Update(voterInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to update voter status: %w", err)
	}

	return &vote.WeightedVoteResponse{
		VoteID:      storedVote.VoteID,
		VoterID:     storedVote.VoterID,
		CandidateID: storedVote.CandidateID,
		Weight:      storedVote.Weight,
	}, nil
}

// GetRangeVotes gets votes for a candidate within a specific time range
func (s *VoteService) GetRangeVotes(candidateID int, from, to string) (*vote.RangeVoteResponse, error) {
	// Parse and validate time strings
	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid from time format: %w", err)
	}

	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, fmt.Errorf("invalid to time format: %w", err)
	}

	// Validate that from is before to
	if fromTime.After(toTime) {
		return nil, fmt.Errorf("invalid interval: from > to")
	}

	// Get votes count in range
	votesGained, err := s.voteRepo.GetVotesInRange(candidateID, from, to)
	if err != nil {
		return nil, err
	}

	return &vote.RangeVoteResponse{
		CandidateID: candidateID,
		From:        from,
		To:          to,
		VotesGained: votesGained,
	}, nil
}
