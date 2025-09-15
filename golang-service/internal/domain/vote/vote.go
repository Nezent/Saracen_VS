package vote

import (
	"time"
)

// Vote represents the domain model for a vote
type Vote struct {
	VoteID      int       `json:"vote_id"`
	VoterID     int       `json:"voter_id"`
	CandidateID int       `json:"candidate_id"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// VoteTimelineItem represents a single vote in the timeline
type VoteTimelineItem struct {
	VoteID    int    `json:"vote_id"`
	Timestamp string `json:"timestamp"`
}

// VoteTimelineResponse represents the response for vote timeline
type VoteTimelineResponse struct {
	CandidateID int                `json:"candidate_id"`
	Timeline    []VoteTimelineItem `json:"timeline"`
}

// WeightedVoteRequest represents the request for casting a weighted vote
type WeightedVoteRequest struct {
	VoterID     int `json:"voter_id"`
	CandidateID int `json:"candidate_id"`
}

// WeightedVoteResponse represents the response for casting a weighted vote
type WeightedVoteResponse struct {
	VoteID      int `json:"vote_id"`
	VoterID     int `json:"voter_id"`
	CandidateID int `json:"candidate_id"`
	Weight      int `json:"weight"`
}

// RangeVoteResponse represents the response for range vote queries
type RangeVoteResponse struct {
	CandidateID int    `json:"candidate_id"`
	From        string `json:"from"`
	To          string `json:"to"`
	VotesGained int    `json:"votes_gained"`
}

// Repository defines the interface for vote data operations
type Repository interface {
	GetTimelineByCandidateID(candidateID int) ([]*Vote, error)
	CreateWeightedVote(vote *Vote) error
	HasVoted(voterID int) (bool, error)
	GetVotesInRange(candidateID int, from, to string) (int, error)
	GetByID(voteID int) (*Vote, error)
}

// Service defines the interface for vote business logic
type Service interface {
	GetVoteTimeline(candidateID int) (*VoteTimelineResponse, error)
	CastWeightedVote(req WeightedVoteRequest) (*WeightedVoteResponse, error)
	GetRangeVotes(candidateID int, from, to string) (*RangeVoteResponse, error)
}
