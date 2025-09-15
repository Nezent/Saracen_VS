-- Database setup for Saracen Voting System
-- This file will be executed when the PostgreSQL container starts

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Voters table
CREATE TABLE IF NOT EXISTS voters (
    voter_id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL CHECK (age >= 18),
    has_voted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on voter_id for performance
CREATE INDEX IF NOT EXISTS idx_voters_voter_id ON voters(voter_id);

-- Candidates table
CREATE TABLE IF NOT EXISTS candidates (
    candidate_id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    party VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on candidate_id for performance
CREATE INDEX IF NOT EXISTS idx_candidates_candidate_id ON candidates(candidate_id);
CREATE INDEX IF NOT EXISTS idx_candidates_party ON candidates(party);

-- Ballots table
CREATE TABLE IF NOT EXISTS ballots (
    ballot_id VARCHAR(20) PRIMARY KEY DEFAULT 'bal_' || nextval('ballot_seq'),
    voter_id INTEGER NOT NULL REFERENCES voters(voter_id) ON DELETE CASCADE,
    candidate_id INTEGER NOT NULL REFERENCES candidates(candidate_id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'accepted',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_voter_ballot UNIQUE (voter_id)
);

-- Create sequence for ballot IDs
CREATE SEQUENCE IF NOT EXISTS ballot_seq START 1001;

-- Create indexes on ballots
CREATE INDEX IF NOT EXISTS idx_ballots_voter_id ON ballots(voter_id);
CREATE INDEX IF NOT EXISTS idx_ballots_candidate_id ON ballots(candidate_id);
CREATE INDEX IF NOT EXISTS idx_ballots_status ON ballots(status);

-- Audits table
CREATE TABLE IF NOT EXISTS audits (
    audit_id VARCHAR(20) PRIMARY KEY,
    type VARCHAR(10) NOT NULL CHECK (type IN ('rla')),
    risk_limit DECIMAL(3,2) CHECK (risk_limit > 0 AND risk_limit <= 1),
    status VARCHAR(20) DEFAULT 'queued' CHECK (status IN ('queued', 'planned', 'in_progress', 'completed')),
    initial_sample_size INTEGER,
    sampling_plan TEXT,
    test VARCHAR(50),
    findings TEXT,
    p_value DECIMAL(10,9) CHECK (p_value >= 0 AND p_value <= 1),
    passed BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Audit samples table
CREATE TABLE IF NOT EXISTS audit_samples (
    id SERIAL PRIMARY KEY,
    audit_id VARCHAR(20) NOT NULL REFERENCES audits(audit_id) ON DELETE CASCADE,
    sample_numbers INTEGER[] NOT NULL,
    accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ballot manifests table
CREATE TABLE IF NOT EXISTS ballot_manifests (
    id SERIAL PRIMARY KEY,
    audit_id VARCHAR(20) NOT NULL REFERENCES audits(audit_id) ON DELETE CASCADE,
    counties TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit tables
CREATE INDEX IF NOT EXISTS idx_audits_status ON audits(status);
CREATE INDEX IF NOT EXISTS idx_audits_type ON audits(type);
CREATE INDEX IF NOT EXISTS idx_audit_samples_audit_id ON audit_samples(audit_id);
CREATE INDEX IF NOT EXISTS idx_ballot_manifests_audit_id ON ballot_manifests(audit_id);

-- System settings table for admin functionality
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial system settings
INSERT INTO system_settings (key, value) VALUES 
    ('system_initialized', 'true'),
    ('uptime_start', EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::TEXT)
ON CONFLICT (key) DO NOTHING;

-- Function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_voters_updated_at BEFORE UPDATE ON voters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_candidates_updated_at BEFORE UPDATE ON candidates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_audits_updated_at BEFORE UPDATE ON audits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_system_settings_updated_at BEFORE UPDATE ON system_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically update voter has_voted status when ballot is cast
CREATE OR REPLACE FUNCTION update_voter_voted_status()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE voters 
    SET has_voted = TRUE 
    WHERE voter_id = NEW.voter_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for ballot insertion
CREATE TRIGGER update_voter_status_on_ballot AFTER INSERT ON ballots
    FOR EACH ROW EXECUTE FUNCTION update_voter_voted_status();

-- Create view for election results
CREATE OR REPLACE VIEW election_results AS
SELECT 
    c.candidate_id,
    c.name,
    c.party,
    COUNT(b.ballot_id) as votes
FROM candidates c
LEFT JOIN ballots b ON c.candidate_id = b.candidate_id
GROUP BY c.candidate_id, c.name, c.party
ORDER BY votes DESC;

-- Grant permissions (adjust as needed for your specific setup)
-- These are basic permissions; you might want to create specific roles
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO saracen_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO saracen_user;