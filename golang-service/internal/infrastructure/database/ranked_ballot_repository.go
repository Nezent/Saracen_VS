package database

import (
	"database/sql"
	"fmt"

	"github.com/Nezent/Saracen_Voting_System/internal/domain/ballot"
)

// RankedBallotPostgresRepository implements the RankedBallotRepository interface
type RankedBallotPostgresRepository struct {
	db *sql.DB
}

// NewRankedBallotRepository creates a new ranked ballot repository
func NewRankedBallotRepository(db *sql.DB) ballot.RankedBallotRepository {
	return &RankedBallotPostgresRepository{db: db}
}

// Create stores a new ranked ballot with its rankings in a transaction
func (r *RankedBallotPostgresRepository) Create(rankedBallot *ballot.RankedBallot, rankings []ballot.BallotRanking) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert ranked ballot
	ballotQuery := `
		INSERT INTO ranked_ballots 
		(ballot_id, election_id, voter_id, timestamp, status)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(
		ballotQuery,
		rankedBallot.BallotID,
		rankedBallot.ElectionID,
		rankedBallot.VoterID,
		rankedBallot.Timestamp,
		rankedBallot.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create ranked ballot: %v", err)
	}

	// Insert ballot rankings
	if len(rankings) > 0 {
		rankingQuery := `
			INSERT INTO ballot_rankings 
			(ballot_id, candidate_id, rank_position)
			VALUES ($1, $2, $3)`

		for _, ranking := range rankings {
			_, err = tx.Exec(
				rankingQuery,
				ranking.BallotID,
				ranking.CandidateID,
				ranking.RankPosition,
			)
			if err != nil {
				return fmt.Errorf("failed to create ballot ranking: %v", err)
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetByBallotID retrieves a ranked ballot with its rankings by ballot ID
func (r *RankedBallotPostgresRepository) GetByBallotID(ballotID string) (*ballot.RankedBallot, []ballot.BallotRanking, error) {
	// Get the ballot
	ballotQuery := `
		SELECT ballot_id, election_id, voter_id, timestamp, status
		FROM ranked_ballots
		WHERE ballot_id = $1`

	row := r.db.QueryRow(ballotQuery, ballotID)

	var rankedBallot ballot.RankedBallot
	err := row.Scan(
		&rankedBallot.BallotID,
		&rankedBallot.ElectionID,
		&rankedBallot.VoterID,
		&rankedBallot.Timestamp,
		&rankedBallot.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("ranked ballot not found: %s", ballotID)
		}
		return nil, nil, fmt.Errorf("failed to get ranked ballot: %v", err)
	}

	// Get the rankings
	rankingsQuery := `
		SELECT id, ballot_id, candidate_id, rank_position
		FROM ballot_rankings
		WHERE ballot_id = $1
		ORDER BY rank_position ASC`

	rows, err := r.db.Query(rankingsQuery, ballotID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query ballot rankings: %v", err)
	}
	defer rows.Close()

	var rankings []ballot.BallotRanking
	for rows.Next() {
		var ranking ballot.BallotRanking
		err := rows.Scan(
			&ranking.ID,
			&ranking.BallotID,
			&ranking.CandidateID,
			&ranking.RankPosition,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan ballot ranking: %v", err)
		}
		rankings = append(rankings, ranking)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating ballot rankings: %v", err)
	}

	return &rankedBallot, rankings, nil
}

// GetByElectionID retrieves all ranked ballots with rankings for an election
func (r *RankedBallotPostgresRepository) GetByElectionID(electionID string) ([]ballot.RankedBallotWithRankings, error) {
	query := `
		SELECT rb.ballot_id, rb.election_id, rb.voter_id, rb.timestamp, rb.status,
			   br.id, br.ballot_id, br.candidate_id, br.rank_position
		FROM ranked_ballots rb
		LEFT JOIN ballot_rankings br ON rb.ballot_id = br.ballot_id
		WHERE rb.election_id = $1
		ORDER BY rb.ballot_id, br.rank_position ASC`

	rows, err := r.db.Query(query, electionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ranked ballots: %v", err)
	}
	defer rows.Close()

	ballotMap := make(map[string]*ballot.RankedBallotWithRankings)

	for rows.Next() {
		var ballotData ballot.RankedBallot
		var rankingData ballot.BallotRanking
		var rankingID sql.NullInt32
		var rankingBallotID sql.NullString
		var candidateID sql.NullInt32
		var rankPosition sql.NullInt32

		err := rows.Scan(
			&ballotData.BallotID,
			&ballotData.ElectionID,
			&ballotData.VoterID,
			&ballotData.Timestamp,
			&ballotData.Status,
			&rankingID,
			&rankingBallotID,
			&candidateID,
			&rankPosition,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ranked ballot: %v", err)
		}

		// Check if ballot exists in map
		if _, exists := ballotMap[ballotData.BallotID]; !exists {
			ballotMap[ballotData.BallotID] = &ballot.RankedBallotWithRankings{
				Ballot:   ballotData,
				Rankings: []ballot.BallotRanking{},
			}
		}

		// Add ranking if it exists (LEFT JOIN might return null rankings)
		if rankingID.Valid {
			rankingData.ID = int(rankingID.Int32)
			rankingData.BallotID = rankingBallotID.String
			rankingData.CandidateID = int(candidateID.Int32)
			rankingData.RankPosition = int(rankPosition.Int32)
			ballotMap[ballotData.BallotID].Rankings = append(ballotMap[ballotData.BallotID].Rankings, rankingData)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ranked ballots: %v", err)
	}

	// Convert map to slice
	var results []ballot.RankedBallotWithRankings
	for _, ballotWithRankings := range ballotMap {
		results = append(results, *ballotWithRankings)
	}

	return results, nil
}

// GetByVoterID retrieves all ranked ballots for a voter
func (r *RankedBallotPostgresRepository) GetByVoterID(voterID int) ([]*ballot.RankedBallot, error) {
	query := `
		SELECT ballot_id, election_id, voter_id, timestamp, status
		FROM ranked_ballots
		WHERE voter_id = $1
		ORDER BY timestamp DESC`

	rows, err := r.db.Query(query, voterID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ranked ballots: %v", err)
	}
	defer rows.Close()

	var ballots []*ballot.RankedBallot
	for rows.Next() {
		var rankedBallot ballot.RankedBallot
		err := rows.Scan(
			&rankedBallot.BallotID,
			&rankedBallot.ElectionID,
			&rankedBallot.VoterID,
			&rankedBallot.Timestamp,
			&rankedBallot.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ranked ballot: %v", err)
		}
		ballots = append(ballots, &rankedBallot)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ranked ballots: %v", err)
	}

	return ballots, nil
}
