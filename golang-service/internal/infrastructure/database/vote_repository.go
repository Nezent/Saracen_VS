package database

import (
	"database/sql"
	"fmt"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/vote"
	_ "github.com/lib/pq"
)

// PostgresVoteRepository implements the vote.Repository interface
type PostgresVoteRepository struct {
	db *sql.DB
}

// NewPostgresVoteRepository creates a new PostgreSQL vote repository
func NewPostgresVoteRepository(db *sql.DB) vote.Repository {
	return &PostgresVoteRepository{db: db}
}

// GetTimelineByCandidateID retrieves all votes for a specific candidate ordered by timestamp
func (r *PostgresVoteRepository) GetTimelineByCandidateID(candidateID int) ([]*vote.Vote, error) {
	query := `
		SELECT vote_id, voter_id, candidate_id, weight, created_at, updated_at
		FROM votes
		WHERE candidate_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, candidateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote timeline: %w", err)
	}
	defer rows.Close()

	var votes []*vote.Vote
	for rows.Next() {
		v := &vote.Vote{}
		err := rows.Scan(&v.VoteID, &v.VoterID, &v.CandidateID, &v.Weight, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating votes: %w", err)
	}

	return votes, nil
}

// CreateWeightedVote inserts a new weighted vote into the database
func (r *PostgresVoteRepository) CreateWeightedVote(v *vote.Vote) error {
	query := `
		INSERT INTO votes (voter_id, candidate_id, weight, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING vote_id
	`

	err := r.db.QueryRow(query, v.VoterID, v.CandidateID, v.Weight, v.CreatedAt, v.UpdatedAt).Scan(&v.VoteID)
	if err != nil {
		return fmt.Errorf("failed to create weighted vote: %w", err)
	}

	return nil
}

// GetByID retrieves a vote by its ID
func (r *PostgresVoteRepository) GetByID(voteID int) (*vote.Vote, error) {
	query := `
		SELECT vote_id, voter_id, candidate_id, weight, created_at, updated_at
		FROM votes
		WHERE vote_id = $1
	`

	var v vote.Vote
	err := r.db.QueryRow(query, voteID).Scan(
		&v.VoteID,
		&v.VoterID,
		&v.CandidateID,
		&v.Weight,
		&v.CreatedAt,
		&v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vote with id: %d was not found", voteID)
		}
		return nil, fmt.Errorf("failed to get vote by ID: %w", err)
	}

	return &v, nil
}

// HasVoted checks if a voter has already cast a vote
func (r *PostgresVoteRepository) HasVoted(voterID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM votes WHERE voter_id = $1)`

	var exists bool
	err := r.db.QueryRow(query, voterID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if voter has voted: %w", err)
	}

	return exists, nil
}

// GetVotesInRange counts votes for a candidate within a specific time range
func (r *PostgresVoteRepository) GetVotesInRange(candidateID int, from, to string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM votes
		WHERE candidate_id = $1 
		AND created_at >= $2 
		AND created_at <= $3
	`

	var count int
	err := r.db.QueryRow(query, candidateID, from, to).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get votes in range: %w", err)
	}

	return count, nil
}
