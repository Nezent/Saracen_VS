package ballot

import (
	"fmt"
	"sort"
	"time"
)

// RankedBallot represents a ranked ballot for Q19
type RankedBallot struct {
	BallotID   string    `json:"ballot_id" db:"ballot_id"`
	ElectionID string    `json:"election_id" db:"election_id"`
	VoterID    int       `json:"voter_id" db:"voter_id"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Status     string    `json:"status" db:"status"`
}

// BallotRanking represents individual candidate rankings within a ballot
type BallotRanking struct {
	ID           int    `json:"id" db:"id"`
	BallotID     string `json:"ballot_id" db:"ballot_id"`
	CandidateID  int    `json:"candidate_id" db:"candidate_id"`
	RankPosition int    `json:"rank_position" db:"rank_position"`
}

// RankedBallotRequest represents the request payload for Q19
type RankedBallotRequest struct {
	ElectionID string    `json:"election_id" validate:"required"`
	VoterID    int       `json:"voter_id" validate:"required,min=1"`
	Ranking    []int     `json:"ranking" validate:"required,min=1"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
}

// RankedBallotResponse represents the response for Q19
type RankedBallotResponse struct {
	BallotID string `json:"ballot_id"`
	Status   string `json:"status"`
}

// SchulzeResult represents the result of Schulze method calculation
type SchulzeResult struct {
	ElectionID string                 `json:"election_id"`
	Winners    []int                  `json:"winners"`
	Rankings   []SchulzeCandidateRank `json:"rankings"`
	Matrix     [][]int                `json:"pairwise_matrix,omitempty"`
}

// SchulzeCandidateRank represents a candidate's rank in Schulze results
type SchulzeCandidateRank struct {
	CandidateID int `json:"candidate_id"`
	Rank        int `json:"rank"`
	Score       int `json:"score"`
}

// Validate validates the ranked ballot request
func (req *RankedBallotRequest) Validate() error {
	if req.ElectionID == "" {
		return fmt.Errorf("election_id is required")
	}

	if req.VoterID <= 0 {
		return fmt.Errorf("voter_id must be positive")
	}

	if len(req.Ranking) == 0 {
		return fmt.Errorf("ranking array cannot be empty")
	}

	// Validate ranking contains unique candidate IDs
	candidateSet := make(map[int]bool)
	for i, candidateID := range req.Ranking {
		if candidateID <= 0 {
			return fmt.Errorf("candidate_id at position %d must be positive", i)
		}
		if candidateSet[candidateID] {
			return fmt.Errorf("candidate_id %d appears multiple times in ranking", candidateID)
		}
		candidateSet[candidateID] = true
	}

	if req.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	return nil
}

// ToRankedBallot converts request to domain model with generated ID
func (req *RankedBallotRequest) ToRankedBallot() (*RankedBallot, []BallotRanking, error) {
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	ballotID := generateBallotID("rb_")

	ballot := &RankedBallot{
		BallotID:   ballotID,
		ElectionID: req.ElectionID,
		VoterID:    req.VoterID,
		Timestamp:  req.Timestamp,
		Status:     "accepted",
	}

	// Create individual rankings
	rankings := make([]BallotRanking, len(req.Ranking))
	for i, candidateID := range req.Ranking {
		rankings[i] = BallotRanking{
			BallotID:     ballotID,
			CandidateID:  candidateID,
			RankPosition: i + 1, // Rankings start from 1
		}
	}

	return ballot, rankings, nil
}

// ToResponse converts domain model to response
func (rb *RankedBallot) ToResponse() *RankedBallotResponse {
	return &RankedBallotResponse{
		BallotID: rb.BallotID,
		Status:   rb.Status,
	}
}

// CalculateSchulze implements the Schulze method for ranked choice voting
func CalculateSchulze(ballots []RankedBallotWithRankings) *SchulzeResult {
	if len(ballots) == 0 {
		return &SchulzeResult{Winners: []int{}, Rankings: []SchulzeCandidateRank{}}
	}

	// Get all unique candidates
	candidateSet := make(map[int]bool)
	for _, ballot := range ballots {
		for _, ranking := range ballot.Rankings {
			candidateSet[ranking.CandidateID] = true
		}
	}

	candidates := make([]int, 0, len(candidateSet))
	for candidate := range candidateSet {
		candidates = append(candidates, candidate)
	}
	sort.Ints(candidates)

	n := len(candidates)
	if n == 0 {
		return &SchulzeResult{Winners: []int{}, Rankings: []SchulzeCandidateRank{}}
	}

	// Create candidate index mapping
	candidateIndex := make(map[int]int)
	for i, candidate := range candidates {
		candidateIndex[candidate] = i
	}

	// Initialize pairwise comparison matrix
	d := make([][]int, n)
	for i := range d {
		d[i] = make([]int, n)
	}

	// Count pairwise preferences
	for _, ballot := range ballots {
		rankMap := make(map[int]int)
		for _, ranking := range ballot.Rankings {
			rankMap[ranking.CandidateID] = ranking.RankPosition
		}

		// Compare all pairs of candidates
		for i, candidateA := range candidates {
			for j, candidateB := range candidates {
				if i != j {
					rankA, hasA := rankMap[candidateA]
					rankB, hasB := rankMap[candidateB]

					// If both candidates are ranked and A is ranked higher (lower number)
					if hasA && hasB && rankA < rankB {
						d[i][j]++
					}
				}
			}
		}
	}

	// Floyd-Warshall algorithm to find strongest paths
	p := make([][]int, n)
	for i := range p {
		p[i] = make([]int, n)
		for j := range p[i] {
			if i != j {
				if d[i][j] > d[j][i] {
					p[i][j] = d[i][j]
				} else {
					p[i][j] = 0
				}
			}
		}
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && i != k && j != k {
					p[i][j] = max(p[i][j], min(p[i][k], p[k][j]))
				}
			}
		}
	}

	// Find winners (Condorcet winners)
	winners := []int{}
	rankings := []SchulzeCandidateRank{}

	for i, candidate := range candidates {
		isWinner := true
		score := 0

		for j := 0; j < n; j++ {
			if i != j {
				if p[j][i] > p[i][j] {
					isWinner = false
				}
				score += p[i][j]
			}
		}

		if isWinner {
			winners = append(winners, candidate)
		}

		rankings = append(rankings, SchulzeCandidateRank{
			CandidateID: candidate,
			Score:       score,
		})
	}

	// Sort rankings by score (descending)
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	// Assign ranks
	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return &SchulzeResult{
		Winners:  winners,
		Rankings: rankings,
		Matrix:   d,
	}
}

// RankedBallotWithRankings combines ballot with its rankings for calculation
type RankedBallotWithRankings struct {
	Ballot   RankedBallot    `json:"ballot"`
	Rankings []BallotRanking `json:"rankings"`
}

// RankedBallotRepository defines repository interface for ranked ballots
type RankedBallotRepository interface {
	Create(ballot *RankedBallot, rankings []BallotRanking) error
	GetByBallotID(ballotID string) (*RankedBallot, []BallotRanking, error)
	GetByElectionID(electionID string) ([]RankedBallotWithRankings, error)
	GetByVoterID(voterID int) ([]*RankedBallot, error)
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
