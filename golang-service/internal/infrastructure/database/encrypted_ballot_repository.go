package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/ballot"
)

// EncryptedBallotPostgresRepository implements the EncryptedBallotRepository interface
type EncryptedBallotPostgresRepository struct {
	db *sql.DB
}

// NewEncryptedBallotRepository creates a new encrypted ballot repository
func NewEncryptedBallotRepository(db *sql.DB) ballot.EncryptedBallotRepository {
	return &EncryptedBallotPostgresRepository{db: db}
}

// Create stores a new encrypted ballot
func (r *EncryptedBallotPostgresRepository) Create(encryptedBallot *ballot.EncryptedBallot) error {
	query := `
		INSERT INTO encrypted_ballots 
		(ballot_id, election_id, voter_id, ciphertext, zk_proof, voter_pubkey, nullifier, signature, status, anchored_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(
		query,
		encryptedBallot.BallotID,
		encryptedBallot.ElectionID,
		encryptedBallot.VoterID,
		encryptedBallot.Ciphertext,
		encryptedBallot.ZKProof,
		encryptedBallot.VoterPubkey,
		encryptedBallot.Nullifier,
		encryptedBallot.Signature,
		encryptedBallot.Status,
		encryptedBallot.AnchoredAt,
	)
	if err != nil {
		// Check for unique constraint violation on nullifier (double voting prevention)
		if strings.Contains(err.Error(), "nullifier") && strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("duplicate nullifier: ballot already submitted for this voter")
		}
		return fmt.Errorf("failed to create encrypted ballot: %v", err)
	}

	return nil
}

// GetByBallotID retrieves an encrypted ballot by ballot ID
func (r *EncryptedBallotPostgresRepository) GetByBallotID(ballotID string) (*ballot.EncryptedBallot, error) {
	query := `
		SELECT ballot_id, election_id, voter_id, ciphertext, zk_proof, voter_pubkey, 
			   nullifier, signature, status, anchored_at
		FROM encrypted_ballots
		WHERE ballot_id = $1`

	row := r.db.QueryRow(query, ballotID)

	var encryptedBallot ballot.EncryptedBallot
	err := row.Scan(
		&encryptedBallot.BallotID,
		&encryptedBallot.ElectionID,
		&encryptedBallot.VoterID,
		&encryptedBallot.Ciphertext,
		&encryptedBallot.ZKProof,
		&encryptedBallot.VoterPubkey,
		&encryptedBallot.Nullifier,
		&encryptedBallot.Signature,
		&encryptedBallot.Status,
		&encryptedBallot.AnchoredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("encrypted ballot not found: %s", ballotID)
		}
		return nil, fmt.Errorf("failed to get encrypted ballot: %v", err)
	}

	return &encryptedBallot, nil
}

// GetByNullifier retrieves an encrypted ballot by nullifier
func (r *EncryptedBallotPostgresRepository) GetByNullifier(nullifier string) (*ballot.EncryptedBallot, error) {
	query := `
		SELECT ballot_id, election_id, voter_id, ciphertext, zk_proof, voter_pubkey, 
			   nullifier, signature, status, anchored_at
		FROM encrypted_ballots
		WHERE nullifier = $1`

	row := r.db.QueryRow(query, nullifier)

	var encryptedBallot ballot.EncryptedBallot
	err := row.Scan(
		&encryptedBallot.BallotID,
		&encryptedBallot.ElectionID,
		&encryptedBallot.VoterID,
		&encryptedBallot.Ciphertext,
		&encryptedBallot.ZKProof,
		&encryptedBallot.VoterPubkey,
		&encryptedBallot.Nullifier,
		&encryptedBallot.Signature,
		&encryptedBallot.Status,
		&encryptedBallot.AnchoredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("encrypted ballot not found for nullifier: %s", nullifier)
		}
		return nil, fmt.Errorf("failed to get encrypted ballot by nullifier: %v", err)
	}

	return &encryptedBallot, nil
}

// GetByElectionID retrieves all encrypted ballots for an election
func (r *EncryptedBallotPostgresRepository) GetByElectionID(electionID string) ([]*ballot.EncryptedBallot, error) {
	query := `
		SELECT ballot_id, election_id, voter_id, ciphertext, zk_proof, voter_pubkey, 
			   nullifier, signature, status, anchored_at
		FROM encrypted_ballots
		WHERE election_id = $1
		ORDER BY anchored_at ASC`

	rows, err := r.db.Query(query, electionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query encrypted ballots: %v", err)
	}
	defer rows.Close()

	var ballots []*ballot.EncryptedBallot
	for rows.Next() {
		var encryptedBallot ballot.EncryptedBallot
		err := rows.Scan(
			&encryptedBallot.BallotID,
			&encryptedBallot.ElectionID,
			&encryptedBallot.VoterID,
			&encryptedBallot.Ciphertext,
			&encryptedBallot.ZKProof,
			&encryptedBallot.VoterPubkey,
			&encryptedBallot.Nullifier,
			&encryptedBallot.Signature,
			&encryptedBallot.Status,
			&encryptedBallot.AnchoredAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan encrypted ballot: %v", err)
		}
		ballots = append(ballots, &encryptedBallot)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating encrypted ballots: %v", err)
	}

	return ballots, nil
}
