-- PostgreSQL Schema for Job Processing System
-- Generated for local testing and validation

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing table if exists (for testing)
DROP TABLE IF EXISTS jobs CASCADE;

-- Create jobs table
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    payload JSONB NOT NULL,
    result JSONB,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC);
CREATE INDEX idx_jobs_payload ON jobs USING GIN (payload);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_jobs_updated_at BEFORE UPDATE
    ON jobs FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Sample data for testing
INSERT INTO jobs (status, payload) VALUES 
    ('pending', '{"type": "email", "to": "user@example.com", "subject": "Test"}'),
    ('processing', '{"type": "report", "format": "pdf", "data": [1,2,3]}'),
    ('completed', '{"type": "backup", "source": "/data", "destination": "s3://bucket"}'),
    ('failed', '{"type": "export", "format": "csv", "query": "SELECT * FROM users"}');

-- Verify schema
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'jobs'
ORDER BY ordinal_position;