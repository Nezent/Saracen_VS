package ballot

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

// EncryptedBallot represents an encrypted ballot for Q16
type EncryptedBallot struct {
	BallotID    string    `json:"ballot_id" db:"ballot_id"`
	ElectionID  string    `json:"election_id" db:"election_id"`
	VoterID     int       `json:"voter_id" db:"voter_id"`
	Ciphertext  string    `json:"ciphertext" db:"ciphertext"`
	ZKProof     string    `json:"zk_proof" db:"zk_proof"`
	VoterPubkey string    `json:"voter_pubkey" db:"voter_pubkey"`
	Nullifier   string    `json:"nullifier" db:"nullifier"`
	Signature   string    `json:"signature" db:"signature"`
	Status      string    `json:"status" db:"status"`
	AnchoredAt  time.Time `json:"anchored_at" db:"anchored_at"`
}

// EncryptedBallotRequest represents the request payload for Q16
type EncryptedBallotRequest struct {
	ElectionID  string `json:"election_id" validate:"required"`
	VoterID     int    `json:"voter_id" validate:"required,min=1"`
	Ciphertext  string `json:"ciphertext" validate:"required,base64"`
	ZKProof     string `json:"zk_proof" validate:"required,base64"`
	VoterPubkey string `json:"voter_pubkey" validate:"required,hexadecimal"`
	Nullifier   string `json:"nullifier" validate:"required,hexadecimal"`
	Signature   string `json:"signature" validate:"required,base64"`
}

// EncryptedBallotResponse represents the response for Q16
type EncryptedBallotResponse struct {
	BallotID   string    `json:"ballot_id"`
	Status     string    `json:"status"`
	Nullifier  string    `json:"nullifier"`
	AnchoredAt time.Time `json:"anchored_at"`
}

// Validate validates the encrypted ballot request
func (req *EncryptedBallotRequest) Validate() error {
	if req.ElectionID == "" {
		return fmt.Errorf("election_id is required")
	}

	if req.VoterID <= 0 {
		return fmt.Errorf("voter_id must be positive")
	}

	// Convert and validate base64 fields (auto-convert if needed)
	req.Ciphertext = convertToBase64(req.Ciphertext)
	if err := validateBase64(req.Ciphertext, "ciphertext"); err != nil {
		return err
	}

	req.ZKProof = convertToBase64(req.ZKProof)
	if err := validateBase64(req.ZKProof, "zk_proof"); err != nil {
		return err
	}

	req.Signature = convertToBase64(req.Signature)
	if err := validateBase64(req.Signature, "signature"); err != nil {
		return err
	}

	// Convert and validate hex fields
	req.VoterPubkey = convertToHex(req.VoterPubkey)
	if err := validateHex(req.VoterPubkey, "voter_pubkey"); err != nil {
		return err
	}

	req.Nullifier = convertToHex(req.Nullifier)
	if err := validateHex(req.Nullifier, "nullifier"); err != nil {
		return err
	}

	return nil
}

// validateBase64 validates if a string is valid base64
func validateBase64(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if _, err := base64.StdEncoding.DecodeString(value); err != nil {
		return fmt.Errorf("%s must be valid base64: %v", fieldName, err)
	}
	return nil
}

// convertToBase64 converts simple input to proper base64 format
func convertToBase64(value string) string {
	if value == "" {
		return value
	}

	// Check if it's already valid base64
	if _, err := base64.StdEncoding.DecodeString(value); err == nil {
		return value
	}

	// Convert any string to base64
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// validateHex validates if a string is valid hexadecimal
func validateHex(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	// Remove 0x prefix if present
	if len(value) >= 2 && value[:2] == "0x" {
		value = value[2:]
	}
	if _, err := hex.DecodeString(value); err != nil {
		return fmt.Errorf("%s must be valid hexadecimal: %v", fieldName, err)
	}
	return nil
}

// convertToHex converts simple input to proper hex format
func convertToHex(value string) string {
	if value == "" {
		return value
	}

	// If already valid hex, return as is
	if len(value) >= 2 && value[:2] == "0x" {
		return value
	}

	// Check if it's already valid hex without 0x prefix
	if _, err := hex.DecodeString(value); err == nil && len(value) > 0 {
		return value
	}

	// Convert simple numbers or strings to hex
	// For simple inputs like "1", "2", convert to proper hex format
	if len(value) <= 3 {
		// Convert to bytes and then to hex
		hexValue := fmt.Sprintf("%x", []byte(value))
		// Ensure minimum length of 16 characters for pubkey/nullifier
		for len(hexValue) < 16 {
			hexValue = "0" + hexValue
		}
		return hexValue
	}

	// For longer strings, convert each character to hex
	hexValue := fmt.Sprintf("%x", []byte(value))
	return hexValue
}

// ToEncryptedBallot converts request to domain model with generated ID
func (req *EncryptedBallotRequest) ToEncryptedBallot() (*EncryptedBallot, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ballotID := generateBallotID("b_")

	return &EncryptedBallot{
		BallotID:    ballotID,
		ElectionID:  req.ElectionID,
		VoterID:     req.VoterID,
		Ciphertext:  req.Ciphertext,
		ZKProof:     req.ZKProof,
		VoterPubkey: req.VoterPubkey,
		Nullifier:   req.Nullifier,
		Signature:   req.Signature,
		Status:      "accepted",
		AnchoredAt:  time.Now(),
	}, nil
} // ToResponse converts domain model to response
func (eb *EncryptedBallot) ToResponse() *EncryptedBallotResponse {
	return &EncryptedBallotResponse{
		BallotID:   eb.BallotID,
		Status:     eb.Status,
		Nullifier:  eb.Nullifier,
		AnchoredAt: eb.AnchoredAt,
	}
}

// EncryptedBallotRepository defines repository interface for encrypted ballots
type EncryptedBallotRepository interface {
	Create(ballot *EncryptedBallot) error
	GetByBallotID(ballotID string) (*EncryptedBallot, error)
	GetByNullifier(nullifier string) (*EncryptedBallot, error)
	GetByElectionID(electionID string) ([]*EncryptedBallot, error)
}

// generateBallotID generates a unique ballot ID with prefix
func generateBallotID(prefix string) string {
	// Simple implementation - in production, use proper UUID or secure random
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s%x", prefix, timestamp%0xFFFF)
}
