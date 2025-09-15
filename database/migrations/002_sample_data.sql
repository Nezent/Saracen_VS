-- Sample data for testing the Saracen Voting System
-- This file is optional and can be used for development/testing

-- Insert sample voters
INSERT INTO voters (voter_id, name, age) VALUES
    (1, 'Alice Johnson', 22),
    (2, 'Bob Smith', 30),
    (3, 'Charlie Brown', 45),
    (4, 'Diana Prince', 28),
    (5, 'Edward Wilson', 55)
ON CONFLICT (voter_id) DO NOTHING;

-- Insert sample candidates
INSERT INTO candidates (candidate_id, name, party) VALUES
    (1, 'John Doe', 'Green Party'),
    (2, 'Jane Smith', 'Blue Party'),
    (3, 'Michael Johnson', 'Red Party'),
    (4, 'Sarah Davis', 'Yellow Party')
ON CONFLICT (candidate_id) DO NOTHING;

-- Insert sample ballot (Alice votes for Jane Smith)
INSERT INTO ballots (voter_id, candidate_id) VALUES
    (1, 2)
ON CONFLICT (voter_id) DO NOTHING;

-- Insert sample audit request
INSERT INTO audits (audit_id, type, risk_limit, status) VALUES
    ('rla_88a1', 'rla', 0.1, 'queued')
ON CONFLICT (audit_id) DO NOTHING;

-- Update system uptime start time
UPDATE system_settings 
SET value = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::TEXT 
WHERE key = 'uptime_start';