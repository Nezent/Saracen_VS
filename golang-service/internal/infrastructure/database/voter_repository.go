package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/voter"
	_ "github.com/lib/pq"
)

// PostgresVoterRepository implements the voter.Repository interface
type PostgresVoterRepository struct {
	db *sql.DB
}

// NewPostgresVoterRepository creates a new PostgreSQL voter repository
func NewPostgresVoterRepository(db *sql.DB) voter.Repository {
	return &PostgresVoterRepository{db: db}
}

// Create inserts a new voter into the database
func (r *PostgresVoterRepository) Create(v *voter.Voter) error {
	query := `
		INSERT INTO voter (voter_id, name, age, has_voted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	v.HasVoted = false

	_, err := r.db.Exec(query, v.VoterID, v.Name, v.Age, v.HasVoted, v.CreatedAt, v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create voter: %w", err)
	}

	return nil
}

// GetByID retrieves a voter by their ID
func (r *PostgresVoterRepository) GetByID(voterID int) (*voter.Voter, error) {
	query := `
		SELECT voter_id, name, age, has_voted, created_at, updated_at
		FROM voter
		WHERE voter_id = $1
	`

	v := &voter.Voter{}
	err := r.db.QueryRow(query, voterID).Scan(
		&v.VoterID, &v.Name, &v.Age, &v.HasVoted, &v.CreatedAt, &v.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("voter with id: %d was not found", voterID)
		}
		return nil, fmt.Errorf("failed to get voter: %w", err)
	}

	return v, nil
}

// GetAll retrieves all voters from the database
func (r *PostgresVoterRepository) GetAll() ([]*voter.Voter, error) {
	query := `
		SELECT voter_id, name, age, has_voted, created_at, updated_at
		FROM voter
		ORDER BY voter_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all voters: %w", err)
	}
	defer rows.Close()

	var voters []*voter.Voter
	for rows.Next() {
		v := &voter.Voter{}
		err := rows.Scan(&v.VoterID, &v.Name, &v.Age, &v.HasVoted, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan voter: %w", err)
		}
		voters = append(voters, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating voters: %w", err)
	}

	return voters, nil
}

// Update updates an existing voter
func (r *PostgresVoterRepository) Update(v *voter.Voter) error {
	query := `
		UPDATE voter
		SET name = $2, age = $3, updated_at = $4
		WHERE voter_id = $1
	`

	v.UpdatedAt = time.Now()
	result, err := r.db.Exec(query, v.VoterID, v.Name, v.Age, v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update voter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("voter with id: %d was not found", v.VoterID)
	}

	return nil
}

// Delete removes a voter from the database
func (r *PostgresVoterRepository) Delete(voterID int) error {
	query := `DELETE FROM voter WHERE voter_id = $1`

	result, err := r.db.Exec(query, voterID)
	if err != nil {
		return fmt.Errorf("failed to delete voter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("voter with id: %d was not found", voterID)
	}

	return nil
}

// ExistsByID checks if a voter exists with the given ID
func (r *PostgresVoterRepository) ExistsByID(voterID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM voter WHERE voter_id = $1)`

	var exists bool
	err := r.db.QueryRow(query, voterID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check voter existence: %w", err)
	}

	return exists, nil
}
